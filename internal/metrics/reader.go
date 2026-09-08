package metrics

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

const automationMilestoneKeyBatchSize = 900

type StoredEvent struct {
	ID             int64
	Type           EventType
	PayloadVersion int
	OccurredAt     time.Time
	CreatedAt      time.Time
	SessionID      string
	AutomationKey  string
	EvaluationID   string
	OperationID    string
	PayloadJSON    json.RawMessage
}

type Scope uint8

const (
	ScopeAll Scope = iota
	ScopeAutomations
	ScopeSystem
)

type EventFilter struct {
	BeforeID      int64
	AfterID       int64
	SessionID     string
	AutomationKey string
	EvaluationID  string
	EventTypes    []EventType
	Search        string
	Scope         Scope
	Limit         int
	OldestFirst   bool
}

type EvaluationSummary struct {
	EvaluationID  string
	SessionID     string
	AutomationKey string
	FirstID       int64
	LastID        int64
	FirstOccurred time.Time
	LastOccurred  time.Time
	EventCount    int
}

type EvaluationFilter struct {
	BeforeID      int64
	SessionID     string
	AutomationKey string
	Search        string
	Limit         int
}

type AutomationMilestones struct {
	AutomationKey string
	LastTriggered time.Time
	LastCompleted time.Time
}

type Reader struct {
	database *sql.DB
}

// Events returns committed rows. It defaults to the newest 100 rows and caps a
// request at 1,000 rows; BeforeID supports the UI's explicit "Load older"
// pagination model.
func (reader *Reader) Events(ctx context.Context, filter EventFilter) ([]StoredEvent, error) {
	limit := normalizedLimit(filter.Limit)
	query := strings.Builder{}
	query.WriteString(`
		SELECT id, event_type, payload_version, occurred_at, created_at,
		       session_id, automation_key, evaluation_id, operation_id, payload_json
		FROM events
		WHERE 1 = 1`)
	arguments := make([]any, 0, 12)

	appendEventPredicates(&query, &arguments, filter)
	if filter.OldestFirst {
		query.WriteString(` ORDER BY occurred_at ASC, id ASC LIMIT ?`)
	} else {
		query.WriteString(` ORDER BY id DESC LIMIT ?`)
	}
	arguments = append(arguments, limit)

	rows, err := reader.database.QueryContext(ctx, query.String(), arguments...)
	if err != nil {
		return nil, fmt.Errorf("query metrics events: %w", err)
	}
	defer rows.Close()

	events := make([]StoredEvent, 0, limit)
	for rows.Next() {
		event, err := scanEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate metrics events: %w", err)
	}
	return events, nil
}

// EvaluationTrace returns every committed fact for one evaluation in observed
// chronological order, including operations that arrived after callback return.
func (reader *Reader) EvaluationTrace(ctx context.Context, evaluationID string) ([]StoredEvent, error) {
	if strings.TrimSpace(evaluationID) == "" {
		return nil, errors.New("evaluation ID is required")
	}
	rows, err := reader.database.QueryContext(ctx, `
		SELECT id, event_type, payload_version, occurred_at, created_at,
		       session_id, automation_key, evaluation_id, operation_id, payload_json
		FROM events
		WHERE evaluation_id = ?
		ORDER BY occurred_at ASC, id ASC`, evaluationID)
	if err != nil {
		return nil, fmt.Errorf("query metrics evaluation trace: %w", err)
	}
	defer rows.Close()

	events := make([]StoredEvent, 0, 16)
	for rows.Next() {
		event, err := scanEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate metrics evaluation trace: %w", err)
	}
	return events, nil
}

// Evaluations returns committed evaluation anchors for the Activity UI. The UI
// can fetch a selected anchor's complete event stream with EvaluationTrace.
func (reader *Reader) Evaluations(ctx context.Context, filter EvaluationFilter) ([]EvaluationSummary, error) {
	limit := normalizedLimit(filter.Limit)
	query := strings.Builder{}
	query.WriteString(`
		SELECT evaluation_id,
		       COALESCE(MAX(session_id), ''),
		       COALESCE(MAX(automation_key), ''),
		       MIN(id), MAX(id), MIN(occurred_at), MAX(occurred_at), COUNT(*)
		FROM events
		WHERE evaluation_id IS NOT NULL`)
	arguments := make([]any, 0, 6)
	if filter.SessionID != "" {
		query.WriteString(` AND session_id = ?`)
		arguments = append(arguments, filter.SessionID)
	}
	if filter.AutomationKey != "" {
		query.WriteString(` AND automation_key = ?`)
		arguments = append(arguments, filter.AutomationKey)
	}
	if filter.Search != "" {
		pattern := "%" + escapeLike(filter.Search) + "%"
		query.WriteString(` AND (
			automation_key LIKE ? ESCAPE '\' OR
			event_type LIKE ? ESCAPE '\' OR
			payload_json LIKE ? ESCAPE '\'
		)`)
		arguments = append(arguments, pattern, pattern, pattern)
	}
	query.WriteString(` GROUP BY evaluation_id`)
	if filter.BeforeID > 0 {
		query.WriteString(` HAVING MAX(id) < ?`)
		arguments = append(arguments, filter.BeforeID)
	}
	query.WriteString(` ORDER BY MAX(id) DESC LIMIT ?`)
	arguments = append(arguments, limit)

	rows, err := reader.database.QueryContext(ctx, query.String(), arguments...)
	if err != nil {
		return nil, fmt.Errorf("query metrics evaluations: %w", err)
	}
	defer rows.Close()

	summaries := make([]EvaluationSummary, 0, limit)
	for rows.Next() {
		var summary EvaluationSummary
		var firstOccurred, lastOccurred int64
		if err := rows.Scan(
			&summary.EvaluationID,
			&summary.SessionID,
			&summary.AutomationKey,
			&summary.FirstID,
			&summary.LastID,
			&firstOccurred,
			&lastOccurred,
			&summary.EventCount,
		); err != nil {
			return nil, fmt.Errorf("scan metrics evaluation: %w", err)
		}
		summary.FirstOccurred = time.Unix(0, firstOccurred).UTC()
		summary.LastOccurred = time.Unix(0, lastOccurred).UTC()
		summaries = append(summaries, summary)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate metrics evaluations: %w", err)
	}
	return summaries, nil
}

// LatestAutomationMilestones returns the most recent persisted trigger and
// normal callback completion times for the requested automations. Large
// catalogs are queried in bounded batches to stay below SQLite bind limits.
func (reader *Reader) LatestAutomationMilestones(ctx context.Context, automationKeys []string) ([]AutomationMilestones, error) {
	keys := make([]string, 0, len(automationKeys))
	seen := make(map[string]struct{}, len(automationKeys))
	for _, key := range automationKeys {
		if key == "" {
			continue
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		keys = append(keys, key)
	}
	if len(keys) == 0 {
		return nil, nil
	}

	milestones := make([]AutomationMilestones, 0, len(keys))
	for start := 0; start < len(keys); start += automationMilestoneKeyBatchSize {
		end := min(start+automationMilestoneKeyBatchSize, len(keys))
		batch, err := reader.latestAutomationMilestonesBatch(ctx, keys[start:end])
		if err != nil {
			return nil, err
		}
		milestones = append(milestones, batch...)
	}
	sort.Slice(milestones, func(i, j int) bool {
		return milestones[i].AutomationKey < milestones[j].AutomationKey
	})
	return milestones, nil
}

func (reader *Reader) latestAutomationMilestonesBatch(ctx context.Context, keys []string) ([]AutomationMilestones, error) {
	query := strings.Builder{}
	query.WriteString(`
		SELECT automation_key,
		       MAX(CASE WHEN event_type = ? THEN occurred_at END),
		       MAX(CASE WHEN event_type = ? THEN occurred_at END)
		FROM events
		WHERE event_type IN (?, ?)
		  AND automation_key IN (`)
	arguments := make([]any, 0, len(keys)+4)
	arguments = append(arguments,
		EventAutomationTriggered,
		EventAutomationCallbackFinished,
		EventAutomationTriggered,
		EventAutomationCallbackFinished,
	)
	for index, key := range keys {
		if index > 0 {
			query.WriteByte(',')
		}
		query.WriteByte('?')
		arguments = append(arguments, key)
	}
	query.WriteString(`)
		GROUP BY automation_key
		ORDER BY automation_key`)

	rows, err := reader.database.QueryContext(ctx, query.String(), arguments...)
	if err != nil {
		return nil, fmt.Errorf("query latest automation milestones: %w", err)
	}

	milestones := make([]AutomationMilestones, 0, len(keys))
	for rows.Next() {
		var milestone AutomationMilestones
		var lastTriggered, lastCompleted sql.NullInt64
		if err := rows.Scan(&milestone.AutomationKey, &lastTriggered, &lastCompleted); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan latest automation milestones: %w", err)
		}
		if lastTriggered.Valid {
			milestone.LastTriggered = time.Unix(0, lastTriggered.Int64).UTC()
		}
		if lastCompleted.Valid {
			milestone.LastCompleted = time.Unix(0, lastCompleted.Int64).UTC()
		}
		milestones = append(milestones, milestone)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, fmt.Errorf("iterate latest automation milestones: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("close latest automation milestones: %w", err)
	}
	return milestones, nil
}

func appendEventPredicates(query *strings.Builder, arguments *[]any, filter EventFilter) {
	if filter.BeforeID > 0 {
		query.WriteString(` AND id < ?`)
		*arguments = append(*arguments, filter.BeforeID)
	}
	if filter.AfterID > 0 {
		query.WriteString(` AND id > ?`)
		*arguments = append(*arguments, filter.AfterID)
	}
	if filter.SessionID != "" {
		query.WriteString(` AND session_id = ?`)
		*arguments = append(*arguments, filter.SessionID)
	}
	if filter.AutomationKey != "" {
		query.WriteString(` AND automation_key = ?`)
		*arguments = append(*arguments, filter.AutomationKey)
	}
	if filter.EvaluationID != "" {
		query.WriteString(` AND evaluation_id = ?`)
		*arguments = append(*arguments, filter.EvaluationID)
	}
	if len(filter.EventTypes) > 0 {
		query.WriteString(` AND event_type IN (`)
		for index, eventType := range filter.EventTypes {
			if index > 0 {
				query.WriteByte(',')
			}
			query.WriteByte('?')
			*arguments = append(*arguments, eventType)
		}
		query.WriteByte(')')
	}
	if filter.Search != "" {
		pattern := "%" + escapeLike(filter.Search) + "%"
		query.WriteString(` AND (
			automation_key LIKE ? ESCAPE '\' OR
			event_type LIKE ? ESCAPE '\' OR
			payload_json LIKE ? ESCAPE '\'
		)`)
		*arguments = append(*arguments, pattern, pattern, pattern)
	}
	switch filter.Scope {
	case ScopeAll:
	case ScopeAutomations:
		query.WriteString(` AND automation_key IS NOT NULL`)
	case ScopeSystem:
		query.WriteString(` AND automation_key IS NULL`)
	default:
		query.WriteString(` AND 0 = 1`)
	}
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanEvent(row rowScanner) (StoredEvent, error) {
	var event StoredEvent
	var occurredAt, createdAt int64
	var sessionID, automationKey, evaluationID, operationID sql.NullString
	var payloadJSON string
	if err := row.Scan(
		&event.ID,
		&event.Type,
		&event.PayloadVersion,
		&occurredAt,
		&createdAt,
		&sessionID,
		&automationKey,
		&evaluationID,
		&operationID,
		&payloadJSON,
	); err != nil {
		return StoredEvent{}, fmt.Errorf("scan metrics event: %w", err)
	}
	event.OccurredAt = time.Unix(0, occurredAt).UTC()
	event.CreatedAt = time.Unix(0, createdAt).UTC()
	event.SessionID = sessionID.String
	event.AutomationKey = automationKey.String
	event.EvaluationID = evaluationID.String
	event.OperationID = operationID.String
	event.PayloadJSON = json.RawMessage(payloadJSON)
	return event, nil
}

func normalizedLimit(limit int) int {
	if limit <= 0 {
		return 100
	}
	if limit > 1000 {
		return 1000
	}
	return limit
}

func escapeLike(value string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return replacer.Replace(value)
}

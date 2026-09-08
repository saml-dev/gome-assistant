package metrics

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestOpenRecorderCreatesDatabaseAndConfiguresSQLite(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "nested", "metrics.sqlite")
	recorder, err := OpenRecorder(databasePath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = recorder.Close(nil) })

	assertPragmaString(t, recorder.store.writer, "journal_mode", "wal")
	assertPragmaInt(t, recorder.store.writer, "synchronous", 1)
	assertPragmaInt(t, recorder.store.writer, "busy_timeout", 5000)
	assertPragmaInt(t, recorder.store.writer, "foreign_keys", 1)

	connections := make([]*sql.Conn, 0, 3)
	for range 3 {
		connection, err := recorder.store.reader.Conn(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		connections = append(connections, connection)
		var foreignKeys, busyTimeout int
		if err := connection.QueryRowContext(context.Background(), `PRAGMA foreign_keys`).Scan(&foreignKeys); err != nil {
			t.Fatal(err)
		}
		if err := connection.QueryRowContext(context.Background(), `PRAGMA busy_timeout`).Scan(&busyTimeout); err != nil {
			t.Fatal(err)
		}
		if foreignKeys != 1 || busyTimeout != 5000 {
			t.Fatalf("connection settings = foreign_keys %d, busy_timeout %d", foreignKeys, busyTimeout)
		}
	}
	for _, connection := range connections {
		_ = connection.Close()
	}

	rows, err := recorder.store.reader.Query(`
		SELECT name FROM sqlite_master
		WHERE type = 'table' AND name NOT LIKE 'sqlite_%'
		ORDER BY name`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			t.Fatal(err)
		}
		tables = append(tables, table)
	}
	if !reflect.DeepEqual(tables, []string{"events", "schema_migrations"}) {
		t.Fatalf("application tables = %v", tables)
	}

	columns, err := recorder.store.reader.Query(`PRAGMA table_info(events)`)
	if err != nil {
		t.Fatal(err)
	}
	defer columns.Close()
	var names []string
	for columns.Next() {
		var cid, notNull, primaryKey int
		var name, columnType string
		var defaultValue any
		if err := columns.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			t.Fatal(err)
		}
		names = append(names, name)
	}
	wantColumns := []string{
		"id", "event_type", "payload_version", "occurred_at", "created_at",
		"session_id", "automation_key", "evaluation_id", "operation_id", "payload_json",
	}
	if !reflect.DeepEqual(names, wantColumns) {
		t.Fatalf("event columns = %v", names)
	}

	if _, err := OpenRecorder("relative.sqlite"); err == nil || !strings.Contains(err.Error(), "absolute") {
		t.Fatalf("relative open error = %v", err)
	}
}

func TestMigrationsRejectNewerSchemaWithoutModification(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "newer.sqlite")
	database, err := sql.Open("sqlite", sqliteDSN(databasePath))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`
		CREATE TABLE schema_migrations (version INTEGER PRIMARY KEY, applied_at INTEGER NOT NULL);
		INSERT INTO schema_migrations(version, applied_at) VALUES (?, ?);
		CREATE TABLE future_data (value TEXT)`, schemaVersion+1, time.Now().UnixNano()); err != nil {
		t.Fatal(err)
	}
	if err := database.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = OpenRecorder(databasePath)
	if !errors.Is(err, ErrNewerSchema) {
		t.Fatalf("open error = %v", err)
	}

	database, err = sql.Open("sqlite", sqliteDSN(databasePath))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	var version int
	if err := database.QueryRow(`SELECT version FROM schema_migrations`).Scan(&version); err != nil {
		t.Fatal(err)
	}
	if version != schemaVersion+1 {
		t.Fatalf("schema version changed to %d", version)
	}
	var eventTables int
	if err := database.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'events'`).Scan(&eventTables); err != nil {
		t.Fatal(err)
	}
	if eventTables != 0 {
		t.Fatal("newer-schema rejection created events table")
	}
}

func TestMigrationIsTransactional(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "broken.sqlite")
	database, err := sql.Open("sqlite", sqliteDSN(databasePath))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`CREATE TABLE events (wrong TEXT)`); err != nil {
		t.Fatal(err)
	}
	if err := database.Close(); err != nil {
		t.Fatal(err)
	}

	if _, err := OpenRecorder(databasePath); err == nil {
		t.Fatal("migration unexpectedly succeeded")
	}

	database, err = sql.Open("sqlite", sqliteDSN(databasePath))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	var migrationsTable int
	if err := database.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'schema_migrations'`).Scan(&migrationsTable); err != nil {
		t.Fatal(err)
	}
	if migrationsTable != 0 {
		t.Fatal("failed migration was not rolled back")
	}
}

func TestRestartPersistenceAndAutoincrement(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "metrics.sqlite")
	recorder := newTestRecorder(t, databasePath, recorderOptions{})
	changes, cancel := recorder.SubscribeChanges()
	defer cancel()
	if !recorder.Record(testEvent(EventSessionStarting, "first")) {
		t.Fatal("first event not admitted")
	}
	waitSignal(t, changes)
	events, err := recorder.Reader().Events(context.Background(), EventFilter{SessionID: "first"})
	if err != nil || len(events) != 1 {
		t.Fatalf("first query = %v, %v", events, err)
	}
	firstID := events[0].ID
	if err := recorder.Close(nil); err != nil {
		t.Fatal(err)
	}

	recorder = newTestRecorder(t, databasePath, recorderOptions{})
	changes, cancel = recorder.SubscribeChanges()
	defer cancel()
	if !recorder.Record(testEvent(EventSessionReady, "second")) {
		t.Fatal("second event not admitted")
	}
	waitSignal(t, changes)
	events, err = recorder.Reader().Events(context.Background(), EventFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 || events[0].ID <= firstID {
		t.Fatalf("persisted events = %#v", events)
	}
}

func TestAcceptedOrderDeterminesCommittedIDs(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "ordered.sqlite")
	entered := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	recorder := newTestRecorder(t, databasePath, recorderOptions{
		beforeWrite: func(ctx context.Context, kind writeKind) error {
			if kind == writeBatch {
				once.Do(func() { close(entered) })
				select {
				case <-release:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		},
	})

	for index := 1; index <= 20; index++ {
		event := testEvent(EventAutomationTriggered, "session")
		event.EvaluationID = fmt.Sprintf("evaluation-%02d", index)
		event.OccurredAt = time.Unix(0, int64(100-index))
		if !recorder.Record(event) {
			t.Fatalf("event %d was not admitted", index)
		}
	}
	<-entered
	close(release)
	if err := recorder.Close(nil); err != nil {
		t.Fatal(err)
	}

	database, err := sql.Open("sqlite", sqliteDSN(databasePath))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	rows, err := database.Query(`SELECT evaluation_id FROM events WHERE evaluation_id IS NOT NULL ORDER BY id`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	index := 1
	for rows.Next() {
		var evaluationID string
		if err := rows.Scan(&evaluationID); err != nil {
			t.Fatal(err)
		}
		want := fmt.Sprintf("evaluation-%02d", index)
		if evaluationID != want {
			t.Fatalf("row %d evaluation = %q, want %q", index, evaluationID, want)
		}
		index++
	}
	if index != 21 {
		t.Fatalf("read %d accepted events", index-1)
	}
}

func TestInvalidationOnlyAfterCommit(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "commit.sqlite")
	entered := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	recorder := newTestRecorder(t, databasePath, recorderOptions{
		beforeWrite: func(ctx context.Context, kind writeKind) error {
			if kind == writeBatch {
				once.Do(func() { close(entered) })
				select {
				case <-release:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		},
	})
	changes, cancel := recorder.SubscribeChanges()
	defer cancel()
	if !recorder.Record(testEvent(EventSessionReady, "session")) {
		t.Fatal("event not admitted")
	}
	<-entered
	select {
	case <-changes:
		t.Fatal("received invalidation before commit")
	case <-time.After(30 * time.Millisecond):
	}
	close(release)
	waitSignal(t, changes)
}

func TestQueueSaturationIsNonblockingAndAggregated(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "saturation.sqlite")
	entered := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	recorder := newTestRecorder(t, databasePath, recorderOptions{
		queueCapacity: 1,
		batchSize:     1,
		beforeWrite: func(ctx context.Context, kind writeKind) error {
			if kind == writeBatch {
				once.Do(func() { close(entered) })
				select {
				case <-release:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		},
	})
	if !recorder.Record(testEvent(EventSessionStarting, "one")) {
		t.Fatal("first event not admitted")
	}
	<-entered
	if !recorder.Record(testEvent(EventSessionReady, "two")) {
		t.Fatal("second event not admitted")
	}
	started := time.Now()
	if recorder.Record(testEvent(EventSessionReady, "three")) {
		t.Fatal("saturated event unexpectedly admitted")
	}
	if time.Since(started) > 100*time.Millisecond {
		t.Fatal("saturated producer blocked")
	}
	close(release)

	eventually(t, func() bool {
		events, err := recorder.Reader().Events(context.Background(), EventFilter{EventTypes: []EventType{EventMetricsEventsDropped}})
		if err != nil || len(events) == 0 {
			return false
		}
		var payload MetricsEventsDroppedPayload
		return json.Unmarshal(events[0].PayloadJSON, &payload) == nil && payload.Counts.QueueFull == 1
	})
	if recorder.Health().TotalLosses.QueueFull != 1 {
		t.Fatalf("health losses = %#v", recorder.Health().TotalLosses)
	}
}

func TestSerializationSafetyAndAggregation(t *testing.T) {
	value := map[string]any{"kept": "before"}
	encoded := EncodeValue(value)
	value["kept"] = "after"
	if string(encoded) != `{"kept":"before"}` {
		t.Fatalf("captured value changed: %s", encoded)
	}

	unavailable := EncodeValue(map[int]string{1: "silently stringified by encoding/json"})
	var unavailableValue Unavailable
	if err := json.Unmarshal(unavailable, &unavailableValue); err != nil {
		t.Fatal(err)
	}
	if !unavailableValue.TelemetryUnavailable || unavailableValue.GoType != "map[int]string" ||
		!strings.Contains(unavailableValue.Error, ErrNonStringMapKey.Error()) {
		t.Fatalf("unavailable = %s", unavailable)
	}
	if _, err := EncodePayload(map[int]string{1: "bad"}); !errors.Is(err, ErrNonStringMapKey) {
		t.Fatalf("non-string map error = %v", err)
	}
	panicValue := EncodeValue(panickingMarshaler{})
	if !strings.Contains(string(panicValue), `"_telemetry_unavailable":true`) {
		t.Fatalf("panic was not contained: %s", panicValue)
	}

	databasePath := filepath.Join(t.TempDir(), "serialization.sqlite")
	recorder := newTestRecorder(t, databasePath, recorderOptions{})
	good := testEvent(EventStateReadFinished, "session")
	good.Payload = StateReadFinishedPayload{
		Method: "get",
		Result: unavailable,
		Error:  nil,
	}
	if !recorder.Record(good) {
		t.Fatal("event containing Unavailable not admitted")
	}
	bad := testEvent(EventSessionReady, "session")
	bad.Payload = map[int]string{1: "bad"}
	if recorder.Record(bad) {
		t.Fatal("invalid payload admitted")
	}

	eventually(t, func() bool {
		events, err := recorder.Reader().Events(context.Background(), EventFilter{EventTypes: []EventType{EventMetricsEventsDropped}})
		if err != nil || len(events) == 0 {
			return false
		}
		var payload MetricsEventsDroppedPayload
		return json.Unmarshal(events[0].PayloadJSON, &payload) == nil && payload.Counts.SerializationError == 1
	})
}

func TestDatabaseDegradationAndRecovery(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "recovery.sqlite")
	var failBatch atomic.Bool
	failBatch.Store(true)
	var allowRecovery atomic.Bool
	recorder := newTestRecorder(t, databasePath, recorderOptions{
		probeInitial: 10 * time.Millisecond,
		probeMaximum: 20 * time.Millisecond,
		beforeWrite: func(_ context.Context, kind writeKind) error {
			if kind == writeBatch && failBatch.CompareAndSwap(true, false) {
				return errors.New("injected database failure")
			}
			if kind == writeRecovery && !allowRecovery.Load() {
				return errors.New("database still unavailable")
			}
			return nil
		},
	})
	if !recorder.Record(testEvent(EventSessionStarting, "failed")) {
		t.Fatal("event not admitted before failure")
	}
	eventually(t, func() bool { return recorder.Health().Degraded })
	if recorder.Record(testEvent(EventSessionReady, "during-degradation")) {
		t.Fatal("event admitted during degradation")
	}
	allowRecovery.Store(true)
	eventually(t, func() bool { return !recorder.Health().Degraded })

	if !recorder.Record(testEvent(EventSessionReady, "after-recovery")) {
		t.Fatal("event not admitted after recovery")
	}
	eventually(t, func() bool {
		events, err := recorder.Reader().Events(context.Background(), EventFilter{})
		if err != nil {
			return false
		}
		seen := map[EventType]bool{}
		for _, event := range events {
			seen[event.Type] = true
		}
		return seen[EventMetricsDegraded] && seen[EventMetricsEventsDropped] &&
			seen[EventMetricsRecovered] && seen[EventSessionReady]
	})
	if recorder.Health().TotalLosses.DatabaseError != 2 {
		t.Fatalf("database losses = %#v", recorder.Health().TotalLosses)
	}
}

func TestCloseIsBounded(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "bounded.sqlite")
	entered := make(chan struct{})
	var once sync.Once
	recorder := newTestRecorder(t, databasePath, recorderOptions{
		closeBudget:  40 * time.Millisecond,
		probeInitial: time.Second,
		beforeWrite: func(ctx context.Context, kind writeKind) error {
			if kind != writeBatch {
				return nil
			}
			once.Do(func() { close(entered) })
			<-ctx.Done()
			return ctx.Err()
		},
	})
	if !recorder.Record(testEvent(EventSessionStarting, "session")) {
		t.Fatal("event not admitted")
	}
	<-entered
	started := time.Now()
	err := recorder.Close(nil)
	elapsed := time.Since(started)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("close error = %v", err)
	}
	if elapsed > 500*time.Millisecond {
		t.Fatalf("close took %s", elapsed)
	}
	if recorder.Record(testEvent(EventSessionReady, "late")) {
		t.Fatal("event admitted after permanent close")
	}
	eventually(t, func() bool { return recorder.Health().Closed })
}

func TestReaderFiltersAndEvaluationTrace(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "reader.sqlite")
	recorder := newTestRecorder(t, databasePath, recorderOptions{})
	evaluationID := "evaluation-a"
	for index, eventType := range []EventType{
		EventAutomationTriggered,
		EventAutomationCallbackStarted,
		EventAutomationCallbackFinished,
	} {
		event := testEvent(eventType, "session-a")
		event.AutomationKey = "kitchen"
		event.EvaluationID = evaluationID
		event.OccurredAt = time.Unix(0, int64(3-index))
		if !recorder.Record(event) {
			t.Fatal("event not admitted")
		}
	}
	if !recorder.Record(testEvent(EventSessionReady, "session-a")) {
		t.Fatal("system event not admitted")
	}
	eventually(t, func() bool {
		events, err := recorder.Reader().Events(context.Background(), EventFilter{})
		return err == nil && len(events) == 4
	})

	trace, err := recorder.Reader().EvaluationTrace(context.Background(), evaluationID)
	if err != nil {
		t.Fatal(err)
	}
	if len(trace) != 3 || trace[0].OccurredAt.After(trace[1].OccurredAt) || trace[1].OccurredAt.After(trace[2].OccurredAt) {
		t.Fatalf("trace is not chronological: %#v", trace)
	}
	system, err := recorder.Reader().Events(context.Background(), EventFilter{Scope: ScopeSystem})
	if err != nil || len(system) != 1 || system[0].Type != EventSessionReady {
		t.Fatalf("system events = %#v, %v", system, err)
	}
	summaries, err := recorder.Reader().Evaluations(context.Background(), EvaluationFilter{Search: "kitchen"})
	if err != nil || len(summaries) != 1 || summaries[0].EventCount != 3 {
		t.Fatalf("evaluation summaries = %#v, %v", summaries, err)
	}
}

func TestLatestAutomationMilestonesGroupsRequestedKeys(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "milestones.sqlite")
	recorder := newTestRecorder(t, databasePath, recorderOptions{})
	recordEvent := func(eventType EventType, key string, occurredAt time.Time) {
		event := testEvent(eventType, "session")
		event.AutomationKey = key
		event.EvaluationID = key + "-evaluation"
		event.OccurredAt = occurredAt
		if !recorder.Record(event) {
			t.Fatalf("%s for %s was not admitted", eventType, key)
		}
	}

	firstTrigger := time.Date(2026, time.September, 1, 10, 0, 0, 0, time.UTC)
	completion := firstTrigger.Add(time.Minute)
	lastTrigger := completion.Add(time.Minute)
	recordEvent(EventAutomationTriggered, "kitchen", firstTrigger)
	recordEvent(EventAutomationCallbackFinished, "kitchen", completion)
	recordEvent(EventAutomationTriggered, "kitchen", lastTrigger)
	recordEvent(EventAutomationSkipped, "kitchen", lastTrigger.Add(time.Minute))
	recordEvent(EventAutomationCallbackStarted, "kitchen", lastTrigger.Add(2*time.Minute))
	recordEvent(EventAutomationCallbackFinished, "completion-only", completion.Add(time.Hour))
	recordEvent(EventAutomationTriggered, "not-requested", lastTrigger.Add(time.Hour))
	eventually(t, func() bool {
		events, err := recorder.Reader().Events(context.Background(), EventFilter{})
		return err == nil && len(events) == 7
	})

	milestones, err := recorder.Reader().LatestAutomationMilestones(context.Background(), []string{
		"kitchen", "completion-only", "missing", "kitchen", "",
	})
	if err != nil {
		t.Fatal(err)
	}
	byKey := make(map[string]AutomationMilestones, len(milestones))
	for _, milestone := range milestones {
		byKey[milestone.AutomationKey] = milestone
	}
	if len(byKey) != 2 {
		t.Fatalf("milestones = %#v, want only requested keys with milestones", milestones)
	}
	if got := byKey["kitchen"]; !got.LastTriggered.Equal(lastTrigger) || !got.LastCompleted.Equal(completion) {
		t.Fatalf("kitchen milestones = %#v, want trigger %s and completion %s", got, lastTrigger, completion)
	}
	if got := byKey["completion-only"]; !got.LastTriggered.IsZero() || !got.LastCompleted.Equal(completion.Add(time.Hour)) {
		t.Fatalf("completion-only milestones = %#v", got)
	}
	empty, err := recorder.Reader().LatestAutomationMilestones(context.Background(), nil)
	if err != nil || len(empty) != 0 {
		t.Fatalf("empty milestone query = %#v, %v", empty, err)
	}
}

func TestLatestAutomationMilestonesBatchesLargeCatalogs(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "milestones-batched.sqlite")
	recorder := newTestRecorder(t, databasePath, recorderOptions{})
	keys := make([]string, automationMilestoneKeyBatchSize+1)
	for index := range keys {
		keys[index] = fmt.Sprintf("automation-%04d", index)
	}

	firstTime := time.Date(2026, time.September, 2, 10, 0, 0, 0, time.UTC)
	lastTime := firstTime.Add(time.Minute)
	for index, occurredAt := range []time.Time{firstTime, lastTime} {
		event := testEvent(EventAutomationTriggered, "session")
		event.AutomationKey = keys[index*automationMilestoneKeyBatchSize]
		event.EvaluationID = fmt.Sprintf("evaluation-%d", index)
		event.OccurredAt = occurredAt
		if !recorder.Record(event) {
			t.Fatalf("trigger for %s was not admitted", event.AutomationKey)
		}
	}
	eventually(t, func() bool {
		events, err := recorder.Reader().Events(context.Background(), EventFilter{})
		return err == nil && len(events) == 2
	})

	milestones, err := recorder.Reader().LatestAutomationMilestones(context.Background(), keys)
	if err != nil {
		t.Fatal(err)
	}
	if len(milestones) != 2 {
		t.Fatalf("milestones = %#v, want results from both batches", milestones)
	}
	if milestones[0].AutomationKey != keys[0] || !milestones[0].LastTriggered.Equal(firstTime) {
		t.Fatalf("first batch milestone = %#v", milestones[0])
	}
	if milestones[1].AutomationKey != keys[automationMilestoneKeyBatchSize] || !milestones[1].LastTriggered.Equal(lastTime) {
		t.Fatalf("second batch milestone = %#v", milestones[1])
	}
}

type panickingMarshaler struct{}

func (panickingMarshaler) MarshalJSON() ([]byte, error) {
	panic("serializer panic")
}

func newTestRecorder(t *testing.T, path string, overrides recorderOptions) *Recorder {
	t.Helper()
	defaults := defaultRecorderOptions()
	if overrides.queueCapacity == 0 {
		overrides.queueCapacity = defaults.queueCapacity
	}
	if overrides.batchSize == 0 {
		overrides.batchSize = defaults.batchSize
	}
	if overrides.closeBudget == 0 {
		overrides.closeBudget = 500 * time.Millisecond
	}
	if overrides.probeInitial == 0 {
		overrides.probeInitial = 10 * time.Millisecond
	}
	if overrides.probeMaximum == 0 {
		overrides.probeMaximum = 20 * time.Millisecond
	}
	if overrides.now == nil {
		overrides.now = defaults.now
	}
	if overrides.logger == nil {
		overrides.logger = discardLogger()
	}
	recorder, err := openRecorder(context.Background(), path, overrides)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if !recorder.Health().Closed {
			_ = recorder.Close(nil)
		}
	})
	return recorder
}

func testEvent(eventType EventType, sessionID string) Event {
	return Event{
		Type:       eventType,
		OccurredAt: time.Now().UTC(),
		SessionID:  sessionID,
		Payload:    EmptyPayload{},
	}
}

func waitSignal(t *testing.T, changes <-chan struct{}) {
	t.Helper()
	select {
	case _, ok := <-changes:
		if !ok {
			t.Fatal("change subscription closed")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for committed-change signal")
	}
}

func eventually(t *testing.T, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("condition did not become true")
}

func assertPragmaInt(t *testing.T, database *sql.DB, name string, want int) {
	t.Helper()
	var got int
	if err := database.QueryRow(`PRAGMA ` + name).Scan(&got); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("PRAGMA %s = %d, want %d", name, got, want)
	}
}

func assertPragmaString(t *testing.T, database *sql.DB, name, want string) {
	t.Helper()
	var got string
	if err := database.QueryRow(`PRAGMA ` + name).Scan(&got); err != nil {
		t.Fatal(err)
	}
	if strings.ToLower(got) != want {
		t.Fatalf("PRAGMA %s = %q, want %q", name, got, want)
	}
}

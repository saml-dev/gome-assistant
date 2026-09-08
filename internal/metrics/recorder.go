package metrics

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultQueueCapacity = 1024
	defaultBatchSize     = 256
	defaultCloseBudget   = 5 * time.Second
)

type Event struct {
	Type          EventType
	OccurredAt    time.Time
	SessionID     string
	AutomationKey string
	EvaluationID  string
	OperationID   string
	Payload       any
}

type encodedEvent struct {
	eventType      EventType
	payloadVersion int
	occurredAt     int64
	sessionID      sql.NullString
	automationKey  sql.NullString
	evaluationID   sql.NullString
	operationID    sql.NullString
	payloadJSON    json.RawMessage
}

type writeKind uint8

const (
	writeBatch writeKind = iota
	writeRecovery
)

// recorderOptions is deliberately private. It is a narrow deterministic seam
// for package tests; production always uses the defaults in OpenRecorder.
type recorderOptions struct {
	queueCapacity int
	batchSize     int
	closeBudget   time.Duration
	probeInitial  time.Duration
	probeMaximum  time.Duration
	now           func() time.Time
	logger        *slog.Logger
	beforeWrite   func(context.Context, writeKind) error
}

func defaultRecorderOptions() recorderOptions {
	return recorderOptions{
		queueCapacity: defaultQueueCapacity,
		batchSize:     defaultBatchSize,
		closeBudget:   defaultCloseBudget,
		probeInitial:  100 * time.Millisecond,
		probeMaximum:  5 * time.Second,
		now:           time.Now,
		logger:        slog.Default(),
	}
}

type Recorder struct {
	store   *store
	reader  *Reader
	options recorderOptions

	queue    chan encodedEvent
	control  chan closeRequest
	lossWake chan struct{}
	done     chan struct{}

	changes *broadcaster

	admissionMu sync.RWMutex
	accepting   bool
	degraded    atomic.Bool

	losses lossLedger

	healthMu sync.RWMutex
	health   Health

	logMu       sync.Mutex
	lastLossLog time.Time

	context context.Context
	cancel  context.CancelFunc

	closeMu        sync.Mutex
	closeRequested bool
	closeErr       error
}

type closeRequest struct {
	deadline time.Time
	final    *encodedEvent
}

type Health struct {
	Degraded      bool
	DegradedAt    time.Time
	LastError     *Error
	PendingLosses LossCounts
	TotalLosses   LossCounts
	Closed        bool
}

func OpenRecorder(databasePath string) (*Recorder, error) {
	return openRecorder(context.Background(), databasePath, defaultRecorderOptions())
}

func openRecorder(ctx context.Context, databasePath string, options recorderOptions) (*Recorder, error) {
	defaults := defaultRecorderOptions()
	if options.queueCapacity <= 0 {
		options.queueCapacity = defaults.queueCapacity
	}
	if options.batchSize <= 0 {
		options.batchSize = defaults.batchSize
	}
	if options.closeBudget <= 0 {
		options.closeBudget = defaults.closeBudget
	}
	if options.probeInitial <= 0 {
		options.probeInitial = defaults.probeInitial
	}
	if options.probeMaximum < options.probeInitial {
		options.probeMaximum = defaults.probeMaximum
	}
	if options.now == nil {
		options.now = defaults.now
	}
	if options.logger == nil {
		options.logger = defaults.logger
	}

	database, err := openStore(ctx, databasePath)
	if err != nil {
		return nil, err
	}

	writerContext, cancel := context.WithCancel(context.Background())
	recorder := &Recorder{
		store:     database,
		reader:    &Reader{database: database.reader},
		options:   options,
		queue:     make(chan encodedEvent, options.queueCapacity),
		control:   make(chan closeRequest, 1),
		lossWake:  make(chan struct{}, 1),
		done:      make(chan struct{}),
		changes:   newBroadcaster(),
		accepting: true,
		context:   writerContext,
		cancel:    cancel,
	}
	go recorder.runWriter()
	return recorder, nil
}

// Reader exposes committed storage to the embedded UI without exposing the
// database handle or a JSON/HTTP API.
func (recorder *Recorder) Reader() *Reader {
	return recorder.reader
}

// SubscribeChanges returns coalesced signals emitted only after successful
// commits. Call cancel when the subscriber is no longer used.
func (recorder *Recorder) SubscribeChanges() (<-chan struct{}, func()) {
	return recorder.changes.subscribe()
}

func (recorder *Recorder) Health() Health {
	recorder.healthMu.RLock()
	health := recorder.health
	recorder.healthMu.RUnlock()
	health.PendingLosses = recorder.losses.pendingCounts()
	health.TotalLosses = recorder.losses.totalCounts()
	if health.LastError != nil {
		lastError := *health.LastError
		health.LastError = &lastError
	}
	return health
}

// Record serializes payload synchronously and then attempts a nonblocking
// admission. Its result reports admission only, never eventual durability.
func (recorder *Recorder) Record(event Event) bool {
	recorder.admissionMu.RLock()
	accepting := recorder.accepting
	recorder.admissionMu.RUnlock()
	if !accepting {
		return false
	}

	encoded, err := recorder.encode(event)
	if err != nil {
		// If Close raced the serialization, rejecting the event is intentional
		// and should not be reported as telemetry loss.
		recorder.admissionMu.RLock()
		stillAccepting := recorder.accepting
		recorder.admissionMu.RUnlock()
		if stillAccepting {
			recorder.noteLoss(lossSerialization, 1, recorder.options.now())
		}
		return false
	}

	recorder.admissionMu.RLock()
	defer recorder.admissionMu.RUnlock()
	if !recorder.accepting {
		return false
	}
	if recorder.degraded.Load() {
		recorder.noteLoss(lossDatabase, 1, recorder.options.now())
		return false
	}
	select {
	case recorder.queue <- encoded:
		return true
	default:
		recorder.noteLoss(lossQueue, 1, recorder.options.now())
		return false
	}
}

func (recorder *Recorder) encode(event Event) (encodedEvent, error) {
	if !validEventType(event.Type) {
		return encodedEvent{}, fmt.Errorf("unknown telemetry event type %q", event.Type)
	}
	payloadJSON, err := EncodePayload(event.Payload)
	if err != nil {
		return encodedEvent{}, fmt.Errorf("encode %s payload: %w", event.Type, err)
	}
	occurredAt := event.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = recorder.options.now()
	}
	return encodedEvent{
		eventType:      event.Type,
		payloadVersion: 1,
		occurredAt:     unixNano(occurredAt),
		sessionID:      nullable(event.SessionID),
		automationKey:  nullable(event.AutomationKey),
		evaluationID:   nullable(event.EvaluationID),
		operationID:    nullable(event.OperationID),
		payloadJSON:    payloadJSON,
	}, nil
}

func nullable(value string) sql.NullString {
	return sql.NullString{String: value, Valid: value != ""}
}

// Close permanently stops external admission and makes a bounded effort to
// flush accepted events. final, when non-nil, is internally admitted after all
// external events. Production's fixed budget is five seconds.
func (recorder *Recorder) Close(final *Event) error {
	return recorder.closeWithBudget(recorder.options.closeBudget, final)
}

func (recorder *Recorder) closeWithBudget(budget time.Duration, final *Event) error {
	deadline := time.Now().Add(budget)
	var encodedFinal *encodedEvent
	if final != nil {
		encoded, err := recorder.encode(*final)
		if err != nil {
			recorder.noteLoss(lossSerialization, 1, recorder.options.now())
		} else {
			encodedFinal = &encoded
		}
	}

	recorder.closeMu.Lock()
	if !recorder.closeRequested {
		recorder.admissionMu.Lock()
		recorder.accepting = false
		recorder.admissionMu.Unlock()
		recorder.closeRequested = true
		recorder.control <- closeRequest{deadline: deadline, final: encodedFinal}
	}
	recorder.closeMu.Unlock()

	timer := time.NewTimer(time.Until(deadline))
	defer timer.Stop()
	select {
	case <-recorder.done:
		recorder.closeMu.Lock()
		err := recorder.closeErr
		recorder.closeMu.Unlock()
		return err
	case <-timer.C:
		recorder.cancel()
		return context.DeadlineExceeded
	}
}

func (recorder *Recorder) runWriter() {
	defer func() {
		recorder.cancel()
		recorder.healthMu.Lock()
		recorder.health.Closed = true
		recorder.healthMu.Unlock()
		recorder.changes.close()
		_ = recorder.store.writer.Close()
		go func() {
			_ = recorder.store.reader.Close()
		}()
		close(recorder.done)
	}()

	for {
		// Once permanent close has begun, prioritize its deadline and final
		// internal event over the writer context's simultaneous cancellation.
		select {
		case request := <-recorder.control:
			recorder.setCloseError(recorder.shutdown(request))
			return
		default:
		}
		if recorder.degraded.Load() {
			if recorder.runDegraded() {
				return
			}
			continue
		}

		select {
		case first := <-recorder.queue:
			batch := recorder.drainBatch(first)
			if err := recorder.commitBatch(recorder.context, batch); err != nil {
				recorder.noteLoss(lossDatabase, uint64(len(batch)), recorder.options.now())
				recorder.markDegraded(err)
				recorder.dropQueuedDatabaseLoss()
			}
		case <-recorder.lossWake:
			if err := recorder.commitLossSummary(recorder.context); err != nil {
				recorder.markDegraded(err)
				recorder.dropQueuedDatabaseLoss()
			}
		case request := <-recorder.control:
			recorder.setCloseError(recorder.shutdown(request))
			return
		case <-recorder.context.Done():
			recorder.dropQueuedDatabaseLoss()
			recorder.setCloseError(recorder.context.Err())
			return
		}
	}
}

func (recorder *Recorder) drainBatch(first encodedEvent) []encodedEvent {
	batch := make([]encodedEvent, 0, recorder.options.batchSize)
	batch = append(batch, first)
	for len(batch) < recorder.options.batchSize {
		select {
		case event := <-recorder.queue:
			batch = append(batch, event)
		default:
			return batch
		}
	}
	return batch
}

func (recorder *Recorder) commitBatch(ctx context.Context, batch []encodedEvent) error {
	if len(batch) == 0 {
		return nil
	}
	snapshot := recorder.losses.take()
	events := make([]encodedEvent, 0, len(batch)+1)
	events = append(events, batch...)
	if snapshot.counts.total() > 0 {
		events = append(events, recorder.lossEvent(snapshot))
	}
	if err := recorder.write(ctx, writeBatch, events); err != nil {
		recorder.losses.restore(snapshot)
		return err
	}
	recorder.changes.publish()
	return nil
}

func (recorder *Recorder) commitLossSummary(ctx context.Context) error {
	snapshot := recorder.losses.take()
	if snapshot.counts.total() == 0 {
		return nil
	}
	if err := recorder.write(ctx, writeBatch, []encodedEvent{recorder.lossEvent(snapshot)}); err != nil {
		recorder.losses.restore(snapshot)
		return err
	}
	recorder.changes.publish()
	return nil
}

func (recorder *Recorder) write(ctx context.Context, kind writeKind, events []encodedEvent) error {
	if recorder.options.beforeWrite != nil {
		if err := recorder.options.beforeWrite(ctx, kind); err != nil {
			return err
		}
	}
	tx, err := recorder.store.writer.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin metrics event transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	statement, err := tx.PrepareContext(ctx, `
		INSERT INTO events(
			event_type, payload_version, occurred_at, created_at,
			session_id, automation_key, evaluation_id, operation_id, payload_json
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("prepare metrics event insert: %w", err)
	}
	defer statement.Close()

	for _, event := range events {
		if _, err := statement.ExecContext(ctx,
			event.eventType,
			event.payloadVersion,
			event.occurredAt,
			recorder.options.now().UTC().UnixNano(),
			event.sessionID,
			event.automationKey,
			event.evaluationID,
			event.operationID,
			string(event.payloadJSON),
		); err != nil {
			return fmt.Errorf("insert metrics event: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit metrics event transaction: %w", err)
	}
	return nil
}

func (recorder *Recorder) markDegraded(err error) {
	now := recorder.options.now().UTC()
	recorder.degraded.Store(true)
	recorder.healthMu.Lock()
	if !recorder.health.Degraded {
		recorder.health.Degraded = true
		recorder.health.DegradedAt = now
		recorder.health.LastError = ErrorValue(err)
	}
	recorder.healthMu.Unlock()
	recorder.options.logger.Error("metrics storage degraded", "error", err)
}

func (recorder *Recorder) runDegraded() bool {
	backoff := recorder.options.probeInitial
	timer := time.NewTimer(backoff)
	defer timer.Stop()

	for recorder.degraded.Load() {
		select {
		case request := <-recorder.control:
			recorder.setCloseError(recorder.shutdown(request))
			return true
		default:
		}
		select {
		case <-timer.C:
			probeContext, cancel := context.WithTimeout(recorder.context, defaultCloseBudget)
			err := recorder.recover(probeContext)
			cancel()
			if err == nil {
				return false
			}
			backoff *= 2
			if backoff > recorder.options.probeMaximum {
				backoff = recorder.options.probeMaximum
			}
			timer.Reset(backoff)
		case event := <-recorder.queue:
			recorder.noteLoss(lossDatabase, 1, recorder.options.now())
			_ = event
		case <-recorder.lossWake:
			// Recovery writes the pending aggregate directly.
		case request := <-recorder.control:
			recorder.setCloseError(recorder.shutdown(request))
			return true
		case <-recorder.context.Done():
			recorder.dropQueuedDatabaseLoss()
			recorder.setCloseError(recorder.context.Err())
			return true
		}
	}
	return false
}

func (recorder *Recorder) recover(ctx context.Context) error {
	recorder.healthMu.RLock()
	degradedAt := recorder.health.DegradedAt
	degradedError := recorder.health.LastError
	recorder.healthMu.RUnlock()
	if degradedError == nil {
		degradedError = &Error{GoType: "unknown", Message: "metrics storage unavailable"}
	}

	snapshot := recorder.losses.take()
	now := recorder.options.now().UTC()
	events := []encodedEvent{
		recorder.internalEvent(EventMetricsDegraded, degradedAt, MetricsDegradedPayload{Error: *degradedError}),
	}
	if snapshot.counts.total() > 0 {
		events = append(events, recorder.lossEvent(snapshot))
	}
	events = append(events,
		recorder.internalEvent(EventMetricsRecovered, now, MetricsRecoveredPayload{DegradedOccurredAt: degradedAt.UnixNano()}),
	)

	if err := recorder.write(ctx, writeRecovery, events); err != nil {
		recorder.losses.restore(snapshot)
		return err
	}
	recorder.healthMu.Lock()
	recorder.health.Degraded = false
	recorder.health.LastError = nil
	recorder.healthMu.Unlock()
	recorder.degraded.Store(false)
	recorder.options.logger.Info("metrics storage recovered", "degraded_at", degradedAt)
	recorder.changes.publish()
	return nil
}

func (recorder *Recorder) shutdown(request closeRequest) error {
	ctx, cancel := context.WithDeadline(recorder.context, request.deadline)
	defer cancel()

	if recorder.degraded.Load() {
		remaining := recorder.dropQueuedDatabaseLoss()
		if err := recorder.recover(ctx); err != nil {
			if request.final != nil {
				recorder.noteLoss(lossDatabase, 1, recorder.options.now())
				remaining++
			}
			if remaining > 0 {
				recorder.options.logger.Error("metrics close dropped accepted events", "count", remaining, "error", err)
			}
			return err
		}
	}

	events := recorder.drainAll()
	if request.final != nil {
		events = append(events, *request.final)
	}

	for len(events) > 0 {
		count := recorder.options.batchSize
		if count > len(events) {
			count = len(events)
		}
		batch := events[:count]
		events = events[count:]
		if err := recorder.commitBatch(ctx, batch); err != nil {
			remaining := uint64(len(batch) + len(events))
			recorder.noteLoss(lossDatabase, remaining, recorder.options.now())
			recorder.markDegraded(err)
			recorder.options.logger.Error("metrics close dropped accepted events", "count", remaining, "error", err)
			if recoveryErr := recorder.recover(ctx); recoveryErr != nil {
				return errors.Join(err, recoveryErr)
			}
			break
		}
	}

	if err := recorder.commitLossSummary(ctx); err != nil {
		return err
	}
	return nil
}

func (recorder *Recorder) drainAll() []encodedEvent {
	events := make([]encodedEvent, 0, len(recorder.queue))
	for {
		select {
		case event := <-recorder.queue:
			events = append(events, event)
		default:
			return events
		}
	}
}

func (recorder *Recorder) dropQueuedDatabaseLoss() uint64 {
	count := uint64(len(recorder.drainAll()))
	if count > 0 {
		recorder.noteLoss(lossDatabase, count, recorder.options.now())
	}
	return count
}

func (recorder *Recorder) internalEvent(eventType EventType, occurredAt time.Time, payload any) encodedEvent {
	payloadJSON, err := EncodePayload(payload)
	if err != nil {
		panic(fmt.Sprintf("encode internal metrics payload: %v", err))
	}
	return encodedEvent{
		eventType:      eventType,
		payloadVersion: 1,
		occurredAt:     occurredAt.UTC().UnixNano(),
		payloadJSON:    payloadJSON,
	}
}

func (recorder *Recorder) lossEvent(snapshot lossSnapshot) encodedEvent {
	return recorder.internalEvent(EventMetricsEventsDropped, snapshot.endedAt, MetricsEventsDroppedPayload{
		PeriodStartedAt: snapshot.startedAt.UTC().UnixNano(),
		PeriodEndedAt:   snapshot.endedAt.UTC().UnixNano(),
		Counts:          snapshot.counts,
	})
}

func (recorder *Recorder) noteLoss(kind lossKind, count uint64, at time.Time) {
	if count == 0 {
		return
	}
	recorder.losses.add(kind, count, at.UTC())
	select {
	case recorder.lossWake <- struct{}{}:
	default:
	}

	recorder.logMu.Lock()
	if at.Sub(recorder.lastLossLog) >= time.Second {
		recorder.lastLossLog = at
		recorder.options.logger.Warn("metrics events dropped", "reason", kind.String(), "count", count)
	}
	recorder.logMu.Unlock()
}

func (recorder *Recorder) setCloseError(err error) {
	recorder.closeMu.Lock()
	if recorder.closeErr == nil {
		recorder.closeErr = err
	}
	recorder.closeMu.Unlock()
}

type lossKind uint8

const (
	lossQueue lossKind = iota
	lossDatabase
	lossSerialization
)

func (kind lossKind) String() string {
	switch kind {
	case lossQueue:
		return "queue_full"
	case lossDatabase:
		return "database_error"
	case lossSerialization:
		return "serialization_error"
	default:
		return "unknown"
	}
}

type lossSnapshot struct {
	counts    LossCounts
	startedAt time.Time
	endedAt   time.Time
}

type lossLedger struct {
	mu      sync.Mutex
	pending lossSnapshot
	total   LossCounts
}

func (ledger *lossLedger) add(kind lossKind, count uint64, at time.Time) {
	ledger.mu.Lock()
	defer ledger.mu.Unlock()
	if ledger.pending.counts.total() == 0 {
		ledger.pending.startedAt = at
	}
	ledger.pending.endedAt = at
	addLoss(&ledger.pending.counts, kind, count)
	addLoss(&ledger.total, kind, count)
}

func (ledger *lossLedger) take() lossSnapshot {
	ledger.mu.Lock()
	defer ledger.mu.Unlock()
	snapshot := ledger.pending
	ledger.pending = lossSnapshot{}
	return snapshot
}

func (ledger *lossLedger) restore(snapshot lossSnapshot) {
	if snapshot.counts.total() == 0 {
		return
	}
	ledger.mu.Lock()
	defer ledger.mu.Unlock()
	if ledger.pending.counts.total() == 0 || snapshot.startedAt.Before(ledger.pending.startedAt) {
		ledger.pending.startedAt = snapshot.startedAt
	}
	if snapshot.endedAt.After(ledger.pending.endedAt) {
		ledger.pending.endedAt = snapshot.endedAt
	}
	ledger.pending.counts.QueueFull += snapshot.counts.QueueFull
	ledger.pending.counts.DatabaseError += snapshot.counts.DatabaseError
	ledger.pending.counts.SerializationError += snapshot.counts.SerializationError
}

func (ledger *lossLedger) pendingCounts() LossCounts {
	ledger.mu.Lock()
	defer ledger.mu.Unlock()
	return ledger.pending.counts
}

func (ledger *lossLedger) totalCounts() LossCounts {
	ledger.mu.Lock()
	defer ledger.mu.Unlock()
	return ledger.total
}

func addLoss(counts *LossCounts, kind lossKind, count uint64) {
	switch kind {
	case lossQueue:
		counts.QueueFull += count
	case lossDatabase:
		counts.DatabaseError += count
	case lossSerialization:
		counts.SerializationError += count
	}
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

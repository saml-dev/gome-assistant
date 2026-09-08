package metrics

import (
	"encoding"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"crypto/rand"
)

type EventType string

const (
	EventSessionStarting            EventType = "session_starting"
	EventSessionReady               EventType = "session_ready"
	EventSessionFinished            EventType = "session_finished"
	EventMetricsDegraded            EventType = "metrics_degraded"
	EventMetricsEventsDropped       EventType = "metrics_events_dropped"
	EventMetricsRecovered           EventType = "metrics_recovered"
	EventAutomationTriggered        EventType = "automation_triggered"
	EventAutomationSkipped          EventType = "automation_skipped"
	EventAutomationWaiting          EventType = "automation_waiting"
	EventAutomationWaitCanceled     EventType = "automation_wait_canceled"
	EventAutomationCallbackStarted  EventType = "automation_callback_started"
	EventAutomationCallbackFinished EventType = "automation_callback_finished"
	EventStateReadStarted           EventType = "state_read_started"
	EventStateReadFinished          EventType = "state_read_finished"
	EventServiceCallStarted         EventType = "service_call_started"
	EventServiceCallFinished        EventType = "service_call_finished"
)

var eventTypes = map[EventType]struct{}{
	EventSessionStarting: {}, EventSessionReady: {}, EventSessionFinished: {},
	EventMetricsDegraded: {}, EventMetricsEventsDropped: {}, EventMetricsRecovered: {},
	EventAutomationTriggered: {}, EventAutomationSkipped: {}, EventAutomationWaiting: {},
	EventAutomationWaitCanceled: {}, EventAutomationCallbackStarted: {}, EventAutomationCallbackFinished: {},
	EventStateReadStarted: {}, EventStateReadFinished: {}, EventServiceCallStarted: {}, EventServiceCallFinished: {},
}

type Error struct {
	GoType  string `json:"go_type"`
	Message string `json:"message"`
}

func ErrorValue(err error) *Error {
	if err == nil {
		return nil
	}
	return &Error{GoType: reflect.TypeOf(err).String(), Message: err.Error()}
}

type Unavailable struct {
	TelemetryUnavailable bool   `json:"_telemetry_unavailable"`
	GoType               string `json:"go_type"`
	Error                string `json:"error"`
}

var ErrNonStringMapKey = errors.New("telemetry values may not contain maps with non-string keys")

// EncodeValue captures an arbitrary callback-visible value immediately. Values
// that JSON cannot represent become the versioned Unavailable object. In
// particular, maps with non-string keys are rejected instead of allowing
// encoding/json to silently stringify their keys.
func EncodeValue(value any) json.RawMessage {
	data, err := encodeJSONValue(value)
	if err == nil {
		return data
	}

	unavailable := Unavailable{
		TelemetryUnavailable: true,
		GoType:               goType(value),
		Error:                err.Error(),
	}
	data, marshalErr := marshalJSON(unavailable)
	if marshalErr != nil {
		// All fields above are strings/bools, so this is a defensive last resort.
		return json.RawMessage(`{"_telemetry_unavailable":true,"go_type":"unknown","error":"failed to encode unavailable telemetry value"}`)
	}
	return data
}

// EncodePayload synchronously snapshots a version-1 event payload. Payloads
// must be JSON objects. Use EncodeValue for callback-controlled subvalues.
func EncodePayload(payload any) (json.RawMessage, error) {
	if err := validateMapKeys(reflect.ValueOf(payload), make(map[visit]bool), 0); err != nil {
		return nil, err
	}
	data, err := marshalJSON(payload)
	if err != nil {
		return nil, err
	}
	trimmed := strings.TrimSpace(string(data))
	if len(trimmed) < 2 || trimmed[0] != '{' || trimmed[len(trimmed)-1] != '}' {
		return nil, errors.New("telemetry payload must encode as a JSON object")
	}
	return append(json.RawMessage(nil), data...), nil
}

func encodeJSONValue(value any) (json.RawMessage, error) {
	if err := validateMapKeys(reflect.ValueOf(value), make(map[visit]bool), 0); err != nil {
		return nil, err
	}
	data, err := marshalJSON(value)
	if err != nil {
		return nil, err
	}
	return append(json.RawMessage(nil), data...), nil
}

func marshalJSON(value any) (data []byte, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			data = nil
			err = fmt.Errorf("panic encoding telemetry JSON: %v", recovered)
		}
	}()
	return json.Marshal(value)
}

type visit struct {
	typeOf  reflect.Type
	pointer uintptr
}

var (
	jsonMarshalerType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

func validateMapKeys(value reflect.Value, visiting map[visit]bool, depth int) error {
	if !value.IsValid() || depth > 1000 {
		if depth > 1000 {
			return errors.New("telemetry value exceeds maximum nesting depth")
		}
		return nil
	}

	if value.Kind() == reflect.Interface {
		if value.IsNil() {
			return nil
		}
		return validateMapKeys(value.Elem(), visiting, depth+1)
	}

	typeOf := value.Type()
	if typeOf.Implements(jsonMarshalerType) || typeOf.Implements(textMarshalerType) ||
		(value.CanAddr() && (value.Addr().Type().Implements(jsonMarshalerType) || value.Addr().Type().Implements(textMarshalerType))) {
		return nil
	}

	var marker visit
	switch value.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice:
		if value.IsNil() {
			return nil
		}
		marker = visit{typeOf: typeOf, pointer: value.Pointer()}
		if visiting[marker] {
			return errors.New("telemetry value contains a cycle")
		}
		visiting[marker] = true
		defer delete(visiting, marker)
	}

	switch value.Kind() {
	case reflect.Pointer:
		return validateMapKeys(value.Elem(), visiting, depth+1)
	case reflect.Map:
		if value.Type().Key().Kind() != reflect.String {
			return fmt.Errorf("%w: %s", ErrNonStringMapKey, value.Type())
		}
		iterator := value.MapRange()
		for iterator.Next() {
			if err := validateMapKeys(iterator.Value(), visiting, depth+1); err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		for index := 0; index < value.Len(); index++ {
			if err := validateMapKeys(value.Index(index), visiting, depth+1); err != nil {
				return err
			}
		}
	case reflect.Struct:
		for index := 0; index < value.NumField(); index++ {
			field := value.Type().Field(index)
			if field.PkgPath != "" || field.Tag.Get("json") == "-" {
				continue
			}
			if err := validateMapKeys(value.Field(index), visiting, depth+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func goType(value any) string {
	if value == nil {
		return "<nil>"
	}
	return reflect.TypeOf(value).String()
}

// NewID returns a cryptographically random 128-bit identifier encoded as 32
// lowercase hexadecimal characters.
func NewID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate telemetry identifier: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

type EntityData struct {
	TriggerEntityID string          `json:"trigger_entity_id"`
	FromState       string          `json:"from_state"`
	FromAttributes  json.RawMessage `json:"from_attributes"`
	ToState         string          `json:"to_state"`
	ToAttributes    json.RawMessage `json:"to_attributes"`
	LastChanged     string          `json:"last_changed"`
}

type EventTrigger struct {
	Type         string `json:"type"`
	RawEventJSON string `json:"raw_event_json"`
}

type Trigger struct {
	Kind   string        `json:"kind"`
	Entity *EntityData   `json:"entity,omitempty"`
	Event  *EventTrigger `json:"event,omitempty"`
}

type Condition struct {
	EntityID          string `json:"entity_id"`
	ExpectedState     string `json:"expected_state"`
	RunOnNetworkError bool   `json:"run_on_network_error"`
}

type EmptyPayload struct{}

type SessionStartingPayload struct {
	WebUIEnabled bool `json:"web_ui_enabled"`
}

type SessionFinishedPayload struct {
	Reason string `json:"reason"`
	Error  *Error `json:"error"`
}

type MetricsDegradedPayload struct {
	Error Error `json:"error"`
}

type LossCounts struct {
	QueueFull          uint64 `json:"queue_full"`
	DatabaseError      uint64 `json:"database_error"`
	SerializationError uint64 `json:"serialization_error"`
}

func (counts LossCounts) total() uint64 {
	return counts.QueueFull + counts.DatabaseError + counts.SerializationError
}

type MetricsEventsDroppedPayload struct {
	PeriodStartedAt int64      `json:"period_started_at"`
	PeriodEndedAt   int64      `json:"period_ended_at"`
	Counts          LossCounts `json:"counts"`
}

type MetricsRecoveredPayload struct {
	DegradedOccurredAt int64 `json:"degraded_occurred_at"`
}

type AutomationTriggeredPayload struct {
	Trigger Trigger `json:"trigger"`
}

type AutomationSkippedPayload struct {
	Trigger   Trigger    `json:"trigger"`
	Reason    string     `json:"reason"`
	Condition *Condition `json:"condition"`
	Error     *Error     `json:"error"`
}

type AutomationWaitingPayload struct {
	Trigger     Trigger `json:"trigger"`
	DurationNS  int64   `json:"duration_ns"`
	TargetState string  `json:"target_state"`
}

type AutomationWaitCanceledPayload struct {
	Reason            string   `json:"reason"`
	CauseEvaluationID *string  `json:"cause_evaluation_id"`
	CauseTrigger      *Trigger `json:"cause_trigger"`
}

type AutomationCallbackFinishedPayload struct {
	DurationNS int64 `json:"duration_ns"`
}

type StateReadStartedPayload struct {
	Method  string          `json:"method"`
	Request json.RawMessage `json:"request"`
}

type StateReadFinishedPayload struct {
	Method string          `json:"method"`
	Result json.RawMessage `json:"result"`
	Error  *Error          `json:"error"`
}

type ServiceTarget struct {
	EntityIDs []string `json:"entity_ids"`
	AreaIDs   []string `json:"area_ids"`
	DeviceIDs []string `json:"device_ids"`
}

type ServiceCallRequest struct {
	Domain      string          `json:"domain"`
	Service     string          `json:"service"`
	ServiceData json.RawMessage `json:"service_data"`
	Target      ServiceTarget   `json:"target"`
}

type FireEventRequest struct {
	EventType string          `json:"event_type"`
	EventData json.RawMessage `json:"event_data"`
}

type ServiceCallStartedPayload struct {
	Kind    string              `json:"kind"`
	Request *ServiceCallRequest `json:"request,omitempty"`
	Event   *FireEventRequest   `json:"-"`
}

func (payload ServiceCallStartedPayload) MarshalJSON() ([]byte, error) {
	switch payload.Kind {
	case "call":
		return json.Marshal(struct {
			Kind    string             `json:"kind"`
			Request ServiceCallRequest `json:"request"`
		}{Kind: payload.Kind, Request: dereference(payload.Request)})
	case "fire_event":
		return json.Marshal(struct {
			Kind    string           `json:"kind"`
			Request FireEventRequest `json:"request"`
		}{Kind: payload.Kind, Request: dereference(payload.Event)})
	default:
		return nil, fmt.Errorf("unknown service operation kind %q", payload.Kind)
	}
}

func dereference[T any](value *T) T {
	if value == nil {
		var zero T
		return zero
	}
	return *value
}

type ServiceCallFinishedPayload struct {
	Kind   string          `json:"kind"`
	Result json.RawMessage `json:"result"`
	Error  *Error          `json:"error"`
}

func unixNano(value time.Time) int64 {
	return value.UTC().UnixNano()
}

func validEventType(eventType EventType) bool {
	_, ok := eventTypes[eventType]
	return ok
}

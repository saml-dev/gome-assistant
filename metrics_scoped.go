package gomeassistant

import (
	"context"
	"reflect"

	"saml.dev/gome-assistant/internal/metrics"
	"saml.dev/gome-assistant/internal/services"
)

// scopedState keeps an evaluation's identity on the State object itself. It
// deliberately has no goroutine-local state, so callers may pass it to their
// own child goroutines without losing attribution.
type scopedState struct {
	delegate   *StateImpl
	evaluation *evaluationTrace
}

func (state *scopedState) ListEntities() ([]EntityState, error) {
	operation := state.begin("list_entities", map[string]any{})
	result, err := state.delegate.ListEntities()
	operation.finish(result, err)
	return result, err
}

func (state *scopedState) Get(entityID string) (EntityState, error) {
	return state.getWithContext(context.Background(), entityID)
}

func (state *scopedState) getWithContext(ctx context.Context, entityID string) (EntityState, error) {
	operation := state.begin("get", map[string]any{"entity_id": entityID})
	result, err := state.delegate.getWithContext(ctx, entityID)
	operation.finish(result, err)
	return result, err
}

func (state *scopedState) GetStates(entityIDs []string) ([]string, error) {
	if len(entityIDs) == 0 {
		// StateImpl returns locally for an empty request, so there is no Home
		// Assistant read to trace.
		return state.delegate.GetStates(entityIDs)
	}
	operation := state.begin("get_states", map[string]any{"entity_ids": append([]string{}, entityIDs...)})
	result, err := state.delegate.GetStates(entityIDs)
	operation.finish(result, err)
	return result, err
}

func (state *scopedState) Equals(entityID string, expectedState string) (bool, error) {
	// Get is the one underlying Home Assistant read. Do not create another
	// operation around Equals itself.
	currentState, err := state.Get(entityID)
	if err != nil {
		return false, err
	}
	return currentState.State == expectedState, nil
}

func (state *scopedState) BeforeSunrise(offset ...DurationString) bool {
	return state.delegate.BeforeSunrise(offset...)
}

func (state *scopedState) AfterSunrise(offset ...DurationString) bool {
	return state.delegate.AfterSunrise(offset...)
}

func (state *scopedState) BeforeSunset(offset ...DurationString) bool {
	return state.delegate.BeforeSunset(offset...)
}

func (state *scopedState) AfterSunset(offset ...DurationString) bool {
	return state.delegate.AfterSunset(offset...)
}

type stateOperation struct {
	evaluation  *evaluationTrace
	method      string
	operationID string
}

func (state *scopedState) begin(method string, request any) stateOperation {
	operation := stateOperation{method: method}
	if state == nil || state.evaluation == nil || state.evaluation.recorder == nil {
		return operation
	}
	operation.evaluation = state.evaluation
	operation.operationID = operation.evaluation.newOperationID()
	if operation.operationID == "" {
		return operation
	}
	operation.evaluation.recordOperation(operation.operationID, metrics.EventStateReadStarted, metrics.StateReadStartedPayload{
		Method: method, Request: metrics.EncodeValue(request),
	})
	return operation
}

func (operation stateOperation) finish(result any, err error) {
	if operation.evaluation == nil || operation.operationID == "" {
		return
	}
	var captured any
	if err == nil {
		captured = result
	}
	operation.evaluation.recordOperation(operation.operationID, metrics.EventStateReadFinished, metrics.StateReadFinishedPayload{
		Method: operation.method, Result: metrics.EncodeValue(captured), Error: metrics.ErrorValue(err),
	})
}

// scopedServicesAPI is the narrow generated-service boundary. CallAndForget
// intentionally delegates without telemetry; it is outside metrics scope.
type scopedServicesAPI struct {
	delegate   services.API
	evaluation *evaluationTrace
}

func (api *scopedServicesAPI) CallAndForget(request services.BaseServiceRequest) error {
	return api.delegate.CallAndForget(request)
}

func (api *scopedServicesAPI) Call(ctx context.Context, request services.BaseServiceRequest, result any) error {
	operation := api.beginCall(request)
	err := api.delegate.Call(ctx, request, result)
	operation.finish(callbackVisibleResult(result, err), err)
	return err
}

func (api *scopedServicesAPI) FireEvent(eventType string, eventData map[string]any) error {
	operation := api.beginFireEvent(eventType, eventData)
	err := api.delegate.FireEvent(eventType, eventData)
	operation.finish(nil, err)
	return err
}

type serviceOperation struct {
	evaluation  *evaluationTrace
	operationID string
	kind        string
}

func (api *scopedServicesAPI) beginCall(request services.BaseServiceRequest) serviceOperation {
	operation := serviceOperation{kind: "call"}
	if api == nil || api.evaluation == nil || api.evaluation.recorder == nil {
		return operation
	}
	operation.evaluation = api.evaluation
	operation.operationID = operation.evaluation.newOperationID()
	if operation.operationID == "" {
		return operation
	}
	operation.evaluation.recordOperation(operation.operationID, metrics.EventServiceCallStarted, metrics.ServiceCallStartedPayload{
		Kind: "call",
		Request: &metrics.ServiceCallRequest{
			Domain: request.Domain, Service: request.Service,
			ServiceData: metrics.EncodeValue(request.ServiceData),
			Target: metrics.ServiceTarget{
				EntityIDs: append([]string{}, request.Target.EntityIDs...),
				AreaIDs:   append([]string{}, request.Target.AreaIDs...),
				DeviceIDs: append([]string{}, request.Target.DeviceIDs...),
			},
		},
	})
	return operation
}

func (api *scopedServicesAPI) beginFireEvent(eventType string, eventData map[string]any) serviceOperation {
	operation := serviceOperation{kind: "fire_event"}
	if api == nil || api.evaluation == nil || api.evaluation.recorder == nil {
		return operation
	}
	operation.evaluation = api.evaluation
	operation.operationID = operation.evaluation.newOperationID()
	if operation.operationID == "" {
		return operation
	}
	operation.evaluation.recordOperation(operation.operationID, metrics.EventServiceCallStarted, metrics.ServiceCallStartedPayload{
		Kind: "fire_event",
		Event: &metrics.FireEventRequest{
			EventType: eventType, EventData: metrics.EncodeValue(eventData),
		},
	})
	return operation
}

func (operation serviceOperation) finish(result any, err error) {
	if operation.evaluation == nil || operation.operationID == "" {
		return
	}
	operation.evaluation.recordOperation(operation.operationID, metrics.EventServiceCallFinished, metrics.ServiceCallFinishedPayload{
		Kind: operation.kind, Result: metrics.EncodeValue(result), Error: metrics.ErrorValue(err),
	})
}

func callbackVisibleResult(result any, err error) any {
	if err != nil || result == nil {
		return nil
	}
	value := reflect.ValueOf(result)
	for value.IsValid() && (value.Kind() == reflect.Interface || value.Kind() == reflect.Pointer) {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}
	if !value.IsValid() {
		return nil
	}
	return value.Interface()
}

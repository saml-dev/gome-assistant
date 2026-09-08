package gomeassistant

import (
	"context"
	"sync/atomic"
	"time"

	"saml.dev/gome-assistant/internal"
	"saml.dev/gome-assistant/internal/metrics"
)

type evaluationTrace struct {
	app           *App
	recorder      *metrics.Recorder
	sessionID     string
	automationKey string
	evaluationID  string
	trigger       metrics.Trigger
	waiting       atomic.Bool
}

func (app *App) newEvaluation(automationKey string, trigger metrics.Trigger) *evaluationTrace {
	evaluation := &evaluationTrace{app: app, automationKey: automationKey, trigger: trigger}
	origin := app.metricsOrigin.Load()
	if origin == nil || automationKey == "" {
		return evaluation
	}
	evaluationID, err := metrics.NewID()
	if err != nil {
		return evaluation
	}
	evaluation.recorder = origin.recorder
	evaluation.sessionID = origin.sessionID
	evaluation.evaluationID = evaluationID
	evaluation.record(metrics.EventAutomationTriggered, metrics.AutomationTriggeredPayload{Trigger: trigger})
	app.updateCatalogStatus(automationKey, evaluationID, "running", time.Now())
	return evaluation
}

func (evaluation *evaluationTrace) record(eventType metrics.EventType, payload any) bool {
	if evaluation == nil || evaluation.recorder == nil || evaluation.evaluationID == "" {
		return false
	}
	return evaluation.recorder.Record(metrics.Event{
		Type:          eventType,
		OccurredAt:    time.Now(),
		SessionID:     evaluation.sessionID,
		AutomationKey: evaluation.automationKey,
		EvaluationID:  evaluation.evaluationID,
		Payload:       payload,
	})
}

func (evaluation *evaluationTrace) newOperationID() string {
	if evaluation == nil || evaluation.recorder == nil || evaluation.evaluationID == "" {
		return ""
	}
	id, err := metrics.NewID()
	if err != nil {
		return ""
	}
	return id
}

func (evaluation *evaluationTrace) recordOperation(operationID string, eventType metrics.EventType, payload any) bool {
	if evaluation == nil || evaluation.recorder == nil || evaluation.evaluationID == "" || operationID == "" {
		return false
	}
	return evaluation.recorder.Record(metrics.Event{
		Type:          eventType,
		OccurredAt:    time.Now(),
		SessionID:     evaluation.sessionID,
		AutomationKey: evaluation.automationKey,
		EvaluationID:  evaluation.evaluationID,
		OperationID:   operationID,
		Payload:       payload,
	})
}

func (evaluation *evaluationTrace) skip(reason string, condition *internal.EnabledDisabledInfo, err error) {
	evaluation.record(metrics.EventAutomationSkipped, metrics.AutomationSkippedPayload{
		Trigger: evaluation.trigger, Reason: reason,
		Condition: metricsCondition(condition), Error: metrics.ErrorValue(err),
	})
	if evaluation.app != nil {
		evaluation.app.updateCatalogStatus(evaluation.automationKey, evaluation.evaluationID, "skipped", time.Now())
	}
}

func (evaluation *evaluationTrace) skipCondition(check conditionCheck, mismatchReason, queryErrorReason string) bool {
	if !check.fail {
		return false
	}
	reason := mismatchReason
	if check.err != nil {
		reason = queryErrorReason
	}
	evaluation.skip(reason, check.condition, check.err)
	return true
}

func metricsCondition(condition *internal.EnabledDisabledInfo) *metrics.Condition {
	if condition == nil {
		return nil
	}
	return &metrics.Condition{
		EntityID: condition.Entity, ExpectedState: condition.State,
		RunOnNetworkError: condition.RunOnError,
	}
}

func (evaluation *evaluationTrace) scopedState(state *StateImpl) State {
	if evaluation == nil || evaluation.recorder == nil {
		return state
	}
	return &scopedState{delegate: state, evaluation: evaluation}
}

func (evaluation *evaluationTrace) getWithContext(ctx context.Context, state *StateImpl, entityID string) (EntityState, error) {
	if scoped, ok := evaluation.scopedState(state).(*scopedState); ok {
		return scoped.getWithContext(ctx, entityID)
	}
	return state.getWithContext(ctx, entityID)
}

func (evaluation *evaluationTrace) scopedService(app *App) *Service {
	if evaluation == nil || evaluation.recorder == nil {
		return app.service
	}
	return newService(&scopedServicesAPI{delegate: app, evaluation: evaluation})
}

func (evaluation *evaluationTrace) runCallback(app *App, callback func(*Service, State)) bool {
	if evaluation == nil || evaluation.recorder == nil {
		return app.goTracked(func() { callback(app.service, app.state) })
	}
	state := evaluation.scopedState(app.state)
	service := evaluation.scopedService(app)
	return app.goTracked(func() {
		started := time.Now()
		evaluation.record(metrics.EventAutomationCallbackStarted, metrics.EmptyPayload{})
		app.changeRunning(1)
		defer app.changeRunning(-1)
		callback(service, state)
		evaluation.record(metrics.EventAutomationCallbackFinished, metrics.AutomationCallbackFinishedPayload{
			DurationNS: time.Since(started).Nanoseconds(),
		})
		app.updateCatalogStatus(evaluation.automationKey, evaluation.evaluationID, "completed", time.Now())
	})
}

func (evaluation *evaluationTrace) runEntityCallback(app *App, callback EntityListenerCallback, data EntityData) bool {
	return evaluation.runCallback(app, func(service *Service, state State) { callback(service, state, data) })
}

func (evaluation *evaluationTrace) runEventCallback(app *App, callback EventListenerCallback, data EventData) bool {
	return evaluation.runCallback(app, func(service *Service, state State) { callback(service, state, data) })
}

func (evaluation *evaluationTrace) beginWaiting(duration time.Duration, targetState string) {
	if evaluation.waiting.CompareAndSwap(false, true) {
		evaluation.record(metrics.EventAutomationWaiting, metrics.AutomationWaitingPayload{
			Trigger: evaluation.trigger, DurationNS: duration.Nanoseconds(), TargetState: targetState,
		})
		if evaluation.app != nil && evaluation.recorder != nil && evaluation.evaluationID != "" {
			evaluation.app.changeWaiting(1)
			evaluation.app.updateCatalogStatus(evaluation.automationKey, evaluation.evaluationID, "waiting", time.Now())
		}
	}
}

func (evaluation *evaluationTrace) finishWaiting() {
	if evaluation.waiting.CompareAndSwap(true, false) && evaluation.app != nil && evaluation.recorder != nil && evaluation.evaluationID != "" {
		evaluation.app.changeWaiting(-1)
	}
}

func (evaluation *evaluationTrace) cancelWaiting(reason string, cause *evaluationTrace) {
	if !evaluation.waiting.CompareAndSwap(true, false) {
		return
	}
	var causeID *string
	var causeTrigger *metrics.Trigger
	if cause != nil && cause.evaluationID != "" {
		id := cause.evaluationID
		trigger := cause.trigger
		causeID, causeTrigger = &id, &trigger
	}
	evaluation.record(metrics.EventAutomationWaitCanceled, metrics.AutomationWaitCanceledPayload{
		Reason: reason, CauseEvaluationID: causeID, CauseTrigger: causeTrigger,
	})
	if evaluation.app != nil && evaluation.recorder != nil && evaluation.evaluationID != "" {
		evaluation.app.changeWaiting(-1)
		evaluation.app.updateCatalogStatus(evaluation.automationKey, evaluation.evaluationID, "canceled", time.Now())
	}
}

func dailyTrigger() metrics.Trigger    { return metrics.Trigger{Kind: "daily"} }
func intervalTrigger() metrics.Trigger { return metrics.Trigger{Kind: "interval"} }

func entityTrigger(kind string, data EntityData) metrics.Trigger {
	return metrics.Trigger{Kind: kind, Entity: &metrics.EntityData{
		TriggerEntityID: data.TriggerEntityID,
		FromState:       data.FromState, FromAttributes: metrics.EncodeValue(data.FromAttributes),
		ToState: data.ToState, ToAttributes: metrics.EncodeValue(data.ToAttributes),
		LastChanged: data.LastChanged.Format(time.RFC3339Nano),
	}}
}

func startupEntityTrigger(entityID string) metrics.Trigger {
	return entityTrigger("entity_startup", EntityData{TriggerEntityID: entityID})
}

func eventTrigger(data EventData) metrics.Trigger {
	return metrics.Trigger{Kind: "event", Event: &metrics.EventTrigger{
		Type: data.Type, RawEventJSON: string(data.RawEventJSON),
	}}
}

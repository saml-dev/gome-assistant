package gomeassistant

import (
	"errors"
	"time"
)

type entityListener struct {
	entityIds    []string
	callback     entityListenerCallback
	fromState    string
	toState      string
	betweenStart time.Duration
	betweenEnd   time.Duration
	err          error
}

type entityListenerCallback func(*Service, *EntityData)

// TODO: use this to flatten json sent from HA for trigger event
type EntityData struct {
	TriggerEntityId string
	FromState       string
	FromAttributes  map[string]any
	ToState         string
	ToAttributes    map[string]any
	LastChanged     time.Time
}

type triggerMsg struct {
	Id    int64  `json:"id"`
	Type  string `json:"type"`
	Event struct {
		Variables struct {
			Trigger struct {
				EntityId  string          `json:"entity_id"`
				FromState triggerMsgState `json:"from_state"`
				ToState   triggerMsgState `json:"to_state"`
			}
		} `json:"variables"`
	} `json:"event"`
}

type triggerMsgState struct {
	State       string         `json:"state"`
	Attributes  map[string]any `json:"attributes"`
	LastChanged string         `json:"last_changed"`
}

type subscribeMsg struct {
	Id      int64               `json:"id"`
	Type    string              `json:"type"`
	Trigger subscribeMsgTrigger `json:"trigger"`
}

type subscribeMsgTrigger struct {
	Platform string `json:"platform"`
	EntityId string `json:"entity_id"`
	From     string `json:"from"`
	To       string `json:"to"`
}

/* Builders */

func EntityListenerBuilder() elBuilder1 {
	return elBuilder1{entityListener{}}
}

type elBuilder1 struct {
	entityListener
}

func (b elBuilder1) EntityIds(entityIds ...string) elBuilder2 {
	if len(entityIds) == 0 {
		b.err = errors.New("must pass at least one entityId to EntityIds()")
	} else {
		b.entityListener.entityIds = entityIds
	}
	return elBuilder2(b)
}

type elBuilder2 struct {
	entityListener
}

func (b elBuilder2) Call(callback entityListenerCallback) elBuilder3 {
	if b.err == nil {
		b.entityListener.callback = callback
	}
	return elBuilder3(b)
}

type elBuilder3 struct {
	entityListener
}

func (b elBuilder3) OnlyBetween(start time.Duration, end time.Duration) elBuilder3 {
	b.entityListener.betweenStart = start
	b.entityListener.betweenEnd = end
	return b
}

func (b elBuilder3) FromState(s string) elBuilder3 {
	b.entityListener.fromState = s
	return b
}

func (b elBuilder3) ToState(s string) elBuilder3 {
	b.entityListener.toState = s
	return b
}

func (b elBuilder3) Build() entityListener {
	return b.entityListener
}

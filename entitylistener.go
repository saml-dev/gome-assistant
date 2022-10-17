package gomeassistant

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-module/carbon"
	i "github.com/saml-dev/gome-assistant/internal"
)

type entityListener struct {
	entityIds    []string
	callback     entityListenerCallback
	fromState    string
	toState      string
	betweenStart string
	betweenEnd   string
	err          error
}

type entityListenerCallback func(*Service, EntityData)

// TODO: use this to flatten json sent from HA for trigger event
type EntityData struct {
	TriggerEntityId string
	FromState       string
	FromAttributes  map[string]any
	ToState         string
	ToAttributes    map[string]any
	LastChanged     time.Time
}

type stateChangedMsg struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Event struct {
		Data struct {
			EntityID string   `json:"entity_id"`
			NewState msgState `json:"new_state"`
			OldState msgState `json:"old_state"`
		} `json:"data"`
		EventType string `json:"event_type"`
		Origin    string `json:"origin"`
	} `json:"event"`
}

type msgState struct {
	EntityID    string         `json:"entity_id"`
	LastChanged time.Time      `json:"last_changed"`
	State       string         `json:"state"`
	Attributes  map[string]any `json:"attributes"`
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

/* Methods */

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

func (b elBuilder3) OnlyBetween(start string, end string) elBuilder3 {
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

/* Functions */
func callEntityListeners(app *app, msgBytes []byte) {
	msg := stateChangedMsg{}
	json.Unmarshal(msgBytes, &msg)
	data := msg.Event.Data
	eid := data.EntityID
	listeners, ok := app.entityListeners[eid]
	if !ok {
		// no listeners registered for this id
		return
	}

	for _, l := range listeners {
		// if betweenStart and betweenEnd both set, first account for midnight
		// overlap, then only run if between those times.
		if l.betweenStart != "" && l.betweenEnd != "" {
			start := i.ParseTime(l.betweenStart)
			end := i.ParseTime(l.betweenEnd)

			// check for midnight overlap
			if end.Lt(start) { // example turn on night lights when motion from 23:00 to 07:00
				if end.IsPast() { // such as at 15:00, 22:00
					end = end.AddDay()
				} else {
					start = start.SubDay() // such as at 03:00, 05:00
				}
			}

			// skip callback if not inside the range
			if !carbon.Now().BetweenIncludedStart(start, end) {
				return
			}
		}
		// otherwise, just check if before/after the individual times
		if l.betweenStart != "" && i.ParseTime(l.betweenStart).IsFuture() {
			return
		}
		if l.betweenEnd != "" && i.ParseTime(l.betweenEnd).IsPast() {
			return
		}

		// don't run callback if fromState or toState are set and don't match
		if l.fromState != "" && l.fromState != data.OldState.State {
			return
		}
		if l.toState != "" && l.toState != data.NewState.State {
			return
		}

		entityData := EntityData{
			TriggerEntityId: eid,
			FromState:       data.OldState.State,
			FromAttributes:  data.OldState.Attributes,
			ToState:         data.NewState.State,
			ToAttributes:    data.NewState.Attributes,
			LastChanged:     data.OldState.LastChanged,
		}
		l.callback(app.service, entityData)
	}
}

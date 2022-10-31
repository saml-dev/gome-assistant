package gomeassistant

import (
	"encoding/json"
	"time"

	"github.com/golang-module/carbon"
	"github.com/saml-dev/gome-assistant/internal"
)

type EntityListener struct {
	entityIds    []string
	callback     EntityListenerCallback
	fromState    string
	toState      string
	betweenStart string
	betweenEnd   string
	throttle     time.Duration
	lastRan      carbon.Carbon
	delay        time.Duration
	delayTimer   *time.Timer
}

type EntityListenerCallback func(*Service, EntityData)

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

/* Methods */

func EntityListenerBuilder() elBuilder1 {
	return elBuilder1{EntityListener{
		lastRan: carbon.Now().StartOfCentury(),
	}}
}

type elBuilder1 struct {
	entityListener EntityListener
}

func (b elBuilder1) EntityIds(entityIds ...string) elBuilder2 {
	if len(entityIds) == 0 {
		panic("must pass at least one entityId to EntityIds()")
	} else {
		b.entityListener.entityIds = entityIds
	}
	return elBuilder2(b)
}

type elBuilder2 struct {
	entityListener EntityListener
}

func (b elBuilder2) Call(callback EntityListenerCallback) elBuilder3 {
	b.entityListener.callback = callback
	return elBuilder3(b)
}

type elBuilder3 struct {
	entityListener EntityListener
}

func (b elBuilder3) OnlyBetween(start string, end string) elBuilder3 {
	b.entityListener.betweenStart = start
	b.entityListener.betweenEnd = end
	return b
}

func (b elBuilder3) OnlyAfter(start string) elBuilder3 {
	b.entityListener.betweenStart = start
	return b
}

func (b elBuilder3) OnlyBefore(end string) elBuilder3 {
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

func (b elBuilder3) Duration(s DurationString) elBuilder3 {
	d := internal.ParseDuration(string(s))
	b.entityListener.delay = d
	return b
}

func (b elBuilder3) Throttle(s DurationString) elBuilder3 {
	d := internal.ParseDuration(string(s))
	b.entityListener.throttle = d
	return b
}

func (b elBuilder3) Build() EntityListener {
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
		// Check conditions
		if c := checkWithinTimeRange(l.betweenStart, l.betweenEnd); c.fail {
			continue
		}
		if c := checkStatesMatch(l.fromState, data.OldState.State); c.fail {
			continue
		}
		if c := checkStatesMatch(l.toState, data.NewState.State); c.fail {
			if l.delayTimer != nil {
				l.delayTimer.Stop()
			}
			continue
		}
		if c := checkThrottle(l.throttle, l.lastRan); c.fail {
			continue
		}

		entityData := EntityData{
			TriggerEntityId: eid,
			FromState:       data.OldState.State,
			FromAttributes:  data.OldState.Attributes,
			ToState:         data.NewState.State,
			ToAttributes:    data.NewState.Attributes,
			LastChanged:     data.OldState.LastChanged,
		}

		if l.delay != 0 {
			l.delayTimer = time.AfterFunc(l.delay, func() {
				go l.callback(app.service, entityData)
				l.lastRan = carbon.Now()
			})
			return
		}

		// run now if no delay set
		go l.callback(app.service, entityData)
		l.lastRan = carbon.Now()
	}
}

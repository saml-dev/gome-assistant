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

type entityListenerCallback func(*Service, *Data)

type Data struct{}

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

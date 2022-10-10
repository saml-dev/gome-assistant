package gomeassistant

type entityListener struct {
	entityId     string
	callback     entityListenerCallback
	fromState    string
	toState      string
	betweenStart hourMinute
	betweenEnd   hourMinute
}

type entityListenerCallback func(Service, Data)

type Data struct{}

func (b elBuilder3) OnlyBetween(start hourMinute, end hourMinute) elBuilder3 {
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

func (b elBuilder1) EntityId(eid string) elBuilder2 {
	b.entityListener.entityId = eid
	return elBuilder2(b)
}

type elBuilder2 struct {
	entityListener
}

func (b elBuilder2) Call(callback entityListenerCallback) elBuilder3 {
	b.entityListener.callback = callback
	return elBuilder3(b)
}

type elBuilder3 struct {
	entityListener
}

# Gome-Assistant

Write your [Home Assistant](https://www.home-assistant.io/) automations with a strongly-typed Golang library!

## Disclaimer

Gome-Assistant is a new library, and I'm opening it up early to get some user feedback on the API and help shape the direction. I plan for it to grow to cover all Home Assistant use cases, services, and event types. So it's possible that breaking changes will happen before v1.0.0!

## Quick Start

### Installation

```
go get github.com/saml-dev/gome-assistant
```

### Write your automations

Check out `example/example.go` for an example of the 3 types of automations — schedules, entity listeners, and event listeners.

> ℹ️ Instead of copying and pasting, try typing it yourself to see how autocomplete guides you through the setup using a builder pattern.

### Run your code

Keeping with the simplicity that Go is famous for, you don't need a specific environment or docker container to run Gome-Assistant. You just write your code like any other Go binary you would write. So once you have your automations, you can run it however you like — using `screen` or `tmux`, a cron job, or wrap it up in a docker container if you just can't get enough docker!

> _❗ No promises, but I may provide a Docker image with file watching to automatically restart gome-assistant, to make it easier to use gome-assistant on a fully managed Home Assistant installation._

## gome-assistant Concepts

### Overview

The general flow is

1. Create your app
2. Register automations
3. Start app

```go
import ga "github.com/saml-dev/gome-assistant"

// replace with IP and port of your Home Assistant installation
app := ga.NewApp("0.0.0.0:8123")

// create automations here (see next sections)

// register automations
app.RegisterSchedules(...)
app.RegisterEntityListeners(...)
app.RegisterEventListeners(...)

app.Start()
```

A full reference is available on [pkg.go.dev](https://pkg.go.dev/github.com/saml-dev/gome-assistant), but all you need to know to get started are the three types of automations in gome-assistant.

### Schedules

Schedules are simply a way to run a function on a schedule. The most common schedule is once a day at a certain time.

```go
_7pm := ga.NewSchedule().Call(myFunc).Daily().At("19:00").Build()
```

Schedules can also be run at sunrise or sunset, with an optional [offset](https://pkg.go.dev/time#ParseDuration).

```go
// 30 mins before sunrise
sunrise := ga.NewSchedule().Call(myFunc).Daily().Sunrise(app, "-30m").Build()
```

Schedules are also used to run a function on an interval. Offset is used to offset the first run of a schedule from midnight.

```go
// run every hour at the 30-minute mark
interval := ga.NewSchedule().Call(a).Every("1h").Offset("30m").Build()
```

#### Schedule Callback function

The function passed to `.Call()` must take

- `*ga.Service` used to call home assistant services
- `*ga.State` used to retrieve state from home assistant

```go
func myFunc(se *ga.Service, st *ga.State) {
  ...
}
```

### Entity Listeners

Entity Listeners are used to respond to entities changing state. The simplest entity listener looks like:

```go
etl := ga.NewEntityListener().EntityIds("binary_sensor.front_door").Call(myFunc).Build()
```

Entity listeners have other functions to change the behavior.

| Function                             | Info                                                                                              |
| ------------------------------------ | ------------------------------------------------------------------------------------------------- |
| ToState("on")                        | Function only called if new state matches argument.                                               |
| FromState("on")                      | Function only called if old state matches argument.                                               |
| Throttle("30s")                      | Minimum time between function calls.                                                              |
| Duration("30s")                      | Requires ToState(). Sets how long the entity should be in the state before running your function. |
| OnlyAfter("03:00")                   | Only run your function after a specified time of day.                                             |
| OnlyBefore("03:00")                  | Only run your function before a specified time of day.                                            |
| OnlyBetween("03:00", "14:00")        | Only run your function between two specified times of day.                                        |
| ExceptionDay(time.Time)              | A one time exception on the given date. Time is ignored, applies to whole day.                    |
| ExceptionRange(time.Time, time.Time) | A one time exception between the two date/times. Both date and time are considered.               |
| RunOnStartup()                       | Run your callback once during App.Start()                                                         |

#### Entity Listener Callback function

The function passed to `.Call()` must take

- `*ga.Service` used to call home assistant services
- `*ga.State` used to retrieve state from home assistant
- `ga.EntityData` which is the entity that triggered the listener

```go
func myFunc(se *ga.Service, st *ga.State, e ga.EntityData) {
  ...
}
```

### Event Listeners

Event Listeners are used to respond to entities changing state. The simplest event listener looks like:

```go
evl := ga.NewEntityListener().EntityIds("binary_sensor.front_door").Call(myFunc).Build()
```

Event listeners have other functions to change the behavior.

| Function                             | Info                                                                                |
| ------------------------------------ | ----------------------------------------------------------------------------------- |
| Throttle("30s")                      | Minimum time between function calls.                                                |
| OnlyAfter("03:00")                   | Only run your function after a specified time of day.                               |
| OnlyBefore("03:00")                  | Only run your function before a specified time of day.                              |
| OnlyBetween("03:00", "14:00")        | Only run your function between two specified times of day.                          |
| ExceptionDay(time.Time)              | A one time exception on the given date. Time is ignored, applies to whole day.      |
| ExceptionRange(time.Time, time.Time) | A one time exception between the two date/times. Both date and time are considered. |

#### Event Listener Callback function

The function passed to `.Call()` must take

- `*ga.Service` used to call home assistant services
- `*ga.State` used to retrieve state from home assistant
- `ga.EventData` containing the event data that triggered the listener

```go
func myFunc(se *ga.Service, st *ga.State) {
  ...
}
```

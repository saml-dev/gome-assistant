# Gome-Assistant

Write strongly typed [Home Assistant](https://www.home-assistant.io/) automations in Go!

## Disclaimer

Gome-Assistant is a new library, and I'm opening it up early to get some user feedback on the API and help shape the direction. I plan for it to grow to cover all Home Assistant use cases, services, and event types. So it's possible that breaking changes will happen before v1.0.0!

## Quick Start

### Installation

```
go get saml.dev/gome-assistant
```

### Generate Entity Constants

You can generate type-safe constants for all your Home Assistant entities using `go generate`. This makes it easier to reference entities in your code.

1. Create a `gen.yaml` file in your project root:

```yaml
url: "http://192.168.1.123:8123"  
ha_auth_token: "your_auth_token"  # Or set HA_AUTH_TOKEN env var
home_zone_entity_id: "zone.home"  # Optional: defaults to zone.home

# Optional: List of domains to include when generating constants
# If provided, only these domains will be processed
include_domains: ["light", "switch", "climate"]

# Optional: List of domains to exclude when generating constants
# Only used if include_domains is empty
exclude_domains: ["device_tracker", "person"]
```

2. Add a `//go:generate` comment in your project:

```go
//go:generate go run saml.dev/gome-assistant/cmd/generate
```

Optionally use the `-config` flag to customize the file path of the config file.

3. Run the generator:

```
go generate
```

This will create an `entities` package with type-safe constants for all your Home Assistant entities, organized by domain. For example:

```go
import "your_project/entities"

// Instead of writing "light.living_room" as a string:
entities.Light.LivingRoom // Type-safe constant

// All your entities are organized by domain
entities.Switch.Kitchen
entities.Climate.Bedroom
entities.MediaPlayer.TVRoom
```

The constants are based on the entity ID itself, not the name of the entity in Home Assistant.

### Write your automations

Check out [`example/example.go`](./example/example.go) for an example of the 3 types of automations — schedules, entity listeners, and event listeners.

> ℹ️ Instead of copying and pasting, try typing it yourself to see how autocomplete guides you through the setup using a builder pattern.

### Run your code

Keeping with the simplicity that Go is famous for, you don't need a specific environment or docker container to run Gome-Assistant. You just write and run your code like any other Go binary. So once you build your code, you can run it however you like — using `screen` or `tmux`, a cron job, a linux service, or wrap it up in a docker container if you like!

> _❗ No promises, but I may provide a Docker image with file watching to automatically restart gome-assistant, to make it easier to use gome-assistant on a fully managed Home Assistant installation._

## gome-assistant Concepts

### Overview

The general flow is

1. Create your app
2. Register automations
3. Start app

```go
import ga "saml.dev/gome-assistant"

// replace with IP and port of your Home Assistant installation
app, err := ga.NewApp(ga.NewAppRequest{
	URL:              "http://192.168.1.123:8123",
	HAAuthToken:      os.Getenv("HA_AUTH_TOKEN"),
	HomeZoneEntityId: "zone.home",
})

// create automations here (see next sections)

// register automations
app.RegisterSchedules(...)
app.RegisterEntityListeners(...)
app.RegisterEventListeners(...)
app.RegisterIntervals(...)

app.Start()
```

A full reference is available on [pkg.go.dev](https://pkg.go.dev/saml.dev/gome-assistant), but all you need to know to get started are the four types of automations in gome-assistant.

- [Daily Schedules](#daily-schedule)
- [Entity Listeners](#entity-listener)
- [Event Listeners](#event-listener)
- [Intervals](#interval)

### Daily Schedule

Daily Schedules run at a specific time each day.

```go
_7pm := ga.NewDailySchedule().Call(myFunc).At("19:00").Build()
```

Schedules can also be run at sunrise or sunset, with an optional [offset](https://pkg.go.dev/time#ParseDuration).

```go
// 30 mins before sunrise
sunrise := ga.NewDailySchedule().Call(myFunc).Sunrise(app, "-30m").Build()
// at sunset
sunset := ga.NewDailySchedule().Call(myFunc).Sunset().Build()
```

Daily schedules have other functions to change the behavior.

| Function                                  | Info                                                                                                     |
| ----------------------------------------- | -------------------------------------------------------------------------------------------------------- |
| ExceptionDates(t time.Time, ...time.Time) | Skip the schedule on the given date(s). Functions like a blocklist. Cannot be combined with OnlyOnDates. |
| OnlyOnDates(t time.Time, ...time.Time)    | Run only on the given date(s). Functions like an allowlist. Cannot be combined with ExceptionDates.      |

#### Schedule Callback function

The function passed to `.Call()` must take

- `*ga.Service` used to call home assistant services
- `*ga.State` used to retrieve state from home assistant

```go
func myFunc(se *ga.Service, st *ga.State) {
  // ...
}
```

### Entity Listener

Entity Listeners are used to respond to entities changing state. The simplest entity listener looks like:

```go
etl := ga.NewEntityListener().EntityIds("binary_sensor.front_door").Call(myFunc).Build()
```

Entity listeners have other functions to change the behavior.

| Function                                | Info                                                                                                              |
| --------------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| ToState("on")                           | Function only called if new state matches argument.                                                               |
| FromState("on")                         | Function only called if old state matches argument.                                                               |
| Throttle("30s")                         | Minimum time between function calls.                                                                              |
| Duration("30s")                         | Requires ToState(). Sets how long the entity must be in the state before running your function.                   |
| OnlyAfter("03:00")                      | Only run your function after a specified time of day.                                                             |
| OnlyBefore("03:00")                     | Only run your function before a specified time of day.                                                            |
| OnlyBetween("03:00", "14:00")           | Only run your function between two specified times of day.                                                        |
| ExceptionDates(time.Time, ...time.Time) | A one time exception on the given date. Time is ignored, applies to whole day. Functions like a "blocklist".      |
| ExceptionRange(time.Time, time.Time)    | A one time exception between the two date/times. Both date and time are considered. Functions like a "blocklist". |
| RunOnStartup()                          | Run your callback during `App.Start()`.                                                                           |

#### Entity Listener Callback function

The function passed to `.Call()` must take

- `*ga.Service` used to call home assistant services
- `*ga.State` used to retrieve state from home assistant
- `ga.EntityData` which is the entity that triggered the listener

```go
func myFunc(se *ga.Service, st *ga.State, e ga.EntityData) {
  // ...
}
```

### Event Listeners

Event listeners allow you to respond to Home Assistant events in real-time. You can create an event listener using the builder pattern:

```go
eventListener := ga.
    NewEventListener().
    EventTypes("zwave_js_value_notification"). // Specify one or more event types
    Call(myCallbackFunc).                     // Specify the callback function
    Build()

// Register the listener with your app
app.RegisterEventListeners(eventListener)
```

Event listeners have other functions to change the behavior.

| Function                                | Info                                                                                                              |
| --------------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| OnlyBetween("03:00", "14:00")           | Only run your function between two specified times of day                                                         |
| OnlyAfter("03:00")                      | Only run your function after a specified time of day                                                              |
| OnlyBefore("03:00")                     | Only run your function before a specified time of day                                                             |
| Throttle("30s")                         | Minimum time between function calls                                                                               |
| ExceptionDates(time.Time, ...time.Time) | A one time exception on the given date. Time is ignored, applies to whole day. Functions like a "blocklist"      |
| ExceptionRange(time.Time, time.Time)    | A one time exception between the two date/times. Both date and time are considered. Functions like a "blocklist" |

The callback function receives three parameters:
```go
func myCallback(service *ga.Service, state ga.State, data ga.EventData) {
    // You can unmarshal the raw JSON into a type-safe struct
    ev := ga.EventZWaveJSValueNotification{}
    json.Unmarshal(data.RawEventJSON, &ev)
    
    // Handle the event...
}
```

> 💡 Check `eventTypes.go` for pre-defined event types, or create your own struct type for custom events and contribute them back to gome-assistant with a PR.

### Interval

Intervals are used to run a function on an interval.

```go
// run every hour at the 30-minute mark
interval := ga.NewInterval().Call(myFunc).Every("1h").StartingAt("00:30").Build()
// run every 5 minutes between 10am and 5pm
interval = ga.NewInterval().Call(myFunc).Every("5m").StartingAt("10:00").EndingAt("17:00").Build()
```

Intervals have other functions to change the behavior.

| Function                                | Info                                                                                |
| --------------------------------------- | ----------------------------------------------------------------------------------- |
| StartingAt(TimeString)                  | What time the interval begins to run each day.                                      |
| EndingAt(TimeString)                    | What time the interval stops running each day.                                      |
| ExceptionDates(time.Time, ...time.Time) | A one time exception on the given date. Time is ignored, applies to whole day.      |
| ExceptionRange(time.Time, time.Time)    | A one time exception between the two date/times. Both date and time are considered. |

#### Interval Callback function

The function passed to `.Call()` must take

- `*ga.Service` used to call home assistant services
- `*ga.State` used to retrieve state from home assistant

```go
func myFunc(se *ga.Service, st *ga.State) {
  // ...
}
```

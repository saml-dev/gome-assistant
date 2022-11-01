# Gome-Assistant

Golang ↔️ Home Assistant

Write your [Home Assistant](https://www.home-assistant.io/) automations with a strongly-typed Golang library!

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

## Disclaimer

Gome-Assistant is a new library, and I'm opening it up early to get some user feedback on the API and help shape the direction. I plan for it to grow to cover all Home Assistant use cases, services, and event types. So it's possible — maybe likely — that breaking changes will happen before v1.0.0!

## gome-assistant Concepts (TODO)

First, you'll need to create your app.

```go
import ga "github.com/saml-dev/gome-assistant"

// replace with IP and port of your Home Assistant installation if needed
app := ga.App("0.0.0.0:8123")
```

A full reference is available on [pkg.go.dev](https://pkg.go.dev/github.com/saml-dev/gome-assistant), but all you need to know to get started are the three types of automations in gome-assistant.

### Schedules

Schedules are as you expect, a way to run a function on a schedule. The most common schedule will be once a day.

```go
_7pm := ga.NewSchedule().Call(myFunc).Daily().At("19:00").Build()
```

Schedules can also be run at sunrise or sunset, with an optional [offset](https://pkg.go.dev/time#ParseDuration).

```go
// 30 mins before sunrise
sunrise := ga.NewSchedule().Call(myFunc).Daily().Sunrise(app, "-30m").Build()
```

Schedules are also used to run a function on a certain interval. Offset is used to offset the first run of a schedule from midnight.

```go
// run every hour at the 30-minute mark
interval := ga.NewSchedule().Call(a).Every("1h").Offset("30m").Build()
```

All schedules must be registered with your app. This will panic if there are any issues with the schedule.

```go
app.RegisterSchedules(_7pm, sunrise, interval)
```

#### Schedule Callback function

The function passed to `.Call()` must take

- `*ga.Service` used to call home assistant services
- `*ga.State` used to retrieve state from home assistant

### Entity Listeners

### Event Listeners

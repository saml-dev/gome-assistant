# Gome-Assistant

Golang ↔️ Home Assistant

Write your [Home Assistant](https://www.home-assistant.io/) automations with a strongly-typed Golang library!

## Why?

My Home Assistant automation journey started with [Node-RED](https://nodered.org/). Since I already know how to write javascript, I started stuffing all my logic into code nodes in Node-RED.

Then one day I stumbled on [Appdaemon](https://appdaemon.readthedocs.io/en/latest/), which lets you write automations entirely in Python. I switched all my automations over because of the flexibility I gained from writing my automations in code.

While autocomplete in your IDE seems trivial, it's an important feature many developers take for granted. Python isn't great at it.

I wanted to learn Golang, and had the idea to wrap the Home Assistant websocket API with a Go library. Go's strong typing means I don't have to remember the API because my IDE will fill in the blanks for me. So anyway, here's my first Go project, I hope you find it useful 😁

## Quick Start

### Installation

```
go get github.com/saml-dev/gome-assistant
```

### Write your automations

Check out `example/example.go` for an example of the 3 types of automations — schedules, entity listeners, and event listeners. Instead of copying and pasting, try typing it our yourself to see how autocomplete guides you through the setup using a builder pattern. You can also check out some of the other options you see in the autocomplete.

### Run your code

Keeping with the simplicity that Go is famous for, you don't need a whole environment or docker container to run Gome-Assistant. It's just a binary like any other Go code you would write. So once you have your automations, you can run it however you like — using `screen` or `tmux`, a cron job, or wrap it up in a docker container if you just can't get enough docker!

_Note: I may provide a Docker image in the future with file watching to restart gome-assistant, to make it easier to use gome-assistant on a fully managed Home Assistant installation._

## Disclaimer

Gome-Assistant is a new library, and I'm opening it up early to get some user feedback on the API and help shape the direction. I plan for it to grow to cover all Home Assistant use cases, services, and event types. So it's possible — maybe likely — that breaking changes will happen before v1.0.0!

## API Reference (WIP)

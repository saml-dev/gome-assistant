package gomeassistant

import (
	"context"
	"time"
)

type app struct {
	url       string
	ctx       context.Context
	schedules []Schedule
	listeners []Listener
}

func App(url string) (app, error) {
	// TODO: connect to websocket, return error if fails
	return app{url: url}, nil
}

func (a app) RegisterSchedule(s Schedule) {
	if s.err != nil {
		panic(s.err) // something wasn't configured properly when the schedule was built
	}

	if s.frequency == 0 {
		panic("A schedule must call either Daily() or Every() when built.")
	}

	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()) // start at midnight today

	// apply offset if set
	if s.offset.Hour != 0 || s.offset.Minute != 0 {
		startTime.Add(time.Hour * time.Duration(s.offset.Hour))
		startTime.Add(time.Minute * time.Duration(s.offset.Minute))
	}

	// advance first scheduled time by frequency until it is in the future
	for startTime.Before(now) {
		startTime = startTime.Add(s.frequency)
	}

	// TODO: save realStartTime or _startTime to s, add to list of Schedules
}

const (
	FrequencyMissing time.Duration = 0

	Daily    time.Duration = time.Hour * 24
	Hourly   time.Duration = time.Hour
	Minutely time.Duration = time.Minute
)

type Listener struct {
}

const (
	_0000 string = "0000"
	_0015 string = "0015"
	_0030 string = "0030"
	_0045 string = "0045"
	_0100 string = "0100"
	_0115 string = "0115"
	_0130 string = "0130"
	_0145 string = "0145"
	_0200 string = "0200"
	_0215 string = "0215"
	_0230 string = "0230"
	_0245 string = "0245"
	_0300 string = "0300"
	_0315 string = "0315"
	_0330 string = "0330"
	_0345 string = "0345"
	_0400 string = "0400"
	_0415 string = "0415"
	_0430 string = "0430"
	_0445 string = "0445"
	_0500 string = "0500"
	_0515 string = "0515"
	_0530 string = "0530"
	_0545 string = "0545"
	_0600 string = "0600"
	_0615 string = "0615"
	_0630 string = "0630"
	_0645 string = "0645"
	_0700 string = "0700"
	_0715 string = "0715"
	_0730 string = "0730"
	_0745 string = "0745"
	_0800 string = "0800"
	_0815 string = "0815"
	_0830 string = "0830"
	_0845 string = "0845"
	_0900 string = "0900"
	_0915 string = "0915"
	_0930 string = "0930"
	_0945 string = "0945"
	_1000 string = "1000"
	_1015 string = "1015"
	_1030 string = "1030"
	_1045 string = "1045"
	_1100 string = "1100"
	_1115 string = "1115"
	_1130 string = "1130"
	_1145 string = "1145"
	_1200 string = "1200"
	_1215 string = "1215"
	_1230 string = "1230"
	_1245 string = "1245"
	_1300 string = "1300"
	_1315 string = "1315"
	_1330 string = "1330"
	_1345 string = "1345"
	_1400 string = "1400"
	_1415 string = "1415"
	_1430 string = "1430"
	_1445 string = "1445"
	_1500 string = "1500"
	_1515 string = "1515"
	_1530 string = "1530"
	_1545 string = "1545"
	_1600 string = "1600"
	_1615 string = "1615"
	_1630 string = "1630"
	_1645 string = "1645"
	_1700 string = "1700"
	_1715 string = "1715"
	_1730 string = "1730"
	_1745 string = "1745"
	_1800 string = "1800"
	_1815 string = "1815"
	_1830 string = "1830"
	_1845 string = "1845"
	_1900 string = "1900"
	_1915 string = "1915"
	_1930 string = "1930"
	_1945 string = "1945"
	_2000 string = "2000"
	_2015 string = "2015"
	_2030 string = "2030"
	_2045 string = "2045"
	_2100 string = "2100"
	_2115 string = "2115"
	_2130 string = "2130"
	_2145 string = "2145"
	_2200 string = "2200"
	_2215 string = "2215"
	_2230 string = "2230"
	_2245 string = "2245"
	_2300 string = "2300"
	_2315 string = "2315"
	_2330 string = "2330"
	_2345 string = "2345"
)

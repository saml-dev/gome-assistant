package internal

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"time"

	"github.com/golang-module/carbon"
)

type EnabledDisabledInfo struct {
	Entity     string
	State      string
	RunOnError bool
}

var id int64 = 0

func GetId() int64 {
	id += 1
	return id
}

// Parses a HH:MM string.
func ParseTime(s string) carbon.Carbon {
	t, err := time.Parse("15:04", s)
	if err != nil {
		log.Fatalf("Failed to parse time string \"%s\"; format must be HH:MM.\n", s)
	}
	return carbon.Now().SetTimeMilli(t.Hour(), t.Minute(), 0, 0)
}

func ParseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Couldn't parse string duration: \"%s\" see https://pkg.go.dev/time#ParseDuration for valid time units\n", s))
	}
	return d
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

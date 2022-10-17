package internal

import (
	"log"
	"time"

	"github.com/golang-module/carbon"
)

var id int64 = 0

func GetId() int64 {
	id += 1
	return id
}

func ParseTime(s string) carbon.Carbon {
	t, err := time.Parse("15:04", s)
	if err != nil {
		log.Fatalf("Failed to parse time string \"%s\"; format must be HH:MM.", s)
	}
	return carbon.Now().StartOfDay().SetHour(t.Hour()).SetMinute(t.Minute())
}

package message

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

type TimeStamp time.Time

func (ts TimeStamp) Time() time.Time {
	return time.Time(ts)
}

// UnmarshalJSON unmarshals a timestamp from JSON. HA sometimes
// formats timestamps as RFC 3339 strings, sometimes as fractional
// seconds since the epoch. Handle either one (without recording which
// one it was).
func (ts *TimeStamp) UnmarshalJSON(b []byte) error {
	t := (*time.Time)(ts)
	if err := t.UnmarshalJSON(b); err == nil {
		return nil
	}

	var v float64
	if err := json.Unmarshal(b, &v); err == nil {
		seconds, fraction := math.Modf(v)
		*t = time.Unix(int64(seconds), int64(fraction*1e+9))
		return nil
	}

	return fmt.Errorf("unmarshaling timestamp: '%s'", string(b))
}

package websocket

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

type TimeStamp time.Time

// UnmarshalJSON unmarshals a timestamp from JSON. HA sometimes
// formats timestamps as RFC 3339 strings, sometimes as fractional
// seconds since the epoch. Handle either one (without recording which
// one it was).
func (ts *TimeStamp) UnmarshalJSON(b []byte) error {
	if err := (*time.Time)(ts).UnmarshalJSON(b); err == nil {
		return nil
	}

	var v float64
	if err := json.Unmarshal(b, &v); err == nil {
		seconds := math.Floor(v)
		*(*time.Time)(ts) = time.Unix(int64(seconds), int64((v-seconds)*1e+9))
		return nil
	}

	return fmt.Errorf("unmarshaling timestamp: '%s'", string(b))
}

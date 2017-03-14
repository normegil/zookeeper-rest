package json

import (
	"time"
)

type JSONTime time.Time

func (j JSONTime) MarshalJSON() ([]byte, error) {
	jsonTime := "\"" + j.String() + "\""
	return []byte(jsonTime), nil
}

func (j *JSONTime) UnmarshalJSON(b []byte) error {
	t, err := time.Parse(time.RFC3339, string(b))
	if nil != err {
		return err
	}
	jsonTime := JSONTime(t)
	j = &jsonTime
	return nil
}

func (j JSONTime) String() string {
	return time.Time(j).Format(time.RFC3339)
}

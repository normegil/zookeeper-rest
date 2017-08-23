package formats

import (
	"time"

	"github.com/pkg/errors"
)

type Time time.Time

func (j Time) MarshalJSON() ([]byte, error) {
	json := "\"" + j.String() + "\""
	return []byte(json), nil
}

func (j *Time) UnmarshalJSON(b []byte) error {
	toUnmarshal := j.clean(string(b))
	t, err := time.Parse(time.RFC3339, toUnmarshal)
	if nil != err {
		return errors.Wrapf(err, "Could not Unmarshall %s into Time", toUnmarshal)
	}
	time := Time(t)
	j = &time
	return nil
}

func (j *Time) clean(toClean string) string {
	toReturn := toClean
	if '"' == toReturn[0] {
		toReturn = toReturn[1:]
	}
	if '"' == toReturn[len(toReturn)-1] {
		toReturn = toReturn[:len(toReturn)-1]
	}
	return toReturn
}

func (j Time) String() string {
	return time.Time(j).Format(time.RFC3339)
}

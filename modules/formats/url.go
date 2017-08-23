package formats

import (
	"net/url"
	"strings"
)

type URL struct {
	*url.URL
}

func (j URL) MarshalJSON() ([]byte, error) {
	return []byte("\"" + j.String() + "\""), nil
}

func (j *URL) UnmarshalJSON(b []byte) error {
	toParse := string(b)
	if strings.HasPrefix(toParse, "\"") {
		toParse = toParse[1:]
	}
	if strings.HasSuffix(toParse, "\"") {
		toParse = toParse[:len(toParse)-1]
	}
	u, err := url.Parse(toParse)
	if nil != err {
		return err
	}
	j.URL = u
	return nil
}

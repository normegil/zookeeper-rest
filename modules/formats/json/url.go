package json

import "net/url"

type JSONURL url.URL

func (j JSONURL) MarshalJSON() ([]byte, error) {
	return []byte("\"" + j.String() + "\""), nil
}

func (j *JSONURL) UnmarshalJSON(b []byte) error {
	u, err := url.Parse(string(b))
	if nil != err {
		return err
	}
	jsonURL := JSONURL(*u)
	j = &jsonURL
	return nil
}

func (j JSONURL) String() string {
	original := url.URL(j)
	ptr := &original
	return ptr.String()
}

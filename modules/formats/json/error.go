package json

type ErrorJSON struct {
	Err error
}

func (j ErrorJSON) MarshalJSON() ([]byte, error) {
	jsonErr := "\"" + j.Error() + "\""
	return []byte(jsonErr), nil
}

func (j ErrorJSON) Error() string {
	return j.Err.Error()
}

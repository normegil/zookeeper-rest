package error

import "encoding/json"

// Error can hold an error message
type Error struct {
	Message string
}

func (j Error) MarshalJSON() ([]byte, error) {
	toJson := make(map[string]interface{})
	toJson["@type"] = "BaseError"
	toJson["message"] = j.Message
	return json.Marshal(toJson)
}

func (j *Error) UnmarshalJSON(b []byte) error {
	toError := make(map[string]*json.RawMessage)
	err := json.Unmarshal(b, &toError)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(*toError["message"]), &j.Message)
}

//Error exist to comply with errors.Error interface
func (j Error) Error() string {
	return j.Message
}

package json

import (
	"leaderboard/pkg/encoder"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// jsonEncoder -
type jsonEncoder struct {
	encoderType string
}

// Encoder -
func (j *jsonEncoder) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Decode -
func (j *jsonEncoder) Decode(d []byte, v interface{}) error {
	return json.Unmarshal(d, v)
}

// Type -
func (j jsonEncoder) Type() string {
	return j.encoderType
}

// NewEncoder -
func NewEncoder() encoder.Encoder {
	return &jsonEncoder{
		encoderType: "json",
	}
}

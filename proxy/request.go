package proxy

import (
	"encoding/json"
)

type RequestAddr struct {
	Host      string
	Port      string
	Network   string
	Timestamp string
	Random    string
}

func (r *RequestAddr) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *RequestAddr) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &r)
}

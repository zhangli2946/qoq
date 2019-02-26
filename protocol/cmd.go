package protocol

import "encoding/json"

// Command struct
// Protocol
type Command struct {
	Handle  string          `json:"cmd"`
	Payload json.RawMessage `json:"payload"`
}

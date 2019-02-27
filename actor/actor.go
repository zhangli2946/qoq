package actor

import (
	"encoding/json"
	"errors"
	"qoq/protocol"
)

var (
	// ErrNotSupport def
	ErrNotSupport = errors.New("cmd not support")
)

// GenError message
func GenError(cmd protocol.Command, step int, err error) (string, []byte) {
	if err == nil {
		err = ErrNotSupport
	}
	b, _ := json.Marshal(map[string]interface{}{
		"code": 128,
		"stat": "error",
		"cmd":  cmd,
		"payload": map[string]interface{}{
			"error": err.Error(),
			"step":  step,
		},
	})
	return "error", b
}

// GenAck message
func GenAck(cmd protocol.Command, hdlr int) (string, []byte) {
	b, _ := json.Marshal(map[string]interface{}{
		"code": 0,
		"stat": "ok",
		"cmd":  cmd,
		"payload": map[string]interface{}{
			"hdlr": hdlr,
		},
	})
	return "ack", b
}

package main

import (
	"encoding/json"
	"qoq/actor/emq"

	client "github.com/eclipse/paho.mqtt.golang"
)

type describe struct {
	Version string   `json:"Version"`
	SysHdlr []string `json:"SysHdlr"`
	BizHdlr []string `json:"BizHdlr"`
	DefHdlr []string `json:"DefHdlr"`
}

// Describe handler
func Describe(i *emq.Emq, cmd json.RawMessage) error {
	var e error
	var b []byte
	var t client.Token
	splugins := []string{}
	for k := range i.SysProtocol {
		splugins = append(splugins, k)
	}
	bplugins := []string{}
	for k := range i.BizProtocol {
		bplugins = append(bplugins, k)
	}

	dplugins := []string{}
	for k := range i.BizProtocol {
		dplugins = append(dplugins, k)
	}

	b, e = json.Marshal(describe{
		Version: i.Version,
		SysHdlr: splugins,
		BizHdlr: bplugins,
		DefHdlr: dplugins,
	})
	if e != nil {
		return e
	}
	if t = i.Cli.Publish("describe", 2, false, b); t.Wait() && t.Error() != nil {
		return t.Error()
	}
	return nil
}

func Gen() (string, func(*emq.Emq, json.RawMessage) error) {
	return "describe", Describe
}

func main() {
}

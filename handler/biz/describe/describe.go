package main

import (
	"encoding/json"
	"qoq/actor/emq"
)

type describe struct {
	Version string   `json:"Version"`
	SysHdlr []string `json:"SysHdlr"`
	BizHdlr []string `json:"BizHdlr"`
	DefHdlr []string `json:"DefHdlr"`
}

// Describe handler
func Describe(i *emq.Emq, cmd json.RawMessage) error {
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

	i.Cli.Publish("describe", 2, false, describe{
		Version: i.Version,
		SysHdlr: splugins,
		BizHdlr: bplugins,
		DefHdlr: dplugins,
	})
	return nil
}

func Gen() (string, func(*emq.Emq, json.RawMessage) error) {
	return "describe", Describe
}

func main() {
}

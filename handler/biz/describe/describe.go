package main

import (
	"qoq/actor/emq"
	"encoding/json"
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

func gen() (string, func(*emq.Emq, json.RawMessage) error) {
	return "Describe", Describe
}

func main() {
}

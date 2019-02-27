package load

import (
	"encoding/json"
	"errors"
	"path"
	"plugin"
	"qoq/actor/emq"
)

type load struct {
	Dst   string `json:"dst"`
	Lib   string `json:"lib"`
	Funcs string `json:"func"`
}

func init() {
	emq.DefaultHandler["load"] = Load
}

// Load plugin once
func Load(i *emq.Emq, c json.RawMessage) (e error) {
	var d load
	var f func() (string, func(i *emq.Emq, c json.RawMessage) error)
	var k string
	var o bool
	var p *plugin.Plugin
	var s plugin.Symbol
	var v func(i *emq.Emq, c json.RawMessage) error

	e = json.Unmarshal(c, &d)
	if e != nil {
		return
	}
	p, e = plugin.Open(path.Join([]string{d.Dst, d.Lib}...))
	if e != nil {
		return
	}
	s, e = p.Lookup("Gen")
	if e != nil {
		return
	}
	f, o = s.(func() (string, func(i *emq.Emq, c json.RawMessage) error))
	if !o {
		e = errors.New("type not support")
		return
	}
	k, v = f()
	i.BizProtocol[k] = v
	return nil
}

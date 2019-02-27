package emq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"qoq/actor"
	"qoq/protocol"
	"time"

	client "github.com/eclipse/paho.mqtt.golang"
)

// Handler func
// handler
type Handler func(*Emq, json.RawMessage) error

// FSM interface
type FSM interface {
	Init(opts []Option)
	Run()
	Stop()
}

var (
	// DefaultHandler root
	DefaultHandler     = map[string]Handler{}
	_              FSM = &Emq{}
)

func init() {
	client.DEBUG = log.New(os.Stderr, "DEBUG    ", log.Ltime)
	client.WARN = log.New(os.Stderr, "WARNING  ", log.Ltime)
	client.CRITICAL = log.New(os.Stderr, "CRITICAL ", log.Ltime)
	client.ERROR = log.New(os.Stderr, "ERROR    ", log.Ltime)
}

// Emq type
type Emq struct {
	Cli         client.Client
	Opts        *client.ClientOptions
	Evt         chan client.Message
	SysProtocol map[string]client.MessageHandler
	BizProtocol map[string]Handler
	Cancel      context.CancelFunc
	Ctx         context.Context
	Version     string
}

func (i *Emq) cleanup() {
	i.Cli.Disconnect(500)
}

// Init method
func (i *Emq) Init(opts []Option) {
	i.Evt = make(chan client.Message, 10)
	i.SysProtocol = make(map[string]client.MessageHandler)
	i.BizProtocol = make(map[string]Handler)
	for _, o := range opts {
		if err := o(i); err != nil {
			panic(o)
		}
	}
	if i.Opts == nil {
		panic(errors.New("Not configured"))
	}

	var tk client.Token
	i.Cli = client.NewClient(i.Opts)
	if tk = i.Cli.Connect(); tk.Wait() && tk.Error() != nil {
		panic(tk.Error())
	}

	var rcnt int
	for t, h := range i.SysProtocol {
		rcnt = 0
	loop:
		if tk = i.Cli.Subscribe(t, 0, h); tk.Wait() && tk.Error() != nil {
			rcnt++
			time.Sleep(1 * time.Second)
			if rcnt > 3 {
				panic(tk.Error())
			}
			goto loop
		}
	}
}

// Run method
func (i *Emq) Run() {
	i.Ctx, i.Cancel = context.WithCancel(context.Background())
	go func(ctx context.Context) {
		var tk client.Token
		var t string
		var p []byte
		var e error
		var h Handler
		var g bool
		var c protocol.Command

		for {
			select {
			case evt := <-i.Evt:
				if e = json.Unmarshal(evt.Payload(), &c); e != nil {
					t, p = actor.GenError(protocol.Command{}, 0, e)
					i.Cli.Publish(t, 1, false, p)
					continue
				}
				if h, g = i.BizProtocol[c.Handle]; g {
					if e = h(i, c.Payload); e == nil {
						t, p = actor.GenAck(c, 1)
						i.Cli.Publish(t, 2, false, p)
						continue
					}
				}
				if h, g = DefaultHandler[c.Handle]; g {
					if e = h(i, c.Payload); e == nil {
						t, p = actor.GenAck(c, 0)
						i.Cli.Publish(t, 2, false, p)
						continue
					}
				}
				if e != nil {
					t, p = actor.GenError(c, 1, e)
				} else {
					t, p = actor.GenError(c, 1, nil)
				}
				i.Cli.Publish(t, 1, false, p)
			case <-ctx.Done():
				t, p = actor.GenError(c, 0, ctx.Err())
				tk = i.Cli.Publish(t, 1, false, p)
				if !tk.WaitTimeout(time.Second * 1) {
					fmt.Println("shutdown message may swallow")
				}
				goto cleanup
			}
		}
	cleanup:
		i.cleanup()
	}(i.Ctx)
}

// Stop method
func (i *Emq) Stop() {
	i.Cancel()
}

func gen(c protocol.Command) func(*map[string]Handler, *Handler) bool {
	return func(m *map[string]Handler, h *Handler) (g bool) {
		var H Handler
		H, g = (*m)[c.Handle]
		h = &H
		return
	}
}

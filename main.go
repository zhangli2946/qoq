package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"qoq/actor/emq"
	"qoq/handler/biz/download"
	_ "qoq/handler/biz/load"
	"time"

	client "github.com/eclipse/paho.mqtt.golang"
)

func init() {
	client.DEBUG = log.New(os.Stderr, "DEBUG    ", log.Ltime)
	client.WARN = log.New(os.Stderr, "WARNING  ", log.Ltime)
	client.CRITICAL = log.New(os.Stderr, "CRITICAL ", log.Ltime)
	client.ERROR = log.New(os.Stderr, "ERROR    ", log.Ltime)
}

func main() {
	var i = &emq.Emq{}
	i.Init([]emq.Option{
		emq.AddBizHandler("download", download.Download),
		emq.MQTTSettings(
			os.Getenv("MQTTSRV"),
			os.Getenv("MQTTUSR"),
			os.Getenv("MQTTPWD")),
		emq.AddSysHandler("bus/#", func(cli client.Client, evt client.Message) {
			i.Evt <- evt
		}),
		emq.MQTTOnConnectHandler(func(client.Client) {
			var rcnt int
			var tk client.Token
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
		}),
		emq.MQTTConnectionLostHandler(func(cli client.Client, err error) {
			fmt.Println(err.Error())
		})})

	go i.Run()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		i.Stop()
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"qoq/actor/emq"
	"qoq/handler/biz/download"
	_ "qoq/handler/biz/load"
	"sync"
	"time"

	client "github.com/eclipse/paho.mqtt.golang"
	getter "github.com/hashicorp/go-getter"
)

func main2() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting wd: %s", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	client := &getter.Client{
		Ctx:     ctx,
		Src:     "http://ftps.lowaniot.com/tools/qrsbox_2016_11_24.zip",
		Dst:     "qrsbox",
		Pwd:     pwd,
		Mode:    getter.ClientModeAny,
		Options: []getter.ClientOption{},
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	errChan := make(chan error, 2)
	go func() {
		defer wg.Done()
		defer cancel()
		if err := client.Get(); err != nil {
			errChan <- err
		}
	}()
	select {
	case <-ctx.Done():
		wg.Wait()
		log.Printf("success!")
	case err := <-errChan:
		wg.Wait()
		log.Fatalf("Error downloading: %s", err)
	}
}

func main() {
	var i = &emq.Emq{}
	i.Init([]emq.Option{
		emq.AddBizHandler("download", download.Download),
		emq.MQTTSettings(
			os.Getenv("MQTTSRV"),
			os.Getenv("MQTTUSR"),
			os.Getenv("MQTTPWD"),
			os.Getenv("MQTTCLID")),
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

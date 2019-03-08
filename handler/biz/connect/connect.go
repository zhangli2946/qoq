package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"qoq/actor"
	"qoq/actor/emq"
	"qoq/protocol"

	client "github.com/eclipse/paho.mqtt.golang"
)

type input struct {
	C int `json:"cnt"`
	B int `json:"cps"`
}

// Connect handler
func Connect(i *emq.Emq, cmd json.RawMessage) (err error) {
	var d input
	err = json.Unmarshal(cmd, &d)
	if err != nil {
		return
	}
	if d.B == 0 {

	} else {
		for cnt := d.B; cnt > 0; cnt-- {
			for step := d.C / d.B; step > 0; step-- {
				go func(ctx context.Context, c int, s int) {
					var tk client.Token
					msg := make(chan client.Message)
					cli := client.NewClient(client.NewClientOptions().
						SetClientID(fmt.Sprintf("%d@%d", c, s)).
						SetAutoReconnect(false).
						SetConnectionLostHandler(func(cli client.Client, err error) {
							t, p := actor.GenError(protocol.Command{
								Handle: "connect"}, 0, err)
							i.Cli.Publish(t, 1, false, p)
						}).
						SetOnConnectHandler(func(cli client.Client) {
							if tk := cli.Subscribe("connect/#", 0, func(cli client.Client, evt client.Message) {
								msg <- evt
							}); tk.Wait() && tk.Error() != nil {
								t, p := actor.GenError(protocol.Command{
									Handle: "connect"}, 0, tk.Error())
								i.Cli.Publish(t, 1, false, p)
							}
						}).
						AddBroker(os.Getenv("MQTTSRV")).
						SetUsername(os.Getenv("MQTTUSR")).
						SetPassword(os.Getenv("MQTTPWD")))

					if tk = cli.Connect(); tk.Wait() && tk.Error() != nil {
						fmt.Println(tk.Error(), "@", c, "@", s)
					}
					for {
						select {
						case <-i.Ctx.Done():
							cli.Disconnect(100)
						}
					}

				}(i.Ctx, cnt, step)
			}
		}
	}
	return nil
}

// Gen func
func Gen() (string, func(*emq.Emq, json.RawMessage) error) {
	return "connect", Connect
}

func main() {
}

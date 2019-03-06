package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"qoq/actor/emq"
	"time"

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
		cnt := d.B
		for ; cnt > 0; cnt-- {
			for step := d.C / d.B; step > 0; step-- {
				go func(ctx context.Context, c int, s int) {
					var tk client.Token
					cli := client.NewClient(client.NewClientOptions().
						SetClientID(fmt.Sprintf("%d@%d", c, s)).
						SetAutoReconnect(false).
						SetConnectionLostHandler(func(cli client.Client, err error) {
							fmt.Println(err.Error(), "@", c, "@", s)
						}).
						AddBroker(os.Getenv("MQTTSRV")).
						SetUsername(os.Getenv("MQTTUSR")).
						SetPassword(os.Getenv("MQTTPWD")))

					if tk = cli.Connect(); tk.Wait() && tk.Error() != nil {
						fmt.Println(tk.Error(), "@", c, "@", s)
					}
					time.Sleep(100 * time.Second)
					cli.Disconnect(250)
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

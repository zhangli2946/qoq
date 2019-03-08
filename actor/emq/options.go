package emq

import (
	"time"

	client "github.com/eclipse/paho.mqtt.golang"
)

// Option def
// inject options to actor
type Option func(*Emq) error

// AddSysHandler option
func AddSysHandler(topic string, handler client.MessageHandler) Option {
	return func(i *Emq) error {
		i.SysProtocol[topic] = handler
		return nil
	}
}

// AddBizHandler option
func AddBizHandler(topic string, handler Handler) Option {
	return func(i *Emq) error {
		i.BizProtocol[topic] = handler
		return nil
	}
}

// MQTTSettings option
func MQTTSettings(Broker string, Username string, Password string) Option {
	return func(i *Emq) error {
		i.Opts = client.NewClientOptions().
			AddBroker(Broker).
			SetUsername(Username).
			SetPassword(Password).
			SetMaxReconnectInterval(10 * time.Second).
			SetCleanSession(false)
		return nil
	}
}

// MQTTConnectionLostHandler option
func MQTTConnectionLostHandler(f client.ConnectionLostHandler) Option {
	return func(i *Emq) error {
		i.Opts.SetConnectionLostHandler(f)
		return nil
	}
}

// MQTTOnConnectHandler option
func MQTTOnConnectHandler(f client.OnConnectHandler) Option {
	return func(i *Emq) error {
		i.Opts.SetOnConnectHandler(f)
		return nil
	}
}

// SetOnConnectHandler(func(conn client.Client) {
// 	var token client.Token
// 	rccount := 0
// 	topics := make([]string, len(ist.sysHandler))
// 	i := 0
// 	for topic := range ist.sysHandler {
// 		topics[i] = topic
// 		i++
// 	}
// 	for i = 0; i < len(ist.sysHandler); {
// 		token = conn.Subscribe(topics[i], 0, ist.sysHandler[topics[i]])
// 		if token.Wait() && token.Error() != nil {
// 			fmt.Println(token.Error())
// 			rccount++
// 			time.Sleep(1 * time.Second)
// 			if rccount > 3 {
// 				panic("sub failed exceed 3")
// 			}
// 			continue
// 		}
// 		i++
// 	}
// })

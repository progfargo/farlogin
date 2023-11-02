package app

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

var NC *nats.Conn
var NatsUrl string
var NatsPort int

func Connect(nodeName string) {

	NatsPort = 4222
	NatsUrl = "localhost"

	var err error
	NC, err = nats.Connect(fmt.Sprintf("nats://%s:%d", NatsUrl, NatsPort),
		nats.Timeout(10*time.Second),
		nats.Name(fmt.Sprintf("flnode: %s", nodeName)),
		nats.UserInfo("farlogin", "Pseudo12Login"))

	if err != nil {
		panic("in app: " + err.Error())
	} else {
		println("connected.")
	}
}

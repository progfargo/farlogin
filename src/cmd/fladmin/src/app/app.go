package app

import (
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

var NC *nats.Conn

var NatsUrl string
var NatsPort int

func Connect() {
	var err error

	NatsPort = 4222
	NatsUrl = "localhost"

	name, err := os.Hostname()
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return
	}

	NC, err = nats.Connect(fmt.Sprintf("nats://%s:%d", NatsUrl, NatsPort),
		nats.Timeout(10*time.Second),
		nats.Name(fmt.Sprintf("fladmin: %s", name)),
		nats.UserInfo("farlogin", "Pseudo12Login"))

	if err != nil {
		panic(err)
	}
}

package main

import (
	"flag"
	"fmt"

	//"errors"
	//"io/fs"
	"log"
	"time"

	"flnode/src/app"
	"flnode/src/session"

	"github.com/nats-io/nats.go"
)

func main() {
	nodeName := flag.String("node", "", "node name")
	flag.Parse()

	if *nodeName == "" {
		fmt.Printf("usage: flnode -node=<node name>\n")
		flag.PrintDefaults()
		return
	}

	app.Connect(*nodeName)

	cron(*nodeName)

	sub, err := app.NC.Subscribe(fmt.Sprintf("newSession.%s", *nodeName), func(msg *nats.Msg) {
		ses, err := session.NewSession(msg.Data)
		if err != nil {
			println(err.Error())
			log.Println(err.Error())
		}

		isValid, err := session.IsValidSession(ses.SessionHash)
		if err != nil {
			log.Println(err.Error())
			return
		}

		var res string
		if isValid {
			res = "true"
		} else {
			res = "false"
		}

		err = app.NC.Publish(msg.Reply, []byte(res))
		if err != nil {
			log.Println(err.Error())
		}

		if !isValid {
			fmt.Printf("invalid hash returns.")
			return
		}

		go func() {
			err = ses.Listen()
			if err != nil {
				/*
					var pErr *fs.PathError
					if errors.As(err, &pErr) {
						return
					}
				*/

				log.Println(err.Error())
				return
			}
		}()
	})

	if err != nil {
		log.Println(err.Error())
	}

	defer sub.Unsubscribe()

	select {}
}

func cron(nodeName string) {
	//clear old sessions
	go func() {
		c := time.Tick(time.Duration(5) * time.Second)
		for _ = range c {
			err := app.NC.Publish("hello", []byte(nodeName))
			if err != nil {
				println(err.Error())
				log.Println(err.Error())
			}
		}
	}()
}

package main

import (
	"flag"
	"fmt"
	"time"

	"fladmin/src/app"
	"fladmin/src/client"
	"fladmin/src/session"
)

func main() {
	sessionHash := flag.String("hash", "", "session hash")
	nodeName := flag.String("node", "", "node name")
	flag.Parse()

	app.Connect()

	if *sessionHash == "" || *nodeName == "" {
		fmt.Printf("usage: fladmin -node=<node name> -hash=<session hash>\n")
		flag.PrintDefaults()
		return
	}

	ses := session.NewSession(*sessionHash)
	jsonData, err := ses.Marshal()
	if err != nil {
		panic("1. " + err.Error())
	}

	response, err := app.NC.Request(fmt.Sprintf("newSession.%s", *nodeName), jsonData, 5*time.Second)
	if err != nil {
		panic("3. " + err.Error())
	}

	if string(response.Data) != "true" {
		panic("Invalid session hash.")
	}

	client := client.NewClient(*sessionHash)

	err = client.Run()
	if err != nil {
		panic(err)
	}
}

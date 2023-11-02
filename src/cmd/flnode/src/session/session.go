package session

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"time"

	"flnode/src/app"

	"github.com/creack/pty"
	"github.com/nats-io/nats.go"
)

type session struct {
	UserName    string
	SessionHash string
}

func NewSession(jsonData []byte) (*session, error) {
	rv := new(session)
	err := json.Unmarshal(jsonData, rv)
	if err != nil {
		return nil, err
	}

	return rv, nil
}

func (ses *session) Write(p []byte) (n int, err error) {
	err = app.NC.Publish(fmt.Sprintf("admin.%s", ses.SessionHash), p)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func IsValidSession(sessionHash string) (bool, error) {
	response, err := app.NC.Request("farlogin.isValidHash", []byte(sessionHash), 1*time.Second)
	if err != nil {
		return false, err
	}

	if string(response.Data) == "true" {
		return true, nil
	}

	return false, nil
}

func (ses *session) Listen() error {
	//fmt.Println("session listening...")
	// Create command
	c := exec.Command("bash")

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		println("6")
		return err
	}

	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	isFirstMsg := true
	sub, err := app.NC.Subscribe(fmt.Sprintf("node.%s", ses.SessionHash), func(msg *nats.Msg) {
		if isFirstMsg {
			err = ses.SessionUsed()
			if err != nil {
				println("3", err.Error())
				return
			}

			isFirstMsg = false
		}

		_, err := ptmx.Write(msg.Data)
		if err != nil {
			println("4")
			return
		}
	})

	defer sub.Unsubscribe()

	if err != nil {
		println("1")
		return err
	}

	_, err = io.Copy(ses, ptmx)
	if err != nil {
		return err
	}

	return nil
}

func (ses *session) SessionUsed() error {
	err := app.NC.Publish("farlogin.sessionUsed", []byte(ses.SessionHash))
	if err != nil {
		return err
	}

	return nil
}

package client

import (
	"bufio"
	"fmt"
	"os"

	"fladmin/src/app"

	"github.com/nats-io/nats.go"
	"golang.org/x/crypto/ssh/terminal"
)

type client struct {
	sessionHash string
}

func NewClient(sessionHash string) *client {
	rv := new(client)
	rv.sessionHash = sessionHash

	return rv
}

func (cl *client) Write(p []byte) (n int, err error) {
	err = app.NC.Publish(fmt.Sprintf("node.%s", cl.sessionHash), p)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func (cl *client) Run() error {
	// MakeRaw put the terminal connected to the given file descriptor into raw
	// mode and returns the previous state of the terminal so that it can be
	// restored.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		println("1")
		return err
	}

	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	app.NC.Subscribe(fmt.Sprintf("admin.%s", cl.sessionHash), func(msg *nats.Msg) {
		fmt.Fprintf(os.Stdout, "%s", string(msg.Data))
	})

	cl.Write([]byte{'\n'})

	reader := bufio.NewReader(os.Stdin)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			println("2")
			return err
		}

		cl.Write([]byte{b})

		if b == '\003' {
			break
		}
	}

	return nil
}

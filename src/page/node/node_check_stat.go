package node

import (
	"database/sql"
	"farlogin/src/app"
	"farlogin/src/lib/util"
	"farlogin/src/page/node/node_lib"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func init() {
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%d", app.NatsUrl, app.NatsPort),
		nats.Timeout(10*time.Second), nats.Name("farlogin server"), nats.UserInfo("farlogin", "Pseudo12Login"))
	if err != nil {
		panic(err)
	}

	nc.Subscribe("hello", func(msg *nats.Msg) {
		updateStatus(string(msg.Data))
	})

	nc.Subscribe("farlogin.sessionUsed", func(msg *nats.Msg) {
		err := markSessionUsed(string(msg.Data))
		if err != nil {
			log.Println(err.Error)
		}
	})

	nc.Subscribe("farlogin.isValidHash", func(msg *nats.Msg) {
		var rv string
		if ok := isValidHash(string(msg.Data)); ok {
			rv = "true"
		} else {
			rv = "false"
		}

		err := nc.Publish(msg.Reply, []byte(rv))
		if err != nil {
			log.Printf(err.Error())
		}
	})
}

func updateStatus(name string) {
	tx, err := app.Db.Begin()
	if err != nil {
		log.Println(err.Error())
	}

	now := util.Now()
	sqlStr := `update node set
						lastSeen = ?
					where
						name = ?`

	_, err = tx.Exec(sqlStr, now, name)

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
	}

	tx.Commit()
}

func markSessionUsed(sessionHash string) error {
	tx, err := app.Db.Begin()
	if err != nil {
		return err
	}

	sqlStr := `update nodeSession set
						status = ?
					where
						sessionHash = ?`

	_, err = tx.Exec(sqlStr, "used", sessionHash)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func isValidHash(sessionHash string) bool {
	rec, err := node_lib.GetSessionByHash(sessionHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}

		log.Println(err.Error)
		return false
	}

	if rec.Status != "new" {
		return false
	}

	return true
}

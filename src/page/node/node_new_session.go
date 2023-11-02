package node

import (
	"database/sql"
	"net/http"

	"farlogin/src/app"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
	"farlogin/src/page/node/node_lib"

	"github.com/go-sql-driver/mysql"
)

func NewSession(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("node", "insert") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("nodeId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.ReadCargo()

	nodeId := ctx.Cargo.Int("nodeId")
	rec, err := node_lib.GetNodeRec(nodeId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning("Record could not be found.")
			ctx.Redirect(ctx.U("/node", "key", "pn"))
			return
		}

		panic(err)
	}

	if !node_lib.IsNodeOn(rec.LastSeen) {
		ctx.Msg.Warning("This node is offline.")
		ctx.Redirect(ctx.U("/node_display", "nodeId", "key", "pn"))
		return
	}

	if ctx.Req.Method != "GET" {
		ctx.Msg.Warning("Unknown method.")
		ctx.Redirect(ctx.U("/node", "key", "pn"))
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	now := util.Now()
	sessionHash := util.NewUUID()

	sqlStr := `insert nodeSession(nodeSessionId, nodeId, recordTime, sessionHash, status)
					values(null, ?, ?, ?, ?)`
	_, err = tx.Exec(sqlStr, nodeId, now, sessionHash, "new")

	if err != nil {
		tx.Rollback()
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1452 {
				ctx.Msg.Warning("Could not find parent record.")
				insertForm(ctx)
				return
			}
		}

		panic(err)
	}

	tx.Commit()

	ctx.Msg.Success("Session has been created.")
	ctx.Redirect(ctx.U("/node_display", "nodeId", "key", "pn"))
}

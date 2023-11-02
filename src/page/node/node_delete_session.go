package node

import (
	"database/sql"
	"net/http"

	"farlogin/src/app"
	"farlogin/src/lib/context"
	"farlogin/src/page/node/node_lib"
)

func DeleteSession(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("node", "update") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("nodeId", -1)
	ctx.Cargo.AddInt("sessionId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.ReadCargo()

	nodeId := ctx.Cargo.Int("nodeId")
	_, err := node_lib.GetNodeRec(nodeId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning("Record could not be found.")
			ctx.Redirect(ctx.U("/node", "key", "pn"))
			return
		}

		panic(err)
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sessionId := ctx.Cargo.Int("sessionId")

	sqlStr := `delete from
					nodeSession
				where
					nodeSessionId = ?`

	res, err := tx.Exec(sqlStr, sessionId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning("Record could not be found.")
		ctx.Redirect(ctx.U("/node_display", "nodeId", "key", "pn"))
		return
	}

	tx.Commit()

	ctx.Msg.Success("Session has been deleted.")
	ctx.Redirect(ctx.U("/node_display", "nodeId", "key", "pn"))
}

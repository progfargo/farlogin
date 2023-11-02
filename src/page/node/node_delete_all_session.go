package node

import (
	"net/http"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/content/left_menu"
	"farlogin/src/content/top_menu"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
	"farlogin/src/page/node/node_lib"
)

func DeleteAllSession(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("node", "delete") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("nodeId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddStr("confirm", "no")
	ctx.ReadCargo()

	nodeId := ctx.Cargo.Int("nodeId")
	count := node_lib.CountNodeSessionList(nodeId)
	if count == 0 {
		ctx.Msg.Warning("No session to delete.")
		ctx.Redirect(ctx.U("/node", "key", "pn"))
	}

	if ctx.Cargo.Str("confirm") != "yes" {
		deleteAllSessionConfirm(ctx)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `delete from
					nodeSession
				where
					nodeId = ?`

	_, err = tx.Exec(sqlStr, nodeId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	tx.Commit()

	ctx.Msg.Success("Sessions has been deleted.")
	ctx.Redirect(ctx.U("/node_display", "nodeId", "key", "pn"))
}

func deleteAllSessionConfirm(ctx *context.Ctx) {

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Nodes", "Delete Sessions"))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx.U("/node_display", "nodeId", "key", "pn")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"callout calloutError\">")
	buf.Add("<h4>Please confirm:</h4>")
	buf.Add("<p>Do you realy want to delete all sessions?</p>")
	buf.Add("</div>")
	buf.Add("</div>")

	ctx.Cargo.SetStr("confirm", "yes")
	urlStr := ctx.U("/node_delete_all_session", "nodeId", "key", "pn", "confirm")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"confirmCommand\">")
	buf.Add("<a href=\"%s\" class=\"button buttonError buttonSm\">Yes</a>", urlStr)
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Include(ctx)
	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "node")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

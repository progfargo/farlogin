package node

import (
	"database/sql"
	"net/http"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/content/left_menu"
	"farlogin/src/content/top_menu"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
	"farlogin/src/page/node/node_lib"
)

func Delete(rw http.ResponseWriter, req *http.Request) {
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
	rec, err := node_lib.GetNodeRec(nodeId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning("Record could not be found.")
			ctx.Redirect(ctx.U("/node", "key", "pn"))
			return
		}

		panic(err)
	}

	if ctx.Cargo.Str("confirm") != "yes" {
		deleteConfirm(ctx, rec)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `delete from
					node
				where
					nodeId = ?`

	res, err := tx.Exec(sqlStr, nodeId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning("Record could not be found.")
		ctx.Redirect(ctx.U("/node", "key", "pn"))
		return
	}

	tx.Commit()

	ctx.Msg.Success("Record has been deleted.")
	ctx.Redirect(ctx.U("/node", "key", "pn"))
}

func deleteConfirm(ctx *context.Ctx, rec *node_lib.NodeRec) {
	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Nodes", "Delete Record", rec.Name))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx.U("/node", "key", "pn")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<table>")
	buf.Add("<tbody>")

	name := util.ScrStr(rec.Name)
	exp := util.ScrStr(rec.Exp)

	buf.Add("<tr><th class=\"fixedMiddle\">Name:</th><td>%s</td></tr>", name)
	buf.Add("<tr><th>Explanation:</th><td>%s</td></tr>", exp)
	buf.Add("</tbody>")
	buf.Add("</table>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"callout calloutError\">")
	buf.Add("<h4>Please confirm:</h4>")
	buf.Add("<p>Do you realy want to delete this record?</p>")
	buf.Add("</div>")
	buf.Add("</div>")

	ctx.Cargo.SetStr("confirm", "yes")
	urlStr := ctx.U("/node_delete", "nodeId", "key", "pn", "confirm")

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

	str := "nodeDeletePage"
	ctx.AddHtml("pageName", &str)

	ctx.Render("default.html")
}

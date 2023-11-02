package node

import (
	"database/sql"
	"net/http"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/content/left_menu"
	"farlogin/src/content/node_menu"
	"farlogin/src/content/top_menu"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
	"farlogin/src/page/node/node_lib"

	"github.com/go-sql-driver/mysql"
)

func Update(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("node", "update") {
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

	if ctx.Req.Method == "GET" {
		updateForm(ctx, rec)
		return
	}

	name := ctx.Req.PostFormValue("name")
	exp := ctx.Req.PostFormValue("exp")

	if name == "" {
		ctx.Msg.Warning("You have left one or more fields empty.")
		updateForm(ctx, rec)
		return
	}

	if !node_lib.CheckName(name) {
		ctx.Msg.Warning("Node name can only include lowercase letters, digits and dots.").
			Add("Sample: texas.well106")
		updateForm(ctx, rec)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `update node set
						name = ?,
						exp = ?
					where
						nodeId = ?`

	res, err := tx.Exec(sqlStr, name, exp, nodeId)

	if err != nil {
		tx.Rollback()
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1062 {
				ctx.Msg.Warning("Duplicate record.")
				updateForm(ctx, rec)
				return
			} else if err.Number == 1452 {
				ctx.Msg.Warning("Could not find parent record.")
				updateForm(ctx, rec)
				return
			}
		}

		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning("You have not change the record.")
		updateForm(ctx, rec)
		return
	}

	tx.Commit()

	ctx.Msg.Success("Record has been saved.")
	ctx.Redirect(ctx.U("/node_update", "key", "pn", "nodeId"))
}

func updateForm(ctx *context.Ctx, rec *node_lib.NodeRec) {
	content.Include(ctx)

	var name, exp string

	if ctx.Req.Method == "POST" {
		name = ctx.Req.PostFormValue("name")
		exp = ctx.Req.PostFormValue("exp")
	} else {
		name = rec.Name
		exp = rec.Exp
	}

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Nodes", "Edit Record", rec.Name))
	buf.Add("</div>")

	buf.Add("<div class=\"col lg2\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx.U("/node", "key", "pn")))
	buf.Add("</div>")
	buf.Add("</div>")

	nodeMenu := node_menu.New("nodeId", "key", "pn")
	nodeMenu.Set(ctx, "node_update")

	buf.Add("<div class=\"col lg10\">")
	buf.Add(nodeMenu.Format(ctx))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/node_update", "nodeId", "key", "pn")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Name:</label>")
	buf.Add("<input type=\"text\" name=\"name\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"2\">", util.ScrStr(name))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Explanation:</label>")
	buf.Add("<input type=\"text\" name=\"exp\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"250\" tabindex=\"2\">", util.ScrStr(exp))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"10\">Submit</button>")
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"11\">Reset</button>")
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "node")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	str := "nodeInsertPage"
	ctx.AddHtml("pageName", &str)

	ctx.Render("default.html")
}

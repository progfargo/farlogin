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

	"github.com/go-sql-driver/mysql"
)

func Insert(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("node", "insert") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("nodeId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.ReadCargo()

	if ctx.Req.Method == "GET" {
		insertForm(ctx)
		return
	}

	name := ctx.Req.PostFormValue("name")
	exp := ctx.Req.PostFormValue("exp")

	if name == "" {
		ctx.Msg.Warning("You have left one or more fields empty.")
		insertForm(ctx)
		return
	}

	if !node_lib.CheckName(name) {
		ctx.Msg.Warning("Node name can only include lowercase letters, digits and dots.")
		ctx.Msg.Warning("Sample texas.well106")
		insertForm(ctx)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `insert node(nodeId, name, exp)
					values(null, ?, ?)`
	_, err = tx.Exec(sqlStr, name, exp)

	if err != nil {
		tx.Rollback()
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1062 {
				ctx.Msg.Warning("Duplicate record.")
				insertForm(ctx)
				return
			} else if err.Number == 1452 {
				ctx.Msg.Warning("Could not find parent record.")
				insertForm(ctx)
				return
			}
		}

		panic(err)
	}

	tx.Commit()

	ctx.Msg.Success("Record has been saved.")
	ctx.Redirect(ctx.U("/node", "key", "pn"))
}

func insertForm(ctx *context.Ctx) {
	content.Include(ctx)

	var name, exp string

	if ctx.Req.Method == "POST" {
		name = ctx.Req.PostFormValue("name")
		exp = ctx.Req.PostFormValue("exp")
	} else {
		name = ""
		exp = ""
	}

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Nodes", "New Record"))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx.U("/node", "key", "pn")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/node_insert", "key", "pn", "pn")
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

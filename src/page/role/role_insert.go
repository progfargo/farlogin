package role

import (
	"net/http"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/content/left_menu"
	"farlogin/src/content/top_menu"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"

	"github.com/go-sql-driver/mysql"
)

func Insert(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("role", "insert") {
		app.BadRequest()
	}

	ctx.ReadCargo()

	if ctx.Req.Method == "GET" {
		insertForm(ctx)
		return
	}

	name := ctx.Req.PostFormValue("name")
	exp := ctx.Req.PostFormValue("exp")

	if name == "" || exp == "" {
		ctx.Msg.Warning("You have left one or more fields empty.")
		insertForm(ctx)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `insert into
					role(roleId, name, exp)
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
	ctx.Redirect(ctx.U("/role"))
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
	buf.Add(content.PageTitle("Roles", "New Record"))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx.U("/role")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/role_insert")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Role Name:</label>")
	buf.Add("<input type=\"text\" name=\"name\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"1\" autofocus>", util.ScrStr(name))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Explanation:</label>")
	buf.Add("<input type=\"text\" name=\"exp\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"250\" tabindex=\"2\">", util.ScrStr(exp))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"3\">Submit</button>")
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"4\">Reset</button>")
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "role")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

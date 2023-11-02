package role

import (
	"database/sql"
	"net/http"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/content/left_menu"
	"farlogin/src/content/top_menu"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
	"farlogin/src/page/role/role_lib"

	"github.com/go-sql-driver/mysql"
)

func Update(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("role", "update") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("roleId", -1)
	ctx.ReadCargo()

	roleId := ctx.Cargo.Int("roleId")
	rec, err := role_lib.GetRoleRec(roleId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning("Record could not be found.")
			ctx.Redirect(ctx.U("/role"))
			return
		}

		panic(err)
	}

	if rec.Name == "admin" && !ctx.IsSuperuser() {
		ctx.Msg.Warning("'admin' role can only be updated by superuser.")
		ctx.Redirect(ctx.U("/role"))
		return
	}

	if ctx.Req.Method == "GET" {
		updateForm(ctx, rec)
		return
	}

	name := ctx.Req.PostFormValue("name")
	exp := ctx.Req.PostFormValue("exp")

	if name == "" || exp == "" {
		ctx.Msg.Warning("You have left one or more fields empty.")
		updateForm(ctx, rec)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `update role set
					name = ?,
					exp = ?
				where
					roleId = ?`

	res, err := tx.Exec(sqlStr, name, exp, roleId)
	if err != nil {
		tx.Rollback()
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1062 {
				ctx.Msg.Warning("Duplicate record.")
				updateForm(ctx, rec)
				return
			}
		}

		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning("You did not change the record.")
		updateForm(ctx, rec)
		return
	}

	tx.Commit()

	ctx.Msg.Success("Record has been changed.")
	ctx.Redirect(ctx.U("/role"))
}

func updateForm(ctx *context.Ctx, rec *role_lib.RoleRec) {
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
	buf.Add(content.PageTitle("Roles", "Update Record"))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx.U("/role")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/role_update", "roleId")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Role Name:</label>")
	buf.Add("<input type=\"text\" name=\"name\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"1\" autofocus>", util.ScrStr(name))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Exlanation</label>")
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

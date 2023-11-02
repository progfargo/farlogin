package user

import (
	"database/sql"
	"net/http"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/content/combo"
	"farlogin/src/content/left_menu"
	"farlogin/src/content/top_menu"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
	"farlogin/src/page/user/user_lib"

	"github.com/go-sql-driver/mysql"
)

func Update(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("user", "update") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("userId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddInt("rid", -1)
	ctx.Cargo.AddStr("stat", "default")
	ctx.ReadCargo()

	userId := ctx.Cargo.Int("userId")
	rec, err := user_lib.GetUserRec(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning("Record could not be found.")
			ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
			return
		}

		panic(err)
	}

	if rec.Login == "superuser" {
		ctx.Msg.Warning("'superuser' account can not be updated.")
		ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
		return
	}

	if rec.Login == "testuser" && !ctx.IsSuperuser() {
		ctx.Msg.Warning("'testuser' account can not be updated.")
		ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
		return
	}

	if ctx.Req.Method == "GET" {
		updateForm(ctx, rec)
		return
	}

	name := ctx.Req.PostFormValue("name")
	email := ctx.Req.PostFormValue("email")
	status := ctx.Req.PostFormValue("status")

	if name == "" || email == "" || status == "default" {
		ctx.Msg.Warning("You have left one or more fields empty.")
		updateForm(ctx, rec)
		return
	}

	if err := util.IsValidEmail(email); err != nil {
		ctx.Msg.Warning(err.Error())
		updateForm(ctx, rec)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	var sqlStr string
	if email != rec.Email {
		sqlStr = `update user set
					user.name = ?,
					user.email = ?,
					user.status = ?,
					user.isEmailValidated = 'no'
				where
					user.userId = ?`
	} else {
		sqlStr = `update user set
					user.name = ?,
					user.email = ?,
					user.status = ?
				where
					user.userId = ?`
	}

	res, err := tx.Exec(sqlStr, name, email, status, userId)
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
	ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
}

func updateForm(ctx *context.Ctx, rec *user_lib.UserRec) {
	content.Include(ctx)

	var name, login, email, status string
	if ctx.Req.Method == "POST" {
		name = ctx.Req.PostFormValue("name")
		email = ctx.Req.PostFormValue("email")
		status = ctx.Req.PostFormValue("status")
	} else {
		name = rec.Name
		email = rec.Email
		status = rec.Status
	}

	login = rec.Login

	statusCombo := combo.NewEnumCombo()
	statusCombo.Add("default", "Select User Status")
	statusCombo.Add("active", "Active")
	statusCombo.Add("blocked", "Blocked")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Users", "Update Record"))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx.U("/user", "key", "pn", "rid", "stat")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/user_update", "userId", "key", "pn", "rid", "stat")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">User Name:</label>")
	buf.Add("<input type=\"text\" name=\"name\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"1\" autofocus>", util.ScrStr(name))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Login Name:</label>")
	buf.Add("<input type=\"text\" name=\"login\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"2\" disabled>", util.ScrStr(login))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Email:</label>")
	buf.Add("<input type=\"text\" name=\"email\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"3\">", util.ScrStr(email))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Status:</label>")
	buf.Add("<select name=\"status\" class=\"formControl\" tabindex=\"4\">")

	buf.Add(statusCombo.Format(status))

	buf.Add("</select>")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"5\">Submit</button>")
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"6\">Reset</button>")
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "user")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

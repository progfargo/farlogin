package user

import (
	"database/sql"
	"net/http"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/content/left_menu"
	"farlogin/src/content/top_menu"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
	"farlogin/src/page/user/user_lib"
)

func UpdatePass(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("user", "update_pass") {
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
		updatePassForm(ctx, rec)
		return
	}

	password := ctx.Req.PostFormValue("password")

	if password == "" {
		ctx.Msg.Warning("You have left one or more fields empty.")
		updatePassForm(ctx, rec)
		return
	}

	if err := util.IsValidPassword(password); err != nil {
		ctx.Msg.Warning(err.Error())
		updatePassForm(ctx, rec)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	password = util.PasswordHash(password)

	sqlStr := `update user set
					password = ?
				where
					userId = ?`

	res, err := tx.Exec(sqlStr, password, userId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning("You did not change the record.")
		updatePassForm(ctx, rec)
		return
	}

	tx.Commit()

	ctx.Msg.Success("User password has been changed.")
	ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
}

func updatePassForm(ctx *context.Ctx, rec *user_lib.UserRec) {
	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Users", "Change Password"))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx.U("/user", "key", "pn", "rid", "stat")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/user_update_pass", "userId", "key", "pn", "rid", "stat")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">New Password:</label>")
	buf.Add("<input type=\"password\" name=\"password\" class=\"formControl\"" +
		" value=\"\" maxlength=\"30\" tabindex=\"1\">")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"2\">Submit</button>")
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"3\">Reset</button>")
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Include(ctx)
	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "user")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

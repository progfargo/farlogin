package profile

import (
	"database/sql"
	"net/http"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/content/left_menu"
	"farlogin/src/content/profile_menu"
	"farlogin/src/content/top_menu"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
	"farlogin/src/page/user/user_lib"
)

func UpdatePass(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("profile", "update_password") {
		app.BadRequest()
	}

	ctx.ReadCargo()

	rec, err := user_lib.GetUserRec(ctx.User.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			panic("User record could not be found.")
		}

		panic(err)
	}

	if ctx.Req.Method == "GET" {
		updatePassForm(ctx, rec)
		return
	}

	curPassword := util.LimitStr(ctx.Req.PostFormValue("curPassword"), 30)
	newPassword := util.LimitStr(ctx.Req.PostFormValue("newPassword"), 30)
	reNewPassword := util.LimitStr(ctx.Req.PostFormValue("reNewPassword"), 30)

	if curPassword == "" || newPassword == "" || reNewPassword == "" {
		ctx.Msg.Warning("You have left one or more fields empty.")
		updatePassForm(ctx, rec)
		return
	}

	if util.PasswordHash(curPassword) != rec.Password {
		ctx.Msg.Warning("You have entered wrong current password. Please try again.")
		updatePassForm(ctx, rec)
		return
	}

	if err := util.IsValidPassword(newPassword); err != nil {
		ctx.Msg.Warning(err.Error())
		updatePassForm(ctx, rec)
		return
	}

	if newPassword == curPassword {
		ctx.Msg.Warning("You have entered your old password as new password.")
		updatePassForm(ctx, rec)
		return
	}

	if newPassword != reNewPassword {
		ctx.Msg.Warning("New password and retyped new password mismatch.")
		updatePassForm(ctx, rec)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	newPassword = util.PasswordHash(newPassword)
	sqlStr := `update user set
					password = ?
				where
					userId = ?`

	_, err = tx.Exec(sqlStr, newPassword, ctx.User.UserId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	tx.Commit()

	ctx.Msg.Success("Record has been changed.")
	ctx.Redirect(ctx.U("/profile"))
}

func updatePassForm(ctx *context.Ctx, rec *user_lib.UserRec) {

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Change Password"))
	buf.Add("</div>")

	profileMenu := profile_menu.New()
	profileMenu.Set(ctx, "update_password")

	buf.Add("<div class=\"col\">")
	buf.Add(profileMenu.Format(ctx))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/profile_update_pass")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Current Password:</label>")
	buf.Add("<input type=\"password\" name=\"curPassword\"" +
		" class=\"formControl\" value=\"\" maxlength=\"30\" tabindex=\"1\" autofocus>")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">New Password:</label>")
	buf.Add("<input type=\"password\" name=\"newPassword\"" +
		" class=\"formControl\" value=\"\" maxlength=\"30\" tabindex=\"2\">")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Retype New Password:</label>")
	buf.Add("<input type=\"password\" name=\"reNewPassword\"" +
		" class=\"formControl\" value=\"\" maxlength=\"30\" tabindex=\"3\">")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"4\">Submit</button>")
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"5\">Reset</button>")
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Include(ctx)
	ctx.Css.Add("/asset/css/page/profile.css")

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx)

	tmenu := top_menu.New()
	tmenu.Set(ctx, "profile")

	str := "profilePage"
	ctx.AddHtml("pageName", &str)

	ctx.Render("default.html")
}

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

	"github.com/go-sql-driver/mysql"
)

func UpdateProfile(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("profile", "update") {
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
		updateProfileForm(ctx, rec)
		return
	}

	email := util.LimitStr(ctx.Req.PostFormValue("email"), 100)
	name := util.LimitStr(ctx.Req.PostFormValue("name"), 100)
	password := util.LimitStr(ctx.Req.PostFormValue("password"), 30)

	if email == "" || name == "" || password == "" {
		ctx.Msg.Warning("You have left one or more fields empty.")
		updateProfileForm(ctx, rec)
		return
	}

	if util.PasswordHash(password) != rec.Password {
		ctx.Msg.Error("Wrong password. Please enter your current password to confirm profil update.")
		updateProfileForm(ctx, rec)
		return
	}

	if err := util.IsValidEmail(email); err != nil {
		ctx.Msg.Warning(err.Error())
		updateProfileForm(ctx, rec)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `update user set
					email = ?,
					name = ?
				where
					userId = ?`

	res, err := tx.Exec(sqlStr, email, name, ctx.User.UserId)
	if err != nil {
		tx.Rollback()

		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1062 {
				ctx.Msg.Warning("This email address is being used by another user.")
				updateProfileForm(ctx, rec)
				return
			}
		}

		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning("You did not change the record.")
		updateProfileForm(ctx, rec)
		return
	}

	tx.Commit()

	ctx.Msg.Success("Record has been changed.")
	ctx.Redirect(ctx.U("/profile"))
}

func updateProfileForm(ctx *context.Ctx, rec *user_lib.UserRec) {

	var email, name, password string

	if ctx.Req.Method == "POST" {
		email = ctx.Req.PostFormValue("email")
		name = ctx.Req.PostFormValue("name")
	} else {
		email = rec.Email
		name = rec.Name
	}

	password = ""

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Profile", "Update"))
	buf.Add("</div>")

	profileMenu := profile_menu.New()
	profileMenu.Set(ctx, "profile")

	buf.Add("<div class=\"col\">")
	buf.Add(profileMenu.Format(ctx))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/profile_update")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label>Login Name:</label>")
	buf.Add("<input type=\"text\" name=\"login\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" disabled>", util.ScrStr(rec.Login))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">E-Mail:</label>")
	buf.Add("<input type=\"text\" name=\"email\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"1\" autofocus>", util.ScrStr(email))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">User Name:</label>")
	buf.Add("<input type=\"text\" name=\"name\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"2\">", util.ScrStr(name))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Current Password:</label>")
	buf.Add("<input type=\"password\" name=\"password\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"30\" tabindex=\"3\">", util.ScrStr(password))
	buf.Add("<span class=\"helpBlock\">Please confirm update by entering your current password.</span>")
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

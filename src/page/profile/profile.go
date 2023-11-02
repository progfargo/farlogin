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

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("profile", "browse") {
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

	content.Include(ctx)
	ctx.Css.Add("/asset/css/page/profile.css")

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx)

	tmenu := top_menu.New()
	tmenu.Set(ctx, "profile")

	displayProfile(ctx, rec)

	str := "profilePage"
	ctx.AddHtml("pageName", &str)

	ctx.Render("default.html")
}

func displayProfile(ctx *context.Ctx, rec *user_lib.UserRec) {
	buf := util.NewBuf()

	updateRight := ctx.IsRight("profile", "update")

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Profile"))
	buf.Add("</div>")

	profileMenu := profile_menu.New()
	profileMenu.Set(ctx, "profile")

	buf.Add("<div class=\"col\">")
	buf.Add(profileMenu.Format(ctx))
	buf.Add("</div>")

	var urlStr string
	if updateRight {
		buf.Add("<div class=\"col\">")
		buf.Add("<div class=\"buttonGroupFixed\">")

		urlStr = ctx.U("/profile_update")
		buf.Add("<a href=\"%s\" class=\"button buttonWarning buttonSm\" title=\"Update profile info.\">Update Profile</a>", urlStr)

		buf.Add("</div>")
		buf.Add("</div>")
	}

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<tbody>")

	buf.Add("<tr><th class=\"fixedMiddle\">User Name:</th><td>%s</td></tr>", util.ScrStr(rec.Name))
	buf.Add("<tr><th>Login Name:</th><td>%s</td></tr>", util.ScrStr(rec.Login))
	buf.Add("<tr><th>E-Mail:</th><td>%s</td></tr>", util.ScrStr(rec.Email))

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}

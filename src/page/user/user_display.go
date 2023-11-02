package user

import (
	"database/sql"
	"net/http"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/content/left_menu"
	"farlogin/src/content/top_menu"
	"farlogin/src/lib/context"
	"farlogin/src/lib/user"
	"farlogin/src/lib/util"
	"farlogin/src/page/user/user_lib"
)

func Display(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("user", "browse") {
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

	content.Include(ctx)
	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "user")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	displayUserRec(ctx, rec)

	ctx.Render("default.html")
}

func displayUserRec(ctx *context.Ctx, rec *user_lib.UserRec) {

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Users", "Display Record"))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx.U("/user", "key", "pn", "rid", "stat")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<caption>User Information:</caption>")
	buf.Add("<tbody>")
	buf.Add("<tr><th class=\"fixedMiddle\">User Name:</th><td>%s</td></tr>", util.ScrStr(rec.Name))
	buf.Add("<tr><th>Login Name:</th><td>%s</td></tr>", util.ScrStr(rec.Login))
	buf.Add("<tr><th>Email:</th><td>%s</td></tr>", util.ScrStr(rec.Email))
	buf.Add("<tr><th>Status:</th><td>%s</td></tr>", user_lib.StatusToLabel(rec.Status))

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<caption>User Roles:</caption>")
	buf.Add("<tbody>")

	libUserRec, err := user.New(rec.UserId)
	if err != nil {
		ctx.Msg.Error(err.Error())
	}

	for k, _ := range libUserRec.RoleList {
		buf.Add("<tr>")
		buf.Add("<td>%s</td>", k)
		buf.Add("</tr>")
	}

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}

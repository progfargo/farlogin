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

func Delete(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("user", "delete") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("userId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddStr("confirm", "no")
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
		ctx.Msg.Warning("'superuser' account can not be deleted.")
		ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
		return
	}

	if rec.Login == "testuser" && !ctx.IsSuperuser() {
		ctx.Msg.Warning("'testuser' account can not be deleted.")
		ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
		return
	}

	if user_lib.IsUserRoleExists(userId, "admin") {
		n := user_lib.CountAdmin()
		if n < 2 {
			ctx.Msg.Warning("Last admin user can not be deleted.")
			ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
			return
		}
	}

	if ctx.Cargo.Str("confirm") != "yes" {
		deleteConfirm(ctx, rec)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `delete from
					user
				where
					userId = ?`

	res, err := tx.Exec(sqlStr, userId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning("Record could not be found.")
		ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
		return
	}

	tx.Commit()

	ctx.Msg.Success("Record has been deleted.")
	ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
}

func deleteConfirm(ctx *context.Ctx, rec *user_lib.UserRec) {
	content.Include(ctx)

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Users", "Delete Record"))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx.U("/user", "key", "pn", "rid", "stat")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<table>")
	buf.Add("<tbody>")
	buf.Add("<tr><th class=\"fixedMiddle\">User Name:</th><td>%s</td></tr>", util.ScrStr(rec.Name))
	buf.Add("<tr><th>Login Name:</th><td>%s</td></tr>", util.ScrStr(rec.Login))
	buf.Add("<tr><th>Email:</th><td>%s</td></tr>", util.ScrStr(rec.Email))
	buf.Add("<tr><th>Status:</th><td>%s</td></tr>", user_lib.StatusToLabel(rec.Status))
	buf.Add("</tbody>")
	buf.Add("</table>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"callout calloutError\">")
	buf.Add("<h4>Please confirm:</h4>")
	buf.Add("<p>Do you realy want to delete this record?</p>")
	buf.Add("</div>")
	buf.Add("</div>")

	ctx.Cargo.SetStr("confirm", "yes")
	urlStr := ctx.U("/user_delete", "userId", "key", "pn", "confirm", "rid", "stat")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"confirmCommand\">")
	buf.Add("<a href=\"%s\" class=\"button buttonError buttonSm\">Yes</a>", urlStr)
	buf.Add("</div>")
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

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
	"farlogin/src/page/role/role_lib"
	"farlogin/src/page/user/user_lib"
)

func UserRole(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("user", "role_browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("userId", -1)
	ctx.Cargo.AddInt("roleId", -1)
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

	if rec.Login == "superuser" && !ctx.IsSuperuser() {
		ctx.Msg.Warning("'superuser' account can not be updated.")
		ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
		return
	}

	displayUser(ctx, rec)
}

func displayUser(ctx *context.Ctx, userRec *user_lib.UserRec) {

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")

	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Users", "User Roles"))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx.U("/user", "key", "pn", "rid", "stat")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	//user info
	buf.Add("<table>")
	buf.Add("<caption>User Information:</caption>")
	buf.Add("<tbody>")
	buf.Add("<tr><th class=\"fixedMiddle\">User Name:</th><td>%s</td></tr>", util.ScrStr(userRec.Name))
	buf.Add("<tr><th>Login:</th><td>%s</td></tr>", util.ScrStr(userRec.Login))
	buf.Add("<tr><th>Email:</th><td>%s</td></tr>", util.ScrStr(userRec.Email))
	buf.Add("<tr><th>Status:</th><td>%s</td></tr>", user_lib.StatusToLabel(userRec.Status))

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	//grant role
	if ctx.IsRight("user", "role_grant") {
		sqlStr := `select
						roleId,
						name
					from
						role`

		roleCombo := combo.NewCombo(sqlStr, "Select Role")
		roleCombo.Set()
		if roleCombo.IsEmpty() {
			ctx.Msg.Warning("Role list is empty. You should enter at least one role first.")
			ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
			return
		}

		buf.Add("<table>")
		buf.Add("<caption>Grant Role:</caption>")
		buf.Add("<tbody>")
		buf.Add("<tr>")
		buf.Add("<th class=\"fixedMiddle\">Role Name:</th>")
		buf.Add("<td>")

		urlStr := ctx.U("/user_role_grant", "userId", "key", "pn", "rid", "stat")
		buf.Add("<form class=\"formFlex\" action=\"%s\" method=\"post\">", urlStr)

		buf.Add("<div class=\"formGroup\">")
		buf.Add("<select name=\"roleId\" class=\"formControl\">")

		buf.Add(roleCombo.Format(-1))

		buf.Add("</select>")
		buf.Add("</div>")

		buf.Add("<div class=\"formGroup\">")
		buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\">Grant</button>")
		buf.Add("</div>")

		buf.Add("</form>")

		buf.Add("</td>")
		buf.Add("</tr>")

		buf.Add("</tbody>")
		buf.Add("</table>")
	}

	buf.Add("</div>")

	revokeRight := ctx.IsRight("user", "role_revoke")

	buf.Add("<div class=\"col\">")

	//user role list
	buf.Add("<table>")
	buf.Add("<caption>User Role List:</caption>")
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th class=\"fixedMiddle\">Role Name</th>")
	buf.Add("<th>Explanation</th>")

	if revokeRight {
		buf.AddLater("<th class=\"right fixedZero\">Command</th>")
	}

	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	userRoleList, err := role_lib.GetRoleByUser(userRec.UserId)
	if err != nil {
		panic(err)
	}

	if len(userRoleList) != 0 {
		buf.Forge()

		for _, row := range userRoleList {
			buf.Add("<tr>")
			buf.Add("<td>%s</td><td>%s</td>", row.Name, row.Exp)

			if revokeRight {
				buf.Add("<td class=\"right\">")
				buf.Add("<div class=\"buttonGroupFixed\">")

				buf.Add("<div class=\"revoke\">")
				buf.Add("<a href=\"#\" class=\"button buttonError buttonXs revokeButton\">Revoke</a>")
				buf.Add("</div>")

				buf.Add("<div class=\"revokeConfirm\">")
				buf.Add("Do you realy want to revoke this role?")

				ctx.Cargo.SetInt("roleId", row.RoleId)
				urlStr := ctx.U("/user_role_revoke", "userId", "roleId", "key", "pn", "stat")
				buf.Add("<a href=\"%s\" class=\"button buttonSuccess buttonXs\" title=\"Revoke this role.\">Yes</a>", urlStr)

				buf.Add("<a class=\"button buttonDefault buttonXs cancelButton\">Cancel</a>")

				buf.Add("</div>") //revokeConfirm
				buf.Add("</div>") //buttonGroupFixed
				buf.Add("</td>")
			}

			buf.Add("</tr>")
		}
	} else {
		buf.Add("<tr><td colspan=\"2\"><span class=\"label label-danger label-sd\">Empty</span></td></tr>")
	}

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Include(ctx)
	ctx.Css.Add("/asset/css/page/user_role.css")
	ctx.Js.Add("/asset/js/page/user/user_role.js")

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "user")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	str := "userRolePage"
	ctx.AddHtml("pageName", &str)

	ctx.Render("default.html")
}

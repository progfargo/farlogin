package role

import (
	"net/http"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/content/left_menu"
	"farlogin/src/content/top_menu"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
	"farlogin/src/page/role/role_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("role", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("roleId", -1)
	ctx.ReadCargo()

	content.Include(ctx)

	browseMid(ctx)

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "role")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {
	ctx.Js.Add("/asset/js/page/role/role_right.js")

	insertRight := ctx.IsRight("role", "insert")
	updateRight := ctx.IsRight("role", "update")
	deleteRight := ctx.IsRight("role", "delete")
	roleRightRight := ctx.IsRight("role", "role_right")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Roles"))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")

	if insertRight {
		buf.Add(content.NewButton(ctx.U("/role_insert")))
	}

	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th class=\"fixedMiddle\">Role Name</th>")
	buf.Add("<th>Explanation:</th>")

	if updateRight || deleteRight || roleRightRight {
		buf.Add("<th class=\"right\">Command</th>")
	}

	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	roleList := role_lib.GetRoleList()
	if len(roleList) == 0 {
		ctx.Msg.Warning("Empty list.")
	}

	var name, exp string
	for _, row := range roleList {
		ctx.Cargo.SetInt("roleId", row.RoleId)

		name = util.ScrStr(row.Name)
		exp = util.ScrStr(row.Exp)

		buf.Add("<tr>")
		buf.Add("<td>%s</td>", name)
		buf.Add("<td>%s</td>", exp)

		if updateRight || deleteRight || roleRightRight {
			buf.Add("<td class=\"right fixedSmall\">")
			buf.Add("<div class=\"buttonGroupFixed\">")

			if updateRight {
				urlStr := ctx.U("/role_update", "roleId")
				buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"Edit record.\">Edit</a>", urlStr)
			}

			if updateRight {
				urlStr := ctx.U("/role_right", "roleId")
				buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"Edit role rights.\">Role Rights</a>", urlStr)
			}

			if deleteRight {
				urlStr := ctx.U("/role_delete", "roleId")
				buf.Add("<a href=\"%s\" class=\"button buttonError buttonXs\" title=\"Delete record.\">Delete</a>", urlStr)
			}

			buf.Add("</div>")
			buf.Add("</td>")
		}

		buf.Add("</tr>")
	}

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}

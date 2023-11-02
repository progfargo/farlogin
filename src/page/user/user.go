package user

import (
	"net/http"
	"strings"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/content/combo"
	"farlogin/src/content/left_menu"
	"farlogin/src/content/ruler"
	"farlogin/src/content/top_menu"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
	"farlogin/src/page/user/user_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
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

	content.Include(ctx)

	browseMid(ctx)

	content.Default(ctx)

	content.Search(ctx, ctx.U("/user"))

	lmenu := left_menu.New()
	lmenu.Set(ctx, "user")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {
	ctx.Js.Add("/asset/js/page/user/user.js")

	key := ctx.Cargo.Str("key")
	pageNo := ctx.Cargo.Int("pn")
	rid := ctx.Cargo.Int("rid")
	stat := ctx.Cargo.Str("stat")

	sqlStr := `select
					roleId,
					name
				from
					role
				order by name`

	roleCombo := combo.NewCombo(sqlStr, "All Roles")
	roleCombo.Set()

	if roleCombo.IsEmpty() {
		ctx.Msg.Warning("Role list is empty. You should enter at least one role first.")
	}

	statCombo := combo.NewEnumCombo()
	statCombo.Add("default", "All Statuses")
	statCombo.Add("active", "Active")
	statCombo.Add("blocked", "Blocked")

	totalRows := user_lib.CountUser(key, rid, stat)
	if totalRows == 0 {
		ctx.Msg.Warning("Empty list.")
	}

	pageLen := ctx.Config.Int("pageLen")
	pageNo = ctx.TouchPageNo(pageNo, totalRows, pageLen)

	insertRight := ctx.IsRight("user", "insert")
	updateRight := ctx.IsRight("user", "update")
	updatePassRight := ctx.IsRight("user", "update_pass")
	deleteRight := ctx.IsRight("user", "delete")
	roleBrowseRight := ctx.IsRight("user", "role_browse")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Users"))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")

	if insertRight {
		buf.Add(content.NewButton(ctx.U("/user_insert", "key", "pn", "rid", "stat")))
	}

	//rid form
	buf.Add("<form id=\"ridForm\" class=\"formInline\" action=\"/user\" method=\"get\">")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<select name=\"rid\" class=\"formControl\">")

	buf.Add(roleCombo.Format(rid))

	buf.Add("</select>")
	buf.Add("</div>")
	buf.Add(content.HiddenCargo(ctx, "key", "stat"))
	buf.Add("</form>")

	//stat form
	buf.Add("<form id=\"statForm\" class=\"formInline\" action=\"/user\" method=\"get\">")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<select name=\"stat\" class=\"formControl\">")

	buf.Add(statCombo.Format(stat))

	buf.Add("</select>")
	buf.Add("</div>")
	buf.Add(content.HiddenCargo(ctx, "key", "rid"))
	buf.Add("</form>")

	urlStr := ctx.U("/user")
	buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonSm\" title=\"Reset all filters.\">Clear</a>", urlStr)

	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")

	buf.Add("<thead>")
	buf.Add("<tr>")

	buf.Add("<th>Name</th>")
	buf.Add("<th>Login</th>")
	buf.Add("<th>Email</th>")
	buf.Add("<th>Status</th>")
	buf.Add("<th>Roles</th>")

	if updateRight || roleBrowseRight || updatePassRight || deleteRight {
		buf.Add("<th class=\"right\">Command</th>")
	}

	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	var userRoleList []string

	if totalRows > 0 {
		userList := user_lib.GetUserPage(ctx, key, pageNo, rid, stat)

		var name, login, email string
		for _, row := range userList {
			if row.Login == "superuser" && !ctx.IsSuperuser() {
				continue
			}

			ctx.Cargo.SetInt("userId", row.UserId)

			name = util.ScrStr(row.Name)
			login = util.ScrStr(row.Login)
			email = util.ScrStr(row.Email)

			if key != "" {
				name = content.Find(name, key)
				login = content.Find(login, key)
				email = content.Find(email, key)
			}

			buf.Add("<tr>")

			urlStr := ctx.U("/user_display", "userId")
			buf.Add("<td><a href=\"%s\">%s</a></td>", urlStr, name)
			buf.Add("<td>%s</td>", login)
			buf.Add("<td>%s</td>", email)

			buf.Add("<td class=\"center\">%s</td>", user_lib.StatusToLabel(row.Status))

			userRoleList = user_lib.GetUserRoleList(row.UserId)

			buf.Add("<td>%s</td>", strings.Join(userRoleList, ", "))

			if updateRight || roleBrowseRight || updatePassRight || deleteRight {
				buf.Add("<td class=\"right\">")
				buf.Add("<div class=\"buttonGroupFlex\">")

				if updateRight {
					urlStr = ctx.U("/user_update", "userId", "key", "pn", "rid", "stat")
					buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\""+
						" title=\"Edit record.\">Edit</a>", urlStr)
				}

				if roleBrowseRight {
					urlStr = ctx.U("/user_role", "userId", "key", "pn", "rid", "stat")
					buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\""+
						" title=\"Edit user roles.\">Roles</a>", urlStr)
				}

				if updatePassRight {
					urlStr = ctx.U("/user_update_pass", "userId", "key", "pn", "rid", "stat")
					buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\""+
						" title=\"Change user password.\">Change Pass</a>", urlStr)
				}

				if deleteRight {
					urlStr = ctx.U("/user_delete", "userId", "key", "pn", "rid", "stat")
					buf.Add("<a href=\"%s\" class=\"button buttonError buttonXs\""+
						" title=\"Delete record.\">Delete</a>", urlStr)
				}

				buf.Add("</div>")
				buf.Add("</td>")
			}

			buf.Add("</tr>")
		}
	}

	buf.Add("</tbody>")
	buf.Add("</table>")
	buf.Add("</div>")

	totalPage := ctx.TotalPage(totalRows, pageLen)
	if totalPage > 1 {
		buf.Add("<div class=\"col\">")
		ruler := ruler.NewRuler(totalPage, pageNo, ctx.U("/user", "key", "rid", "stat"))
		ruler.Set(ctx)
		buf.Add(ruler.Format())
		buf.Add("</div>")
	}

	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}

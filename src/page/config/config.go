package config

import (
	"net/http"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/content/left_menu"
	"farlogin/src/content/top_menu"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
	"farlogin/src/page/config/config_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("config", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("configId", -1)
	ctx.Cargo.AddInt("gid", -1)
	ctx.ReadCargo()

	content.Include(ctx)
	content.Default(ctx)

	browseMid(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "config")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {

	configList := config_lib.GetConfigList()
	if len(configList) == 0 {
		ctx.Msg.Warning("Empty list.")
	}

	insertRight := ctx.IsRight("config", "insert")
	updateRight := ctx.IsRight("config", "update")
	deleteRight := ctx.IsRight("config", "delete")
	setRight := ctx.IsRight("config", "set")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Configuration"))
	buf.Add("</div>")

	if insertRight {
		buf.Add("<div class=\"col\">")
		buf.Add("<div class=\"buttonGroupFixed\">")

		buf.Add(content.NewButton(ctx.U("/config_insert")))

		buf.Add("</div>")
		buf.Add("</div>")
	}

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th>Title</th>")
	buf.Add("<th>Value</th>")

	if updateRight || deleteRight {
		buf.Add("<th>Enum</th>")
	}

	buf.Add("<th>Group Name</th>")
	buf.Add("<th>Explanation</th>")

	if setRight || updateRight || deleteRight {
		buf.Add("<th class=\"right\">Command</th>")
	}

	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	var name, title, value, groupName, exp string
	for _, row := range configList {
		ctx.Cargo.SetInt("configId", row.ConfigId)

		name = util.ScrStr(row.Name)
		title = util.ScrStr(row.Title)
		value = util.ScrStr(row.Value)
		exp = util.ScrStr(row.Exp)

		buf.Add("<tr>")
		buf.Add("<td title=\"%s\">%s</td>", name, title)
		buf.Add("<td>%s</td>", value)

		if updateRight || deleteRight {
			buf.Add("<td>%d</td>", row.Enum)
		}

		buf.Add("<td>%s</td>", groupName)
		buf.Add("<td>%s</td>", exp)

		if setRight || updateRight || deleteRight {
			buf.Add("<td class=\"right\">")
			buf.Add("<div class=\"buttonGroupFixed\">")

			if setRight {
				urlStr := ctx.U("/config_set", "configId", "gid")
				buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"Set value.\">Set</a>", urlStr)
			}

			if updateRight {
				urlStr := ctx.U("/config_update", "configId", "gid")
				buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"Edit record.\">Edit</a>", urlStr)
			}

			if deleteRight {
				urlStr := ctx.U("/config_delete", "configId", "gid")
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

package node

import (
	"net/http"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/content/left_menu"
	"farlogin/src/content/ruler"
	"farlogin/src/content/top_menu"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
	"farlogin/src/page/node/node_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("node", "browse") && !ctx.IsRight("node", "browse_all") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("nodeId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.ReadCargo()

	content.Include(ctx)
	ctx.Js.Add("/asset/js/page/node/node.js")

	browseMid(ctx)

	content.Default(ctx)

	content.Search(ctx, "/node")

	lmenu := left_menu.New()
	lmenu.Set(ctx, "node")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	str := "nodePage"
	ctx.AddHtml("pageName", &str)

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {
	key := ctx.Cargo.Str("key")
	pageNo := ctx.Cargo.Int("pn")

	totalRows := node_lib.CountNode(key)
	if totalRows == 0 {
		ctx.Msg.Warning("Empty list.")
	}

	pageLen := ctx.Config.Int("pageLen")
	pageNo = ctx.TouchPageNo(pageNo, totalRows, pageLen)

	insertRight := ctx.IsRight("node", "insert")
	updateRight := ctx.IsRight("node", "update")
	deleteRight := ctx.IsRight("node", "delete")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Nodes"))
	buf.Add("</div>")

	var urlStr string
	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")

	if insertRight {
		buf.Add(content.NewButton(ctx.U("/node_insert", "key", "pn")))
	}

	buf.Add("</div>") //buttonGroupFixed
	buf.Add("</div>") //col

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th class=\"fixedZero\">Id</th>")

	buf.Add("<th>Name</th>")
	buf.Add("<th>Exp</th>")
	buf.Add("<th>Last Seen</th>")
	buf.Add("<th>Status</th>")
	buf.Add("<th class=\"right fixedZero\">Command</th>")
	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	if totalRows > 0 {
		nodeList := node_lib.GetNodePage(ctx, key, pageNo)

		var name string
		for _, row := range nodeList {
			ctx.Cargo.SetInt("nodeId", row.NodeId)

			name = util.ScrStr(row.Name)

			if key != "" {
				name = content.Find(name, key)
			}

			buf.Add("<tr>")
			buf.Add("<td>%d</td>", row.NodeId)

			buf.Add("<td>%s</td>", name)
			buf.Add("<td>%s</td>", row.Exp)
			buf.Add("<td>%s</td>", util.Int64ToTimeStr(row.LastSeen))
			buf.Add("<td class=\"center\">%s</td>", node_lib.NodeStatusToLabel(row.LastSeen))

			buf.Add("<td class=\"right\">")
			buf.Add("<div class=\"buttonGroupFixed\">")

			urlStr = ctx.U("/node_display", "nodeId", "key", "pn")
			buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"Display record.\">Display</a>", urlStr)

			if updateRight {
				urlStr = ctx.U("/node_update", "nodeId", "key", "pn")
				buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"Edit record.\">Edit</a>", urlStr)
			}

			if deleteRight {
				urlStr = ctx.U("/node_delete", "nodeId", "key", "pn")
				buf.Add("<a href=\"%s\" class=\"button buttonError buttonXs\" title=\"Delete record.\">Delete</a>", urlStr)
			}

			buf.Add("</div>")
			buf.Add("</td>")

			buf.Add("</tr>")
		}
	}

	buf.Add("</tbody>")
	buf.Add("</table>")
	buf.Add("</div>")

	totalPage := ctx.TotalPage(totalRows, pageLen)
	if totalPage > 1 {
		buf.Add("<div class=\"col\">")
		ruler := ruler.NewRuler(totalPage, pageNo, ctx.U("/node", "key"))
		ruler.Set(ctx)
		buf.Add(ruler.Format())
		buf.Add("</div>")
	}

	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}

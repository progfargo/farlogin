package node

import (
	"database/sql"
	"net/http"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/content/left_menu"
	"farlogin/src/content/node_menu"
	"farlogin/src/content/top_menu"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
	"farlogin/src/page/node/node_lib"
)

func Display(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("node", "browse_own") && !ctx.IsRight("node", "browse_all") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("nodeId", -1)
	ctx.Cargo.AddInt("sessionId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.ReadCargo()

	nodeId := ctx.Cargo.Int("nodeId")
	rec, err := node_lib.GetNodeRec(nodeId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning("Record could not be found.")
			ctx.Redirect(ctx.U("/node", "key", "pn"))
			return
		}

		panic(err)
	}

	displayNode(ctx, rec)
}

func displayNode(ctx *context.Ctx, rec *node_lib.NodeRec) {
	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Nodes", "Display Record", util.ScrStr(rec.Name)))
	buf.Add("</div>")

	buf.Add("<div class=\"col lg2\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx.U("/node", "key", "pn")))
	buf.Add("</div>")
	buf.Add("</div>")

	nodeMenu := node_menu.New("nodeId", "key", "pn")
	nodeMenu.Set(ctx, "node_display")

	buf.Add("<div class=\"col lg10\">")
	buf.Add(nodeMenu.Format(ctx))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<table>")
	buf.Add("<caption>Node Information:</caption>")
	buf.Add("<tbody>")

	name := util.ScrStr(rec.Name)
	exp := util.ScrStr(rec.Exp)

	buf.Add("<tr><th class=\"fixedMiddle\">Name:</th><td>%s</td></tr>", name)
	buf.Add("<tr><th>Explanation:</th><td>%s</td></tr>", exp)
	buf.Add("<tr><th>Status:</th><td>%s</td></tr>", node_lib.NodeStatusToLabel(rec.LastSeen))

	buf.Add("</tbody>")
	buf.Add("</table>")
	buf.Add("</div>") //col

	nodeSessionList := node_lib.GetNodeSessionList(rec.NodeId)
	if len(nodeSessionList) != 0 {
		buf.Add("<div class=\"col\">")
		buf.Add("<table>")
		buf.Add("<caption>Session List:</caption>")

		buf.Add("<thead>")
		buf.Add("<tr>")
		buf.Add("<th class=\"fixedMiddle\">Time</th>")
		buf.Add("<th>Session Key</th>")
		buf.Add("<th>Status</th>")
		buf.Add("<th class=\"right\">Command</th>")
		buf.Add("</tr>")
		buf.Add("</thead>")

		buf.Add("<tbody>")

		var urlStr string
		for _, v := range nodeSessionList {
			buf.Add("<tr>")
			buf.Add("<td>%s</td>", util.Int64ToTimeStr(v.RecordTime))
			buf.Add("<td class=\"sessionId\">")
			buf.Add("<span class=\"sessionHash\">%s</span>", v.SessionHash)
			buf.Add("<span class=\"copyLink\" title=\"Copy session hash.\"><i class=\"fa-solid fa-copy\"></i></span>")
			buf.Add("</td>")

			buf.Add("<td class=\"center\">%s</td>", node_lib.SessionStatusToLabel(v.Status))

			buf.Add("<td class=\"right\">")
			buf.Add("<div class=\"buttonGroupFixed\">")

			buf.Add("<div class=\"delete\">")
			buf.Add("<a href=\"#\" class=\"button buttonError buttonXs deleteButton\">Delete</a>")
			buf.Add("</div>")

			buf.Add("<div class=\"deleteConfirm\">")
			buf.Add("Do you realy want to delete this session?")

			ctx.Cargo.SetInt("sessionId", v.NodeSessionId)
			urlStr = ctx.U("/node_delete_session", "sessionId", "nodeId", "key", "pn")
			buf.Add("<a href=\"%s\" class=\"button buttonSuccess buttonXs\" title=\"Delete this session.\">Yes</a>", urlStr)
			buf.Add("<a class=\"button buttonDefault buttonXs cancelButton\">Cancel</a>")

			buf.Add("</div>") //deleteConfirm

			buf.Add("</div>")
			buf.Add("</td>")
			buf.Add("</tr>")
		}

		buf.Add("</tbody>")
		buf.Add("</table>")
		buf.Add("</div>") //col
	}

	buf.Add("</div>") //row

	ctx.AddHtml("midContent", buf.String())

	content.Include(ctx)
	ctx.Css.Add("/asset/css/page/node_display.css")
	ctx.Js.Add("/asset/js/page/node/node_display.js")

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "node")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	str := "nodeDisplayPage"
	ctx.AddHtml("pageName", &str)

	ctx.Render("default.html")
}

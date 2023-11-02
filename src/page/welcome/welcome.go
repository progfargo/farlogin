package welcome

import (
	"net/http"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/content/left_menu"
	"farlogin/src/content/top_menu"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsLoggedIn() {
		app.BadRequest()
	}

	ctx.ReadCargo()

	content.Include(ctx)

	browseMid(ctx)

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx)

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Welcome"))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"callout calloutSuccess\">")
	buf.Add("<h4>Welcome. You have logged in.</h4>")
	buf.Add("<p>You can manage nodes by using related links.</p>")
	buf.Add("</div>")
	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}

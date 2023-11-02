package content

import (
	"fmt"
	"regexp"
	"strings"

	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
)

func Include(ctx *context.Ctx) {
	ctx.Css.Add("/asset/fontawesome/css/all.css")
	ctx.Css.Add("/asset/css/style.css")

	ctx.Js.Add("/asset/js/jquery-2.1.3.min.js")
	ctx.Js.Add("/asset/js/flex_watch.js")
	ctx.Js.Add("/asset/js/common.js")
}

func Default(ctx *context.Ctx) {
	HtmlTitle(ctx, "farlogin")
	SiteTitle(ctx)

	if ctx.IsLoggedIn() {
		UserInfo(ctx)
	}

	Logo(ctx)
	FooterNote(ctx)
}

func HtmlTitle(ctx *context.Ctx, str string) {
	ctx.AddHtml("htmlTitle", &str)
}

func SiteTitle(ctx *context.Ctx) {
	link := "<a href=\"/node\">far<span>login</span></a>"
	str := fmt.Sprintf("<div class=\"siteTitle\"><h1>%s</h1></div>", link)
	ctx.AddHtml("siteTitle", &str)
}

func FooterNote(ctx *context.Ctx) {
	str := "<p><small>&copy; 2023 farlogin ver:1.0.3</small></p>"

	ctx.AddHtml("footerNote", &str)
}

func Search(ctx *context.Ctx, url string, args ...string) {
	buf := util.NewBuf()

	key := ctx.Cargo.Str("key")
	buf.Add("<form class=\"formInline formSearch\" action=\"%s\" method=\"get\">", url)
	buf.Add("<input type=\"text\" name=\"key\" placeholder=\"%s\" value=\"%s\">", "search word...", key)
	buf.Add(HiddenCargo(ctx, args...))
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\">%s</button>", "Search")
	buf.Add("</form>")

	ctx.AddHtml("search", buf.String())
}

func UserInfo(ctx *context.Ctx) {
	buf := util.NewBuf()

	var urlStr, link string
	if ctx.Session.UserId != ctx.Session.EffectiveUserId {
		urlStr = ctx.U("/user_unselect")
		link = fmt.Sprintf(" <a href=\"%s\"><i class=\"fas fa-window-close\"></i></a>", urlStr)
	}

	buf.Add("<div class=\"userInfo\">")

	buf.Add("<p><strong>User:</strong> <span title=\"%s\">%s</span>%s</p>", ctx.User.Name, ctx.User.Login, link)

	buf.Add("</div>")

	ctx.AddHtml("userInfo", buf.String())
}

func Logo(ctx *context.Ctx) {
	buf := util.NewBuf()

	buf.Add("<div class=\"logo\">")
	buf.Add("<img src=\"/asset/img/farlogin_logo.png\" alt=\"farlogin logo\">")
	buf.Add("</div>")

	ctx.AddHtml("logo", buf.String())
}

func HiddenCargo(ctx *context.Ctx, args ...string) string {
	buf := util.NewBuf()
	var val string

	for _, name := range args {
		if !ctx.Cargo.IsDefault(name) {
			val = ctx.Cargo.Str(name)
			buf.Add("<input type=\"hidden\" name=\"%s\" value=\"%s\">", name, val)
		}
	}

	return *buf.String()
}

func Find(str, key string) string {
	key = strings.ReplaceAll(key, "%", "")
	key = strings.ReplaceAll(key, "?", "")

	re := regexp.MustCompile("(?i)" + key)
	return re.ReplaceAllStringFunc(str, replaceFunc)
}

func replaceFunc(str string) string {
	return fmt.Sprintf("<span class=\"find\">%s</span>", str)
}

func PageTitle(args ...string) string {
	buf := util.NewBuf()
	buf.Add("<div class=\"pageTitle\">")
	buf.Add("<h3>%s</h3>", strings.Join(args, " <i class=\"fa fa-fw fa-angle-right\"></i> "))
	buf.Add("</div>")

	return *buf.String()
}

func BackButton(urlStr string) string {
	return fmt.Sprintf("<a href=\"%s\" class=\"button buttonAccent buttonSm\">Back</a>", urlStr)
}

func NewButton(urlStr string) string {
	return fmt.Sprintf("<a href=\"%s\" class=\"button buttonPrimary buttonSm\" title=\"New record.\">New</a>", urlStr)
}

func End(ctx *context.Ctx, calloutType, calloutTitle, htmlTitle string, msg ...string) {
	buf := util.NewBuf()

	buf.Add("<div class=\"callout callout%s\">", strings.Title(calloutType))
	buf.Add("<h4>%s</h4>", calloutTitle)

	for _, v := range msg {
		buf.Add("<p>%s</p>", v)
	}

	buf.Add("</div>")

	Include(ctx)
	ctx.Css.Add("/asset/css/page/end.css")
	ctx.AddHtml("midContent", buf.String())

	SiteTitle(ctx)
	HtmlTitle(ctx, htmlTitle)
	FooterNote(ctx)

	str := "endPage"
	ctx.AddHtml("pageName", &str)

	ctx.Render("end.html")
}

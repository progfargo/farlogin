package context

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"farlogin/src/app"
	"farlogin/src/content/css"
	"farlogin/src/content/js"
	"farlogin/src/lib/cargo"
	"farlogin/src/lib/user"
)

type Ctx struct {
	Config    app.ConfigType
	Css       *css.Css
	Js        js.Js
	Msg       *MessageList
	htmlOut   map[string]*string
	jsonOut   map[string]*string
	Rw        http.ResponseWriter
	Req       *http.Request
	SessionId string
	User      *user.UserRec
	RightList map[string]bool
	Session   *sessionType
	Cargo     cargo.CargoList
	Url       *url.URL
}

func NewContext(rw http.ResponseWriter, req *http.Request) *Ctx {
	ctx := new(Ctx)

	ctx.Rw = rw
	ctx.Req = req
	ctx.Config = app.CopyConfig()
	ctx.Cargo = cargo.NewCargo()

	parsedUrl, err := url.ParseRequestURI(ctx.Req.RequestURI)
	if err != nil {
		panic("Can not parse request url.")
	}

	ctx.Url = parsedUrl

	//start session
	ctx.Session = ctx.NewSession()
	ctx.readSession()

	if ctx.SessionId == "" {
		ctx.CreateSession()
	}

	ctx.Msg = NewMessageList()
	ctx.Css = css.New()
	ctx.Js = js.New()

	ctx.htmlOut = make(map[string]*string, 50)
	ctx.jsonOut = make(map[string]*string, 50)

	ctx.readMessages()
	ctx.clearMessages()

	if ctx.Session.EffectiveUserId != 0 { //if there is a userId in session.
		ctx.User, err = user.New(ctx.Session.EffectiveUserId)
		if err != nil {
			ctx.Msg.Error(err.Error())
			ctx.Session.EffectiveUserId = 0
		}
	}

	return ctx
}

func (ctx *Ctx) ReadCargo() {
	values := ctx.Url.Query()

	for k, _ := range ctx.Cargo {
		str := values.Get(k)
		if str != "" {
			ctx.Cargo.SetConvert(k, str)
		}
	}
}

func (ctx *Ctx) AddHtml(name string, cnt *string) {
	if _, ok := ctx.htmlOut[name]; ok {
		panic("Formatted content already exists: " + name)
	}

	ctx.htmlOut[name] = cnt
}

func (ctx *Ctx) AddJson(name string, cnt *string) {
	if _, ok := ctx.jsonOut[name]; ok {
		panic("Formatted content already exists: " + name)
	}

	ctx.jsonOut[name] = cnt
}

func (ctx *Ctx) Redirect(urlStr string) {
	ctx.saveMessages()
	ctx.SaveSession()

	http.Redirect(ctx.Rw, ctx.Req, urlStr, 302)
}

func (ctx *Ctx) Render(pageTmp string) {
	ctx.SaveSession()

	ctx.AddHtml("css", ctx.Css.Format())
	ctx.AddHtml("js", ctx.Js.Format())

	messageStr := ctx.Msg.Format(ctx)
	if len(*messageStr) > 0 {
		ctx.AddHtml("message", messageStr)
	}

	err := app.Tmpl.ExecuteTemplate(ctx.Rw, pageTmp, ctx.htmlOut)
	if err != nil {
		panic(err)
	}
}

func (ctx *Ctx) RenderAjax() {
	ajaxStatus := "success"
	ctx.AddJson("status", &ajaxStatus)

	messageStr := ctx.Msg.Format(ctx)
	if len(*messageStr) > 0 {
		ctx.AddJson("message", messageStr)
	}

	jsonStr, err := json.Marshal(ctx.jsonOut)
	if err != nil {
		panic(err)
	}

	ctx.Rw.Header().Set("Content-Type", "application/json")
	ctx.Rw.Header().Set("Content-Length", strconv.Itoa(len(jsonStr)))
	ctx.Rw.Write(jsonStr)
}

func (ctx *Ctx) RenderAjaxError() {
	ajaxStatus := "error"
	ctx.AddJson("status", &ajaxStatus)

	messageStr := ctx.Msg.Format(ctx)
	if len(*messageStr) > 0 {
		ctx.AddJson("message", messageStr)
	}

	jsonStr, err := json.Marshal(ctx.jsonOut)
	if err != nil {
		panic(err)
	}

	ctx.Rw.Header().Set("Content-Type", "application/json")
	ctx.Rw.Header().Set("Content-Length", strconv.Itoa(len(jsonStr)))
	ctx.Rw.Write(jsonStr)
}

func (ctx *Ctx) RenderEmail(pageTmp string) (string, error) {
	buf := new(bytes.Buffer)
	err := app.Tmpl.ExecuteTemplate(buf, pageTmp, ctx.htmlOut)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (ctx *Ctx) TotalPage(totalRows, pageLen int64) int64 {
	var r int64

	if totalRows%pageLen > 0 {
		r = 1
	} else {
		r = 0
	}

	totalPage := totalRows/pageLen + r

	return totalPage
}

func (ctx *Ctx) TouchPageNo(pageNo, totalRows, pageLen int64) int64 {
	totalPage := ctx.TotalPage(totalRows, pageLen)
	if pageNo < 1 {
		return 1
	} else if pageNo > totalPage {
		return totalPage
	}

	return pageNo
}

func (ctx *Ctx) IsRight(name, function string) bool {
	if ctx.User == nil {
		return false
	}

	return ctx.User.IsRight(name, function)
}

func (ctx *Ctx) IsLoggedIn() bool {
	return ctx.User != nil // if there is no login there is no user info in context.
}

func (ctx *Ctx) IsSuperuser() bool {
	if ctx.User == nil {
		return false
	}

	return ctx.User.IsSuperUser()

}

func (ctx *Ctx) U(urlStr string, args ...string) string {
	return ctx.Cargo.MakeUrl(urlStr, args...)
}

package login

import (
	"database/sql"
	"net/http"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
	"farlogin/src/page/user/user_lib"

	"github.com/dchest/captcha"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if req.URL.Path != "/login" && req.URL.Path != "/" {
		app.NotFound()
	}

	ctx.ReadCargo()

	if ctx.IsLoggedIn() {
		ctx.Msg.Warning("You are already logged in.")
		ctx.Redirect(ctx.U("/welcome"))
		return
	}

	if ctx.Req.Method == "GET" {
		loginForm(ctx)
		return
	}

	login := ctx.Req.PostFormValue("login")
	password := ctx.Req.PostFormValue("password")
	captchaId := ctx.Req.PostFormValue("captchaId")
	captchaAnswer := ctx.Req.PostFormValue("captchaAnswer")

	if login == "" || password == "" {
		ctx.Msg.Warning("You have left one or more fields empty.")
		loginForm(ctx)
		return
	}

	if !captcha.VerifyString(captchaId, captchaAnswer) {
		ctx.Msg.Warning("Wrong security answer. Please try again.")
		loginForm(ctx)
		return
	}

	password = util.PasswordHash(password)
	rec, err := user_lib.GetUserRecByLogin(login, password)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning("Wrong user name or password. Please try again.")
			loginForm(ctx)
			return
		}

		panic(err)
	}

	if rec.Status == "new" {
		ctx.Msg.Warning("Your have not confirmed your registration yet. Please check your e-mail box and follow the instructions.")
		loginForm(ctx)
		return
	}

	if rec.Status != "active" {
		ctx.Msg.Warning("Your account is not active. Please contact administrator to have your account activated.")
		loginForm(ctx)
		return
	}

	ctx.Session.UserId = rec.UserId
	ctx.Session.UserName = rec.Name

	ctx.Session.EffectiveUserId = rec.UserId
	ctx.Session.EffectiveUserName = rec.Name

	ctx.Redirect(ctx.U("/welcome"))
}

func loginForm(ctx *context.Ctx) {

	var login, password string
	if ctx.Req.Method == "POST" {
		login = ctx.Req.PostFormValue("login")
		password = ctx.Req.PostFormValue("password")
	} else {
		login = ""
		password = ""
	}

	buf := util.NewBuf()

	buf.Add("<div class=\"row login\">")
	buf.Add("<div class=\"col colSpan md2 lg3 xl4\"></div>")

	buf.Add("<div class=\"col md8 lg6 xl4 panel panelPrimary\">")

	buf.Add("<div class=\"panelHeading\">")
	buf.Add("<h3>farlogin</h3>")
	buf.Add("</div>")

	buf.Add("<div class=\"panelBody\">")

	buf.Add("<form action=\"/login\" method=\"post\">")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Login Name:</label>")
	buf.Add("<input type=\"text\" name=\"login\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"1\" autofocus>", util.ScrStr(login))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Password:</label>")
	buf.Add("<input type=\"password\" name=\"password\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"30\" tabindex=\"2\">", util.ScrStr(password))
	buf.Add("</div>")

	captchaId := captcha.NewLen(4)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Security Question:</label>")
	buf.Add("<input type=\"text\" name=\"captchaAnswer\"" +
		" class=\"formControl\" value=\"\" maxlength=\"10\" tabindex=\"3\">")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<p><img src=\"/captcha/%s.png\" alt=\"Captcha image.\"></p>", captchaId)
	buf.Add("<input type=hidden name=captchaId value=\"%s\">", captchaId)
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"4\">Submit</button>")
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"5\">Reset</button>")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<a href=\"/forgot_pass\">I forgot my password.</a>")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<a href=\"/register\">Create account.</a>")
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>") //panelBody
	buf.Add("</div>") //col
	buf.Add("<div class=\"col colSpan md2 lg3 xl4\"></div>")
	buf.Add("</div>") //row

	ctx.AddHtml("midContent", buf.String())

	content.Include(ctx)
	ctx.Css.Add("/asset/css/page/login.css")

	str := "loginPage"
	ctx.AddHtml("pageName", &str)

	content.Default(ctx)

	ctx.Render("login.html")
}

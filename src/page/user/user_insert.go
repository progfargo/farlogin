package user

import (
	"net/http"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/content/combo"
	"farlogin/src/content/left_menu"
	"farlogin/src/content/top_menu"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"

	"github.com/go-sql-driver/mysql"
)

func Insert(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("user", "insert") {
		app.BadRequest()
	}

	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddInt("rid", -1)
	ctx.Cargo.AddStr("stat", "default")
	ctx.ReadCargo()

	if ctx.Req.Method == "GET" {
		insertForm(ctx)
		return
	}

	name := ctx.Req.PostFormValue("name")
	login := ctx.Req.PostFormValue("login")
	email := ctx.Req.PostFormValue("email")
	password := ctx.Req.PostFormValue("password")

	status := ctx.Req.PostFormValue("status")

	if name == "" || login == "" || email == "" || password == "" || status == "default" {
		ctx.Msg.Warning("You have left one or more fields empty.")
		insertForm(ctx)
		return
	}

	if err := util.IsValidEmail(email); err != nil {
		ctx.Msg.Warning(err.Error())
		insertForm(ctx)
		return
	}

	if err := util.IsValidPassword(password); err != nil {
		ctx.Msg.Warning(err.Error())
		insertForm(ctx)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	password = util.PasswordHash(password)

	sqlStr := `insert into
					user(userId, name, login, email, password, resetKey, status)
					values(null, ?, ?, ?, ?, '', ?)`

	_, err = tx.Exec(sqlStr, name, login, email, password, status)
	if err != nil {
		tx.Rollback()
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1062 {
				ctx.Msg.Warning("Duplicate record.")
				insertForm(ctx)
				return
			} else if err.Number == 1452 {
				ctx.Msg.Warning("Could not find parent record.")
				insertForm(ctx)
				return
			}
		}

		panic(err)
	}

	tx.Commit()

	ctx.Msg.Success("Record has been saved.")
	ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
}

func insertForm(ctx *context.Ctx) {
	stat := ctx.Cargo.Str("stat")

	var name, login, email, password, status string

	if ctx.Req.Method == "POST" {
		name = ctx.Req.PostFormValue("name")
		login = ctx.Req.PostFormValue("login")
		email = ctx.Req.PostFormValue("email")
		password = ctx.Req.PostFormValue("password")
		status = ctx.Req.PostFormValue("status")
	} else {
		name = ""
		login = ""
		email = ""
		password = ""
		status = stat
	}

	statusCombo := combo.NewEnumCombo()
	statusCombo.Add("default", "Select Status")
	statusCombo.Add("active", "Active")
	statusCombo.Add("blocked", "Blocked")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Users", "New Record"))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx.U("/user", "key", "pn", "rid", "stat")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/user_insert", "key", "pn", "rid", "stat")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">User Name:</label>")
	buf.Add("<input type=\"text\" name=\"name\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"1\" autofocus>", util.ScrStr(name))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Login Name:</label>")
	buf.Add("<input type=\"text\" name=\"login\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"2\">", util.ScrStr(login))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Email:</label>")
	buf.Add("<input type=\"text\" name=\"email\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"3\">", util.ScrStr(email))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Password:</label>")
	buf.Add("<input type=\"text\" name=\"password\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"30\" tabindex=\"4\">", util.ScrStr(password))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Status:</label>")
	buf.Add("<select name=\"status\" class=\"formControl\" tabindex=\"5\">")

	buf.Add(statusCombo.Format(status))

	buf.Add("</select>")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"6\">Submit</button>")
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"7\">Reset</button>")
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Include(ctx)
	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "user")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

package config

import (
	"fmt"
	"net/http"
	"strconv"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/content/left_menu"
	"farlogin/src/content/top_menu"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"

	"github.com/go-sql-driver/mysql"
)

func Insert(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsLoggedIn() || !ctx.IsSuperuser() {
		app.BadRequest()
	}

	ctx.ReadCargo()

	if ctx.Req.Method == "GET" {
		insertForm(ctx)
		return
	}

	enum, err := strconv.ParseInt(ctx.Req.PostFormValue("enum"), 10, 64)
	if err != nil {
		enum = 0
	}

	name := ctx.Req.PostFormValue("name")
	varType := ctx.Req.PostFormValue("type")
	value := ctx.Req.PostFormValue("value")
	title := ctx.Req.PostFormValue("title")
	exp := ctx.Req.PostFormValue("exp")

	if name == "" || varType == "" || value == "" || title == "" || exp == "" {
		ctx.Msg.Warning("You have left one or more fields empty.")
		insertForm(ctx)
		return
	}

	if !util.IsValidIdentifier(name) {
		ctx.Msg.Warning(fmt.Sprintf("%s ^[a-zA-Z]{1}[a-zA-Z0-9]*", "Invalid variable name. Variable name pattern:"))
		insertForm(ctx)
		return
	}

	if varType != "int" && varType != "string" {
		ctx.Msg.Warning("Invalid variable type. Must be 'int' or 'string'.")
		insertForm(ctx)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `insert into
					config(configId, enum, name, type, value, title, exp)
					values(null, ?, ?, ?, ?, ?, ?)`

	_, err = tx.Exec(sqlStr, enum, name, varType, value, title, exp)
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

	app.ReadConfig()
	ctx.Msg.Success("Record has been saved.")
	ctx.Redirect(ctx.U("/config"))
}

func insertForm(ctx *context.Ctx) {
	content.Include(ctx)

	var enum int64
	var name, varType, value, title, exp string
	var err error
	if ctx.Req.Method == "POST" {
		enum, err = strconv.ParseInt(ctx.Req.PostFormValue("enum"), 10, 64)
		if err != nil {
			enum = 0
		}

		name = ctx.Req.PostFormValue("name")
		varType = ctx.Req.PostFormValue("type")
		value = ctx.Req.PostFormValue("value")
		title = ctx.Req.PostFormValue("title")
		exp = ctx.Req.PostFormValue("exp")
	} else {
		enum = 0
		name = ""
		varType = ""
		value = ""
		title = ""
		exp = ""
	}

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Configuration", "New Record"))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx.U("/config")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/config_insert")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Variable Name:</label>")
	buf.Add("<input type=\"text\" name=\"name\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"1\" autofocus>", util.ScrStr(name))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Variable Type:</label>")
	buf.Add("<input type=\"text\" name=\"type\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"10\" tabindex=\"2\">", util.ScrStr(varType))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Value:</label>")
	buf.Add("<input type=\"text\" name=\"value\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"250\" tabindex=\"3\">", util.ScrStr(value))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Title:</label>")
	buf.Add("<input type=\"text\" name=\"title\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"4\">", util.ScrStr(title))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Explanation:</label>")
	buf.Add("<input type=\"text\" name=\"exp\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"250\" tabindex=\"4\">", util.ScrStr(exp))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label>Enumeration:</label>")
	buf.Add("<input type=\"text\" name=\"enum\" class=\"formControl\""+
		" value=\"%d\" maxlength=\"11\" tabindex=\"5\">", enum)
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"6\">Submit</button>")
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"7\">Reset</button>")
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "config")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

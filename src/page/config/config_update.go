package config

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/content/left_menu"
	"farlogin/src/content/top_menu"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
	"farlogin/src/page/config/config_lib"

	"github.com/go-sql-driver/mysql"
)

func Update(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsLoggedIn() || !ctx.IsSuperuser() {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("configId", -1)
	ctx.Cargo.AddInt("gid", -1)
	ctx.ReadCargo()

	configId := ctx.Cargo.Int("configId")
	rec, err := config_lib.GetConfigRec(configId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning("Record could not be found.")
			ctx.Redirect(ctx.U("/config", "gid"))
			return
		}

		panic(err)
	}

	if ctx.Req.Method == "GET" {
		updateForm(ctx, rec)
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
		updateForm(ctx, rec)
		return
	}

	if !util.IsValidIdentifier(name) {
		ctx.Msg.Warning(fmt.Sprintf("%s ^[a-zA-Z]{1}[a-zA-Z0-9]*", "Invalid variable name. Variable name pattern:"))
		updateForm(ctx, rec)
		return
	}

	if varType != "int" && varType != "string" {
		ctx.Msg.Warning("Invalid variable type. Must be 'int' or 'string'.")
		updateForm(ctx, rec)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `update config set
					enum = ?,
					name = ?,
					type = ?,
					value = ?,
					title = ?,
					exp = ?
				where
					configId = ?`

	res, err := tx.Exec(sqlStr, enum, name, varType, value, title, exp, configId)
	if err != nil {
		tx.Rollback()
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1062 {
				ctx.Msg.Warning("Duplicate record.")
				updateForm(ctx, rec)
				return
			}
		} else if err.Number == 1452 {
			ctx.Msg.Warning("Could not find parent record.")
			updateForm(ctx, rec)
			return
		}

		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning("You did not change the record.")
		updateForm(ctx, rec)
		return
	}

	tx.Commit()

	app.ReadConfig()
	ctx.Msg.Success("Record has been changed.")
	ctx.Redirect(ctx.U("/config", "gid"))
}

func updateForm(ctx *context.Ctx, rec *config_lib.ConfigRec) {
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
		enum = rec.Enum
		name = rec.Name
		varType = rec.Type
		value = rec.Value
		title = rec.Title
		exp = rec.Exp
	}

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Configuration", "Update Record"))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx.U("/config")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/config_update", "configId", "gid")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Variable Name:</label>")
	buf.Add("<input type=\"text\" name=\"name\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"2\" autofocus>", util.ScrStr(name))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Variable Type:</label>")
	buf.Add("<input type=\"text\" name=\"type\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"10\" tabindex=\"3\">", util.ScrStr(varType))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Value:</label>")
	buf.Add("<input type=\"text\" name=\"value\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"250\" tabindex=\"4\">", util.ScrStr(value))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Title:</label>")
	buf.Add("<input type=\"text\" name=\"title\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"5\">", util.ScrStr(title))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Explanation:</label>")
	buf.Add("<input type=\"text\" name=\"exp\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"250\" tabindex=\"6\">", util.ScrStr(exp))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label>Enumeration:</label>")
	buf.Add("<input type=\"text\" name=\"enum\" class=\"formControl\""+
		" value=\"%d\" maxlength=\"11\" tabindex=\"7\">", enum)
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"8\">Submit</button>")
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"9\">Reset</button>")
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

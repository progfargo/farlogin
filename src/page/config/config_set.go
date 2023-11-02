package config

import (
	"database/sql"
	"net/http"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/content/left_menu"
	"farlogin/src/content/top_menu"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
	"farlogin/src/page/config/config_lib"
)

func Set(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("config", "set") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("configId", -1)
	ctx.ReadCargo()

	configId := ctx.Cargo.Int("configId")
	rec, err := config_lib.GetConfigRec(configId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning("Record could not be found.")
			ctx.Redirect(ctx.U("/config"))
			return
		}

		panic(err)
	}

	if ctx.Req.Method == "GET" {
		setForm(ctx, rec)
		return
	}

	value := ctx.Req.PostFormValue("value")

	if value == "" {
		ctx.Msg.Warning("You have left one or more fields empty.")
		setForm(ctx, rec)
		return
	}

	if val, ok := config_lib.CheckList[rec.Name]; ok {
		if err := val(value); err != nil {
			ctx.Msg.Warning(err.Error())
			setForm(ctx, rec)
			return
		}
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `update config set
					value = ?
				where
					configId = ?`

	res, err := tx.Exec(sqlStr, value, configId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning("You did not change the record.")
		setForm(ctx, rec)
		return
	}

	tx.Commit()

	app.ReadConfig()
	ctx.Msg.Success("Record has been changed.")
	ctx.Redirect(ctx.U("/config"))
}

func setForm(ctx *context.Ctx, rec *config_lib.ConfigRec) {
	content.Include(ctx)

	var value string
	if ctx.Req.Method == "POST" {
		value = ctx.Req.PostFormValue("value")
	} else {
		value = rec.Value
	}

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Configuration", "Set Value"))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx.U("/config")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/config_set", "configId")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">Value:</label>")
	buf.Add("<input type=\"text\" name=\"value\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"250\" tabindex=\"1\" autofocus>", util.ScrStr(value))
	buf.Add("<span class=\"helpBlock\">%s</span>", util.ScrStr(rec.Exp))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"2\">Submit</button>")
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"3\">Reset</button>")
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

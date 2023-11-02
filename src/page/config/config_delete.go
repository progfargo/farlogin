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

func Delete(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsLoggedIn() || !ctx.IsSuperuser() {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("configId", -1)
	ctx.Cargo.AddStr("confirm", "no")
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

	if ctx.Cargo.Str("confirm") != "yes" {
		deleteConfirm(ctx, rec)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `delete from
					config
				where
					configId = ?`

	res, err := tx.Exec(sqlStr, configId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning("Record could not be found.")
		ctx.Redirect(ctx.U("/config"))
		return
	}

	tx.Commit()

	app.ReadConfig()
	ctx.Msg.Success("Record has been deleted.")
	ctx.Redirect(ctx.U("/config"))
}

func deleteConfirm(ctx *context.Ctx, rec *config_lib.ConfigRec) {
	content.Include(ctx)

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Configuration", "Delete Record"))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx.U("/config")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<table>")
	buf.Add("<tbody>")
	buf.Add("<tr><th class=\"fixedMiddle\">Variable Name:</th><td>%s</td></tr>", util.ScrStr(rec.Name))
	buf.Add("<tr><th>Variable Type:</th><td>%s</td></tr>", util.ScrStr(rec.Type))
	buf.Add("<tr><th>Value:</th><td>%s</td></tr>", util.ScrStr(rec.Value))
	buf.Add("<tr><th>Title:</th><td>%s</td></tr>", util.ScrStr(rec.Title))
	buf.Add("<tr><th>Explanation:</th><td>%s</td></tr>", util.ScrStr(rec.Exp))
	buf.Add("<tr><th>Enumeration:</th><td>%d</td></tr>", rec.Enum)
	buf.Add("</tbody>")
	buf.Add("</table>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"callout calloutError\">")
	buf.Add("<h4>Please confirm:</h4>")
	buf.Add("<p>Do you realy want to delete this record?</p>")
	buf.Add("</div>")
	buf.Add("</div>")

	ctx.Cargo.SetStr("confirm", "yes")
	urlStr := ctx.U("/config_delete", "configId", "confirm")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"confirmCommand\">")
	buf.Add("<a href=\"%s\" class=\"button buttonError buttonSm\">Yes</a>", urlStr)
	buf.Add("</div>")
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

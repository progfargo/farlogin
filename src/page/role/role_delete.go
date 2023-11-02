package role

import (
	"database/sql"
	"net/http"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/content/left_menu"
	"farlogin/src/content/top_menu"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
	"farlogin/src/page/role/role_lib"

	"github.com/go-sql-driver/mysql"
)

func Delete(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("role", "delete") {
		app.BadRequest()
	}

	ctx.Cargo.AddStr("confirm", "no")
	ctx.Cargo.AddInt("roleId", -1)
	ctx.ReadCargo()

	roleId := ctx.Cargo.Int("roleId")
	rec, err := role_lib.GetRoleRec(roleId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning("Record could not be found.")
			ctx.Redirect(ctx.U("/role"))
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

	if rec.Name == "admin" {
		ctx.Msg.Warning("'admin' role can not be deleted.")
		ctx.Redirect(ctx.U("/role"))
		return
	}

	sqlStr := `delete from
					role
				where
					role.roleId = ?`

	res, err := tx.Exec(sqlStr, roleId)
	if err != nil {
		tx.Rollback()

		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1451 {
				ctx.Msg.Warning("Could not delete parent record.")
				ctx.Redirect(ctx.U("/role"))
				return
			}
		}

		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning("Record could not be found.")
		ctx.Redirect(ctx.U("/role"))
		return
	}

	tx.Commit()

	ctx.Msg.Success("Record has been deleted.")
	ctx.Redirect(ctx.U("/role"))
}

func deleteConfirm(ctx *context.Ctx, rec *role_lib.RoleRec) {
	content.Include(ctx)

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle("Roles", "Delete Record"))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx.U("/role")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<table>")
	buf.Add("<tbody>")
	buf.Add("<tr><th class=\"fixedMiddle\">Name:</th><td>%s</td></tr>", util.ScrStr(rec.Name))
	buf.Add("<tr><th>Explanation:</th><td>%s</td></tr>", util.ScrStr(rec.Exp))
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
	urlStr := ctx.U("/role_delete", "roleId", "confirm")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"confirmCommand\">")
	buf.Add("<a href=\"%s\" class=\"button buttonError buttonSm\">Yes</a>", urlStr)
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "role")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

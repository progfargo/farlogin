package user

import (
	"database/sql"
	"net/http"
	"strconv"

	"farlogin/src/app"
	"farlogin/src/lib/context"
	"farlogin/src/page/role/role_lib"
	"farlogin/src/page/user/user_lib"

	"github.com/go-sql-driver/mysql"
)

func UserRoleGrant(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("user", "role_grant") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("userId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddInt("rid", -1)
	ctx.Cargo.AddStr("stat", "default")
	ctx.ReadCargo()

	userId := ctx.Cargo.Int("userId")
	roleId, err := strconv.ParseInt(ctx.Req.PostFormValue("roleId"), 10, 64)
	if err != nil {
		roleId = -1
	}

	if roleId == -1 {
		ctx.Msg.Warning("You did not select a role.")
		ctx.Redirect(ctx.U("/user_role", "userId", "key", "pn", "rid", "stat"))
		return
	}

	if userId == -1 {
		ctx.Msg.Warning("Could not find user record.")
		ctx.Redirect(ctx.U("/user_role", "userId", "key", "pn", "rid", "stat"))
		return
	}

	userRec, err := user_lib.GetUserRec(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning("Could not find user record.")
			ctx.Redirect(ctx.U("/user_role", "userId", "key", "pn", "rid", "stat"))
			return
		}

		panic(err)
	}

	_, err = role_lib.GetRoleRec(roleId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning("Could not find role record.")
			ctx.Redirect(ctx.U("/user_role", "userId", "key", "pn", "rid", "stat"))
			return
		}

		panic(err)
	}

	if userRec.Login == "superuser" {
		ctx.Msg.Warning("'superuser' account can not be updated.")
		ctx.Redirect(ctx.U("/user_role", "userId", "key", "pn", "rid", "stat"))
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `insert into
					userRole(userId, roleId)
					values(?, ?)`

	_, err = tx.Exec(sqlStr, userId, roleId)
	if err != nil {
		tx.Rollback()
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1452 {
				ctx.Msg.Warning("Could not find parent record.")
				ctx.Redirect(ctx.U("/user_role", "userId", "key", "pn", "rid", "stat"))
				return
			} else if err.Number == 1062 {
				ctx.Msg.Warning("Duplicate record.")
				ctx.Redirect(ctx.U("/user_role", "userId", "key", "pn", "rid", "stat"))
				return
			}
		}

		panic(err)
	}

	tx.Commit()

	ctx.Msg.Success("Role has been granted.")
	ctx.Redirect(ctx.U("/user_role", "userId", "key", "pn", "rid", "stat"))
}

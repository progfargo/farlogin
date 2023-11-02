package context

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"farlogin/src/app"
	"farlogin/src/lib/util"
)

type sessionType struct {
	Msg               *MessageList
	UserId            int64
	UserName          string
	EffectiveUserId   int64
	EffectiveUserName string
}

func (ctx *Ctx) NewSession() *sessionType {
	rv := new(sessionType)

	return rv
}

func (ctx *Ctx) CreateSession() {
	now := time.Now()
	epoch := now.Unix()
	expire := now.Add(time.Duration(app.Ini.CookieExpires) * time.Second)
	randString := util.NewUUID()

	cookie := http.Cookie{
		Name:     app.Ini.CookieName,
		Value:    randString,
		Path:     app.Ini.CookiePath,
		Domain:   app.Ini.CookieDomain,
		Expires:  expire,
		Secure:   app.Ini.CookieSecure,
		HttpOnly: app.Ini.CookieHttpOnly,
	}

	http.SetCookie(ctx.Rw, &cookie)

	jsonData, err := json.Marshal(ctx.Session)
	if err != nil {
		panic(err)
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `insert into
					userSession(sessionId, recordTime, data)
					values(?, ?, ?)`

	_, err = tx.Exec(sqlStr, randString, epoch, jsonData)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	tx.Commit()
	ctx.SessionId = randString
}

func (ctx *Ctx) readSession() {
	sid := ""
	cookie, ok := ctx.Req.Cookie(app.Ini.CookieName)
	if ok == nil {
		sid = cookie.Value
	}

	if sid == "" {
		return
	}

	sqlStr := `select
					data
				from
					userSession
				where
					sessionId = ?`

	row := app.Db.QueryRow(sqlStr, sid)

	var data []byte

	err := row.Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}

		panic(err)
	}

	err = json.Unmarshal(data, &ctx.Session)
	if err != nil {
		panic(err)
	}

	ctx.SessionId = sid
}

func (ctx *Ctx) SaveSession() {
	jsonData, err := json.Marshal(ctx.Session)
	if err != nil {
		panic(err)
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `update userSession
					set data = ?
				where
					sessionId = ?`

	_, err = tx.Exec(sqlStr, jsonData, ctx.SessionId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	tx.Commit()
}

func (ctx *Ctx) DeleteSession() {
	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `delete from
					userSession
				where
					sessionId = ?`

	_, err = tx.Exec(sqlStr, ctx.SessionId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	tx.Commit()
}

func (ctx *Ctx) readMessages() {
	if ctx.Session.Msg != nil {
		ctx.Msg = ctx.Session.Msg
	}
}

func (ctx *Ctx) saveMessages() {
	if ctx.Msg != nil {
		ctx.Session.Msg = ctx.Msg
	}
}

func (ctx *Ctx) clearMessages() {
	ctx.Session.Msg = nil
}

package user_lib

import (
	"fmt"

	"farlogin/src/app"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
)

type UserRec struct {
	UserId    int64
	Name      string
	Login     string
	Email     string
	Password  string
	ResetKey  string
	ResetTime int64
	Status    string
}

func GetUserRec(userId int64) (*UserRec, error) {
	sqlStr := `select
					user.userId,
					user.name,
					user.login,
					user.email,
					user.password,
					user.status
				from
					user
				where
					user.userId = ?`

	row := app.Db.QueryRow(sqlStr, userId)

	rec := new(UserRec)
	err := row.Scan(&rec.UserId, &rec.Name, &rec.Login, &rec.Email,
		&rec.Password, &rec.Status)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func CountUser(key string, roleId int64, status string) int64 {
	sqlBuf := util.NewBuf()
	sqlBuf.Add("select count(*)")

	fromBuf := util.NewBuf()
	fromBuf.Add("user")

	conBuf := util.NewBuf()

	if roleId != -1 {
		fromBuf.Add("userRole")
		conBuf.Add("(user.userId = userRole.userId)")
		conBuf.Add("(userRole.roleId = %d)", roleId)
	}

	if status != "default" {
		conBuf.Add("(user.status = '%s')", util.DbStr(status))
	}

	if key != "" {
		key = util.DbStr(key)
		conBuf.Add(`(user.name like('%%%s%%') or
					user.email like('%%%s%%') or
					user.login like('%%%s%%'))`, key, key, key)
	}

	sqlBuf.Add("from " + *fromBuf.StringSep(", "))

	if !conBuf.IsEmpty() {
		sqlBuf.Add("where " + *conBuf.StringSep(" and "))
	}

	row := app.Db.QueryRow(*sqlBuf.String())

	var rv int64
	err := row.Scan(&rv)
	if err != nil {
		panic(err)
	}

	return rv
}

func CountAdmin() int64 {
	sqlStr := `select
					count(*)
				from
					user,
					userRole,
					role
				where
					user.userId = userRole.userId and
					userRole.roleId = role.roleId and
					role.name = 'admin'`

	row := app.Db.QueryRow(sqlStr)

	var rv int64
	err := row.Scan(&rv)
	if err != nil {
		panic(err)
	}

	return rv
}

func IsUserRoleExists(userId int64, roleName string) bool {
	userRoleList := GetUserRoleList(userId)
	for _, v := range userRoleList {
		if v == roleName {
			return true
		}
	}

	return false
}

func GetUserPage(ctx *context.Ctx, key string, pageNo, roleId int64, status string) []*UserRec {

	sqlBuf := util.NewBuf()
	sqlBuf.Add(`select
					user.userId,
					user.name,
					user.login,
					user.email,
					user.status`)

	fromBuf := util.NewBuf()
	fromBuf.Add("user")

	conBuf := util.NewBuf()

	if roleId != -1 {
		fromBuf.Add("userRole")
		conBuf.Add("(user.userId = userRole.userId)")
		conBuf.Add("(userRole.roleId = %d)", roleId)
	}

	if status != "default" {
		conBuf.Add("(user.status = '%s')", util.DbStr(status))
	}

	if key != "" {
		key = util.DbStr(key)
		conBuf.Add(`(user.name like('%%%s%%') or
					user.email like('%%%s%%') or
					user.login like('%%%s%%'))`, key, key, key)
	}

	sqlBuf.Add("from " + *fromBuf.StringSep(", "))

	if !conBuf.IsEmpty() {
		sqlBuf.Add(" where " + *conBuf.StringSep(" and "))
	}

	sqlBuf.Add("order by user.name")

	pageLen := ctx.Config.Int("pageLen")
	start := (pageNo - 1) * pageLen
	sqlBuf.Add("limit %d, %d", start, pageLen)

	rows, err := app.Db.Query(*sqlBuf.String())
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*UserRec, 0, 100)
	for rows.Next() {
		rec := new(UserRec)
		err = rows.Scan(&rec.UserId, &rec.Name, &rec.Login,
			&rec.Email, &rec.Status)
		if err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}

func GetUserRecByLogin(name, password string) (*UserRec, error) {
	sqlStr := `select
					user.userId,
					user.name,
					user.login,
					user.name,
					user.email,
					user.status
				from
					user
				where
					user.login = ? and
					user.password = ?`

	row := app.Db.QueryRow(sqlStr, name, password)

	rec := new(UserRec)
	err := row.Scan(&rec.UserId, &rec.Name, &rec.Login,
		&rec.Name, &rec.Email, &rec.Status)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetUserRecByEmail(email string) (*UserRec, error) {
	sqlStr := `select
					user.userId,
					user.name,
					user.login,
					user.email,
					user.status,
					user.resetKey,
					user.resetTime
				from
					user
				where
					user.email = ?`

	row := app.Db.QueryRow(sqlStr, email)

	rec := new(UserRec)
	err := row.Scan(&rec.UserId, &rec.Name, &rec.Login, &rec.Email,
		&rec.Status, &rec.ResetKey, &rec.ResetTime)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetUserRecByResetKey(key string) (*UserRec, error) {
	sqlStr := `select
					user.userId,
					user.name,
					user.login,
					user.email,
					user.status,
					user.resetKey,
					user.resetTime
				from
					user
				where
					user.resetKey = ?`

	row := app.Db.QueryRow(sqlStr, key)

	rec := new(UserRec)
	err := row.Scan(&rec.UserId, &rec.Name, &rec.Login,
		&rec.Email, &rec.Status, &rec.ResetKey, &rec.ResetTime)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetUserRecByEmailVerifyKey(key string) (*UserRec, error) {
	sqlStr := `select
					user.userId,
					user.name,
					user.login,
					user.email,
					user.status,
					user.resetKey,
					user.resetTime
				from
					user
				where
					user.emailVerifyKey = ?`

	row := app.Db.QueryRow(sqlStr, key)

	rec := new(UserRec)
	err := row.Scan(&rec.UserId, &rec.Name, &rec.Login,
		&rec.Email, &rec.Status, &rec.ResetKey, &rec.ResetTime)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetUserUserRecBySessionId(sessionId string) (*UserRec, error) {
	sqlStr := `select
					user.userId,
					user.name,
					user.login,
					user.email,
					user.status,
					user.resetKey,
					user.resetTime,
				from
					user,
					session
				where
					session.sessionId = ? and
					user.userId = session.userId`

	row := app.Db.QueryRow(sqlStr, sessionId)

	rec := new(UserRec)
	err := row.Scan(&rec.UserId, &rec.Name, &rec.Login, &rec.Email,
		&rec.Status, &rec.Status, &rec.ResetKey, &rec.ResetTime)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetUserRoleList(userId int64) []string {
	sqlStr := `select
					role.name
				from
					role,
					userRole
				where
					role.roleId = userRole.roleId and
					userRole.userId = ? 
				order by
					name`

	rows, err := app.Db.Query(sqlStr, userId)
	if err != nil {
		panic("Could not read user roles." + " " + err.Error())
	}

	defer rows.Close()

	var rv []string
	var roleName string
	for rows.Next() {
		if err = rows.Scan(&roleName); err != nil {
			panic(err.Error())
		}

		rv = append(rv, roleName)
	}

	return rv
}

type UserBillingInfoRec struct {
	BillingInfoId int64
	UserId        int64
	InWhoseName   string
	Address1      string
	Address2      string
	PostCode      string
	City          string
	Country       string
	TaxId         string
	TaxOffice     string
}

func GetUserBillingInfoRec(userId int64) (*UserBillingInfoRec, error) {
	sqlStr := `select
					billingInfo.billingInfoId,
					billingInfo.userId,
					billingInfo.inWhoseName,
					billingInfo.address1,
					billingInfo.address2,
					billingInfo.postCode,
					billingInfo.city,
					billingInfo.country,
					billingInfo.taxId,
					billingInfo.taxOffice
				from
					billingInfo
				where
					billingInfo.userId = ?`

	row := app.Db.QueryRow(sqlStr, userId)

	rec := new(UserBillingInfoRec)
	err := row.Scan(&rec.BillingInfoId, &rec.UserId, &rec.InWhoseName,
		&rec.Address1, &rec.Address2, &rec.PostCode, &rec.City, &rec.Country,
		&rec.TaxId, &rec.TaxOffice)
	if err != nil {
		return rec, err
	}

	return rec, nil
}

type UserCardInfoRec struct {
	CreditCardId   int64
	UserId         int64
	NameOnCard     string
	CardNo         string
	ExpirationDate string
	Cvv            string
}

func GetUserCardInfoRec(userId int64) (*UserCardInfoRec, error) {
	sqlStr := `select
					creditCard.creditCardId,
					creditCard.userId,
					creditCard.nameOnCard,
					creditCard.cardNo,
					creditCard.expirationDate,
					creditCard.cvv
				from
					creditCard
				where
					creditCard.userId = ?`

	row := app.Db.QueryRow(sqlStr, userId)

	rec := new(UserCardInfoRec)
	err := row.Scan(&rec.CreditCardId, &rec.UserId, &rec.NameOnCard,
		&rec.CardNo, &rec.ExpirationDate, &rec.Cvv)
	if err != nil {
		return rec, err
	}

	return rec, nil
}

func StatusToLabel(str string) string {
	switch str {
	case "active":
		return fmt.Sprintf("<span class=\"label labelSuccess labelXs\">active</span>")
	case "blocked":
		return fmt.Sprintf("<span class=\"label labelError labelXs\">blocked</span>")
	default:
		return fmt.Sprintf("<span class=\"label labelDefault labelXs\">unknown</span>")
	}
}

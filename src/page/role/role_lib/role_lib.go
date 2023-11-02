package role_lib

import (
	"farlogin/src/app"
)

type RoleRec struct {
	RoleId int64
	Name   string
	Exp    string
}

type RoleRightRec struct {
	RoleId   int64
	PageName string
	FuncName string
}

func GetRoleRec(roleId int64) (*RoleRec, error) {
	sqlStr := `select
					role.roleId,
					role.name,
					role.exp
				from
					role
				where
					role.roleId = ?`

	row := app.Db.QueryRow(sqlStr, roleId)

	rec := new(RoleRec)
	err := row.Scan(&rec.RoleId, &rec.Name, &rec.Exp)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetRoleList() []*RoleRec {
	sqlStr := `select
					role.roleId,
					role.name,
					role.exp
				from
					role
				order by role.name`

	rows, err := app.Db.Query(sqlStr)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*RoleRec, 0, 20)
	for rows.Next() {
		rec := new(RoleRec)
		if err = rows.Scan(&rec.RoleId, &rec.Name, &rec.Exp); err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}

func CountRoleByUser(userId int64) (int64, error) {
	sqlStr := `select
					count(*)
				from
					role,
					userRole
				where
					role.roleId = userRole.roleId and
					userRole.userId = ?`

	rows := app.Db.QueryRow(sqlStr, userId)

	var rv int64
	if err := rows.Scan(&rv); err != nil {
		return 0, err
	}

	return rv, nil
}

func GetRoleByUser(userId int64) ([]*RoleRec, error) {
	sqlStr := `select
					role.roleId,
					role.name,
					role.exp
				from
					role,
					userRole
				where
					role.roleId = userRole.roleId and
					userRole.userId = ?
				order by
					role.name`

	rows, err := app.Db.Query(sqlStr, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rv := make([]*RoleRec, 0, 10)
	for rows.Next() {
		rec := new(RoleRec)
		if err = rows.Scan(&rec.RoleId, &rec.Name, &rec.Exp); err != nil {
			return nil, err
		}

		rv = append(rv, rec)
	}

	return rv, nil
}

func GetRoleRight(roleId int64) (map[string]bool, error) {
	sqlStr := `select
					roleId,
					pageName,
					funcName
				from
					roleRight
				where
					roleId = ?
				order by
					pageName, funcName`

	rows, err := app.Db.Query(sqlStr, roleId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rv := make(map[string]bool, 50)
	for rows.Next() {
		rec := new(RoleRightRec)
		err = rows.Scan(&rec.RoleId, &rec.PageName, &rec.FuncName)
		if err != nil {
			return nil, err
		}

		rv[app.MakeKey(rec.PageName, rec.FuncName)] = true
	}

	return rv, nil
}

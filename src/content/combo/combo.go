package combo

import (
	"database/sql"

	"farlogin/src/app"
	"farlogin/src/lib/util"
)

type comboItem struct {
	key   int64
	label sql.NullString
}

type combo struct {
	sqlStr        string
	defaultOption sql.NullString
	list          []*comboItem
}

func NewCombo(sqlStr, defaultOption string) *combo {
	rv := new(combo)
	rv.list = make([]*comboItem, 0, 50)
	rv.sqlStr = sqlStr
	rv.defaultOption = sql.NullString{String: defaultOption, Valid: true}

	return rv
}

func (com *combo) IsEmpty() bool {
	return len(com.list) < 2
}

func (com *combo) Set() {
	com.list = append(com.list, &comboItem{-1, com.defaultOption})

	rows, err := app.Db.Query(com.sqlStr)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		item := new(comboItem)
		if err = rows.Scan(&item.key, &item.label); err != nil {
			panic(err)
		}

		com.list = append(com.list, item)
	}
}

func (com *combo) Format(curKey int64) string {
	buf := util.NewBuf()

	var selectedStr, label string
	for _, item := range com.list {
		selectedStr = ""

		if curKey == item.key {
			selectedStr = " selected"
		}

		label = util.NullToString(item.label)
		buf.Add("<option value=\"%d\"%s>%s</option>", item.key, selectedStr, label)
	}

	return *buf.String()
}

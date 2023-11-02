package config_lib

import (
	"errors"
	"strconv"

	"farlogin/src/app"
	"farlogin/src/lib/util"
)

type ConfigRec struct {
	ConfigId int64
	Enum     int64
	Name     string
	Type     string
	Value    string
	Title    string
	Exp      string
}

type checkFunc func(str string) error

var CheckList map[string]checkFunc

func init() {
	CheckList := make(map[string]checkFunc, 10)
	CheckList["maintenanceMode"] = checkMaintenanceMode
	CheckList["pageLen"] = checkPageLen
	CheckList["rulerLen"] = checkRulerLen
}

func GetConfigRec(configId int64) (*ConfigRec, error) {
	sqlStr := `select
					config.configId,
					config.enum,
					config.name,
					config.type,
					config.value,
					config.title,
					config.exp
				from
					config
				where
					config.configId = ?`

	row := app.Db.QueryRow(sqlStr, configId)

	rec := new(ConfigRec)
	if err := row.Scan(&rec.ConfigId, &rec.Enum, &rec.Name, &rec.Type,
		&rec.Value, &rec.Title, &rec.Exp); err != nil {
		return nil, err
	}

	return rec, nil
}

func GetConfigList() []*ConfigRec {
	sqlBuf := util.NewBuf()
	sqlBuf.Add(`select
					config.configId,
					config.enum,
					config.name,
					config.type,
					config.value,
					config.title,
					config.exp
				from
					config`)

	sqlBuf.Add("order by config.enum")

	rows, err := app.Db.Query(*sqlBuf.String())
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*ConfigRec, 0, 100)
	for rows.Next() {
		rec := new(ConfigRec)
		if err = rows.Scan(&rec.ConfigId, &rec.Enum, &rec.Name, &rec.Type,
			&rec.Value, &rec.Title, &rec.Exp); err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}

//check functions
func checkPageLen(val string) error {
	pageLen, err := strconv.ParseInt(val, 10, 64)
	if err == nil && pageLen >= 5 && pageLen <= 100 {
		return nil
	}

	return errors.New("Invalid page length. Please enter a number between 10-100.")
}

func checkRulerLen(val string) error {
	rulerLen, err := strconv.ParseInt(val, 10, 64)
	if err == nil && rulerLen >= 2 && rulerLen <= 10 {
		return nil
	}

	return errors.New("Invalid ruler length. Please enter a number between 2-10.")
}

func checkMaxImageSize(val string) error {
	maxImageSize, err := strconv.ParseInt(val, 10, 64)
	if err == nil && maxImageSize > 5242880 && maxImageSize < 52428800 {
		return nil
	}

	return errors.New("Invalid maximum image size. Size must be between  5242880 (5M) and 52428800 (50M)")
}

func checkMaintenanceMode(val string) error {
	if val == "yes" || val == "no" {
		return nil
	}

	return errors.New("Invalid value. Value must be 'yes' or 'no'.")
}

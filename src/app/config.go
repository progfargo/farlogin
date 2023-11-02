package app

import (
	"fmt"
	"strconv"
)

type configRec struct {
	name       string
	configType string
	value      string
}

type configVar struct {
	varType string
	value   string
}

type ConfigType map[string]configVar

var config = make(ConfigType, 10)

func ReadConfig() {
	sqlStr := `select
					config.name,
					config.type,
					config.value
				from
					config
				order by config.name`

	rows, err := Db.Query(sqlStr)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rec := new(configRec)
	for rows.Next() {
		err = rows.Scan(&rec.name, &rec.configType, &rec.value)
		if err != nil {
			panic(err)
		}

		config.Set(rec.name, rec.configType, rec.value)
	}
}

func CopyConfig() ConfigType {
	rv := make(ConfigType, 10)

	for k, v := range config {
		rv[k] = v
	}

	return rv
}

func (cnf ConfigType) Set(name, configType, value string) {
	cnf[name] = configVar{configType, value}
}

func (cnf ConfigType) Str(name string) string {
	v, ok := cnf[name]
	if !ok {
		panic(fmt.Errorf("Could not find configuration variable. name: %s", name))
	}

	if v.varType != "string" {
		panic(fmt.Errorf("Configuration variable type mismatch. name: %s", name))
	}

	return v.value
}

func (cnf ConfigType) Int(name string) int64 {
	v, ok := cnf[name]
	if !ok {
		panic(fmt.Errorf("could not find configuration variable. name: %s", name))
	}

	if v.varType != "int" {
		panic(fmt.Errorf("Configuration variable type mismatch. name: %s", name))
	}

	num, err := strconv.ParseInt(v.value, 10, 64)
	if err != nil {
		panic(err)
	}

	return num
}

func (cnf ConfigType) Float(name string) float64 {
	v, ok := cnf[name]
	if !ok {
		panic(fmt.Errorf("could not find configuration variable. name: %s", name))
	}

	if v.varType != "float" {
		panic(fmt.Errorf("Configuration variable type mismatch. name: %s", name))
	}

	num, err := strconv.ParseFloat(v.value, 64)
	if err != nil {
		panic(err)
	}

	return num
}

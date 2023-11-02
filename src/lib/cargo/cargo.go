package cargo

import (
	"net/url"
	"reflect"
	"strconv"
)

type cargoVar struct {
	Name         string
	Type         reflect.Kind
	Value        interface{}
	DefaultValue interface{}
}

type CargoList map[string]*cargoVar

func NewCargo() CargoList {
	rv := make(CargoList, 10)

	return rv
}

func (cargo CargoList) AddInt(name string, defaultValue int64) {
	_, ok := cargo[name]
	if ok {
		panic("Cargo variable already exists: " + name)
	}

	cargo[name] = &cargoVar{Name: name, DefaultValue: defaultValue, Type: reflect.Int64}
}

func (cargo CargoList) AddStr(name string, defaultValue string) {
	_, ok := cargo[name]
	if ok {
		panic("Cargo variable already exists: " + name)
	}

	cargo[name] = &cargoVar{Name: name, DefaultValue: defaultValue, Type: reflect.String}
}

func (cargo CargoList) IsExists(name string) bool {
	_, ok := cargo[name]
	if ok {
		return true
	}

	return false
}

func (cargo CargoList) SetInt(name string, value int64) {
	_, ok := cargo[name]
	if !ok {
		panic("Cargo variable does not exist: " + name)
	}

	if cargo[name].Type != reflect.Int64 {
		panic("Cargo variable type mismatch: " + name)
	}

	cargo[name].Value = value
}

func (cargo CargoList) SetStr(name string, value string) {
	_, ok := cargo[name]
	if !ok {
		panic("Cargo variable does not exist: " + name)
	}

	if cargo[name].Type != reflect.String {
		panic("Cargo variable type mismatch: " + name)
	}

	cargo[name].Value = value
}

func (cargo CargoList) SetConvert(name string, str string) {
	_, ok := cargo[name]
	if !ok {
		panic("Cargo variable does not exist: " + name)
	}

	switch cargo[name].Type {
	case reflect.Int64:
		v, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			cargo[name].Value = v
		}

	case reflect.String:
		cargo[name].Value = str
	}
}

func (cargo CargoList) Int(name string) int64 {
	val, ok := cargo[name]
	if !ok {
		panic("Cargo variable does not exists: " + name)
	}

	switch val.Type {
	case reflect.Int64:
		if val.Value == nil {
			return val.DefaultValue.(int64)
		}

		return val.Value.(int64)
	case reflect.String:
		var str string
		if val.Value != nil {
			rv, err := strconv.ParseInt(val.Value.(string), 10, 64)

			if err == nil {
				return rv
			}
		}

		str = val.DefaultValue.(string)
		rv, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			panic("Could not convert string to int64: " + name)
		}

		return rv
	default:
		panic("Unknown cargo variable type: " + name)
	}
}

func (cargo CargoList) Str(name string) string {
	val, ok := cargo[name]
	if !ok {
		panic("Cargo variable does not exists: " + name)
	}

	switch val.Type {
	case reflect.Int64:
		var num int64
		if val.Value == nil {
			num = val.DefaultValue.(int64)
		} else {
			num = val.Value.(int64)
		}

		return strconv.FormatInt(num, 10)
	case reflect.String:
		if val.Value == nil {
			return val.DefaultValue.(string)
		}

		return val.Value.(string)
	default:
		panic("Unknown cargo variable type: " + name)
	}
}

func (cargo CargoList) IsDefault(name string) bool {
	if _, ok := cargo[name]; !ok {
		panic("Unknown cargo variable: " + name)
	}

	switch cargo[name].Type {
	case reflect.Int64:
		if cargo[name].Value == nil || cargo[name].Value.(int64) == cargo[name].DefaultValue.(int64) {
			return true
		}
	case reflect.String:
		if cargo[name].Value == nil || cargo[name].Value.(string) == cargo[name].DefaultValue.(string) {
			return true
		}
	default:
		panic("Unnown cargo variable type.")
	}

	return false
}

func (cargo CargoList) MakeUrl(urlStr string, args ...string) string {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		panic(err)
	}

	q := u.Query()

	for _, v := range args {
		if _, ok := cargo[v]; !ok {
			panic("Unknown cargo variable: " + v)
		}

		if cargo.IsDefault(v) {
			continue
		}

		q.Set(v, cargo.Str(v))
	}

	u.RawQuery = q.Encode()

	rv := u.RequestURI()

	return rv
}

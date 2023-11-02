package app

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ini struct {
	HomeDir string
	BaseUrl string
	Host    string
	Port    string

	DbProtocol string
	DbHost     string
	DbPort     int64
	DbName     string
	DbUser     string
	DbPassword string

	CookieName     string
	CookiePath     string
	CookieDomain   string
	CookieExpires  int64
	CookieSecure   bool
	CookieHttpOnly bool

	Cache bool
	Debug bool
}

var Ini = new(ini)

func readIni() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	file := dir + "/farlogin.ini"

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
		os.Exit(1)
	}

	err = json.Unmarshal(buf, Ini)
	if err != nil {
		panic(err)
		os.Exit(1)
	}
}

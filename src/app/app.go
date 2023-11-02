package app

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/natefinch/lumberjack"

	"github.com/nats-io/nats-server/v2/server"
)

type BadRequestError error
type NotFoundError error
type MaintenanceError error

func BadRequest() {
	panic(new(BadRequestError))
}

func NotFound() {
	panic(new(NotFoundError))
}

var Debug = true
var SiteName = "farlogin"

var Tmpl *template.Template
var Db *sql.DB

var NatsUrl string
var NatsPort int

func init() {
	readIni()
	connect()
	ReadConfig()

	log.SetOutput(&lumberjack.Logger{
		Filename:   Ini.HomeDir + "/logs/error.log",
		MaxSize:    10, // megabytes
		MaxBackups: 6,
		MaxAge:     28, //days
	})

	funcmap := template.FuncMap{
		"raw": func(str string) template.HTML {
			return template.HTML(str)
		},
	}

	Tmpl = template.Must(template.New("").Funcs(funcmap).ParseGlob(Ini.HomeDir + "/view/*.html"))

	NatsPort = 4222
	NatsUrl = "localhost"

	//nats server
	opts := &server.Options{
		ServerName: "farlogin",
		Host:       NatsUrl,
		Port:       NatsPort,
	}

	ns, err := server.NewServer(opts)

	if err != nil {
		panic(err)
	}

	go ns.Start()

	if !ns.ReadyForConnections(10 * time.Second) {
		panic("nats.io is not ready for connection")
	}
}

func connect() {
	conStr := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?autocommit=false",
		Ini.DbUser, Ini.DbPassword, Ini.DbProtocol,
		Ini.DbHost, Ini.DbPort, Ini.DbName)
	db, err := sql.Open("mysql", conStr)

	if err != nil {
		panic("Could not connec to database.")
	}

	db.SetMaxIdleConns(0)
	Db = db
}

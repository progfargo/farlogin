package main

import (
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
	"farlogin/src/lib/watch"
	"farlogin/src/page/config"
	"farlogin/src/page/login"
	"farlogin/src/page/logout"
	"farlogin/src/page/node"
	"farlogin/src/page/profile"

	"farlogin/src/page/role"
	"farlogin/src/page/user"
	"farlogin/src/page/welcome"

	"github.com/dchest/captcha"
)

type HandleFuncType func(http.ResponseWriter, *http.Request)

func main() {
	cron()

	if app.Ini.Debug {
		watch.WatchLess()
	}

	//static constent
	fsAdmin := http.FileServer(http.Dir(app.Ini.HomeDir + "/asset/"))
	http.Handle("/asset/", http.StripPrefix("/asset/", fsAdmin))

	http.Handle("/captcha/", captcha.Server(200, 90))

	//admin
	http.HandleFunc("/role", makeHandle(role.Browse))
	http.HandleFunc("/role_insert", makeHandle(role.Insert))
	http.HandleFunc("/role_update", makeHandle(role.Update))
	http.HandleFunc("/role_delete", makeHandle(role.Delete))
	http.HandleFunc("/role_right", makeHandle(role.RoleRight))

	http.HandleFunc("/user", makeHandle(user.Browse))
	http.HandleFunc("/user_insert", makeHandle(user.Insert))
	http.HandleFunc("/user_update", makeHandle(user.Update))
	http.HandleFunc("/user_update_pass", makeHandle(user.UpdatePass))
	http.HandleFunc("/user_delete", makeHandle(user.Delete))
	http.HandleFunc("/user_role", makeHandle(user.UserRole))
	http.HandleFunc("/user_role_revoke", makeHandle(user.UserRoleRevoke))
	http.HandleFunc("/user_role_grant", makeHandle(user.UserRoleGrant))
	http.HandleFunc("/user_display", makeHandle(user.Display))

	http.HandleFunc("/config", makeHandle(config.Browse))
	http.HandleFunc("/config_insert", makeHandle(config.Insert))
	http.HandleFunc("/config_update", makeHandle(config.Update))
	http.HandleFunc("/config_delete", makeHandle(config.Delete))
	http.HandleFunc("/config_set", makeHandle(config.Set))

	http.HandleFunc("/login", makeHandle(login.Browse))
	http.HandleFunc("/welcome", makeHandle(welcome.Browse))
	http.HandleFunc("/logout", makeHandle(logout.Browse))

	http.HandleFunc("/profile", makeHandle(profile.Browse))
	http.HandleFunc("/profile_update", makeHandle(profile.UpdateProfile))
	http.HandleFunc("/profile_update_pass", makeHandle(profile.UpdatePass))

	http.HandleFunc("/node", makeHandle(node.Browse))
	http.HandleFunc("/node_display", makeHandle(node.Display))
	http.HandleFunc("/node_insert", makeHandle(node.Insert))
	http.HandleFunc("/node_update", makeHandle(node.Update))
	http.HandleFunc("/node_delete", makeHandle(node.Delete))
	http.HandleFunc("/node_new_session", makeHandle(node.NewSession))
	http.HandleFunc("/node_delete_session", makeHandle(node.DeleteSession))
	http.HandleFunc("/node_delete_all_session", makeHandle(node.DeleteAllSession))

	http.HandleFunc("/", makeHandle(login.Browse))

	//timeout for old browsers and old mobile devices
	srv := &http.Server{
		Addr:         app.Ini.Host + ":" + app.Ini.Port,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Printf("%s starting server on port: %s\n", app.SiteName, app.Ini.Port)
	log.Printf("%s starting server on port: %s\n", app.SiteName, app.Ini.Port)
	err := srv.ListenAndServe()

	if err != nil {
		println(err.Error())
		log.Println(err.Error())
	}
}

func makeHandle(hf HandleFuncType) HandleFuncType {
	return func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				switch err := err.(type) {
				case string:
					fatalError(rw, req, err)
				case error:
					fatalError(rw, req, err.Error())
				case *app.BadRequestError:
					badRequest(rw, req)
				case *app.NotFoundError:
					notFound(rw, req)
				default:
					fmt.Printf("%s %v", "unknown error", err)
					log.Printf("%s %v\n", "unknown error", err)
				}
			}
		}()

		if app.Ini.Cache {
			rw.Header().Add("Cache-Control", "public, max-age=3600, must-revalidate")
		} else {
			rw.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
		}

		hf(rw, req)
	}
}

func BasicAuth(hf HandleFuncType, username, password, realm string) HandleFuncType {

	return func(rw http.ResponseWriter, req *http.Request) {

		user, pass, ok := req.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			rw.Header().Set("WWW-Authenticate", "Basic realm=\""+realm+"\"")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorised."))
			return
		}

		hf(rw, req)
	}
}

func fatalError(rw http.ResponseWriter, req *http.Request, msg string) {
	if req.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		list := strings.Split(string(debug.Stack()[:]), "\n")
		errorMsg := fmt.Sprintf("Fatal Error: %s", msg)

		if app.Ini.Debug {
			errorMsg += "\n" + strings.Join(list[6:], "\n")
		}

		http.Error(rw, errorMsg, http.StatusInternalServerError)
		return
	}

	buf := new(util.Buf)

	buf.Add("<!DOCTYPE html>")
	buf.Add("<html lang=\"en\">")

	buf.Add("<head>")
	buf.Add("<meta charset=\"utf-8\">")
	buf.Add("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">")
	buf.Add("<link rel=\"icon\" href=\"/asset/img/favicon.png\">")
	buf.Add("<link rel=\"stylesheet\" href=\"/asset/css/style.css\">")
	buf.Add("<link rel=\"stylesheet\" href=\"/asset/css/page/fatal_error.css\">")
	buf.Add("<title>Fatal Error</title>")
	buf.Add("</head>")

	buf.Add("<body id=\"fatalErrorPage\">")

	buf.Add("<main>")
	buf.Add("<div class=\"row\">")

	buf.Add("<div class=\"col colSpan md1 lg2 xl3\"></div>")

	buf.Add("<div class=\"col md10 lg8 xl6 panel panelError\">")

	buf.Add("<div class=\"panelHeading\">")
	buf.Add("<h3>Fatal Error</h3>")
	buf.Add("</div>")

	buf.Add("<div class=\"panelBody\">")
	buf.Add("<p>%s</p>", msg)

	list := strings.Split(string(debug.Stack()[:]), "\n")
	buf.Add("<p>%s</p>", strings.Join(list[6:], "<br>"))
	buf.Add("</div>")

	buf.Add("</div>") //end of col

	buf.Add("<div class=\"col colSpan md1 lg2 xl3\"></div>")

	buf.Add("</div>") //end of row

	buf.Add("</main>")

	buf.Add("</body>")
	buf.Add("</html>")

	fmt.Fprintf(rw, *buf.String())
}

func badRequest(rw http.ResponseWriter, req *http.Request) {
	if req.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := context.NewContext(rw, req)

	ctx.Rw.WriteHeader(http.StatusBadRequest)
	buf := util.NewBuf()
	buf.Add("<div class=\"callout calloutError\">")
	buf.Add("<h3>Bad request.</h3>")
	buf.Add("<p>You don't have right to access this page.</p>")
	buf.Add("<p><a href=\"/login\">Login page.</a></p>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	str := "endPage"
	ctx.AddHtml("pageName", &str)

	content.Include(ctx)
	ctx.Css.Add("/asset/css/page/end.css")
	content.Default(ctx)

	ctx.Render("end.html")
}

func notFound(rw http.ResponseWriter, req *http.Request) {
	if req.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	ctx := context.NewContext(rw, req)
	ctx.Rw.WriteHeader(http.StatusNotFound)

	buf := util.NewBuf()
	buf.Add("<div class=\"callout calloutError\">")
	buf.Add("<h3>404 error.</h3>")
	buf.Add("<p><img src=\"/asset/img/404.jpg\" width=\"100%%\"></p>")
	buf.Add("<p>The page you are looking for does not exists.</p>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	str := "endPage"
	ctx.AddHtml("pageName", &str)

	content.Include(ctx)
	ctx.Css.Add("/asset/css/page/end.css")
	content.Default(ctx)

	ctx.Render("end.html")
}

func cron() {
	//clear old sessions
	go func() {
		c := time.Tick(time.Duration(app.Ini.CookieExpires) * time.Second)
		for _ = range c {
			now := time.Now()
			epoch := now.Unix()

			tx, err := app.Db.Begin()
			if err != nil {
				panic(err)
			}

			sqlStr := `delete from 	userSession
							where (? - userSession.recordTime) > ?`
			_, err = tx.Exec(sqlStr, epoch, app.Ini.CookieExpires)
			if err != nil {
				tx.Rollback()
				log.Printf("%s %s", "error while deleting expired web sessions:", err.Error())
				return
			}

			tx.Commit()
		}
	}()
}

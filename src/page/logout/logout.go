package logout

import (
	"net/http"

	"farlogin/src/app"
	"farlogin/src/content"
	"farlogin/src/lib/context"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsLoggedIn() {
		app.BadRequest()
	}

	ctx.DeleteSession()

	content.End(ctx, "success", "Goodbye", app.SiteName,
		"<p>Your session has been closed.</p>",
		"<a href=\"/login\">Login page.</a>")
}

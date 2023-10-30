package logoutcontroller

import (
	"net/http"

	"github.com/kataras/go-sessions"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	session.Clear()
	sessions.Destroy(w, r)
	http.Redirect(w, r, "/home", http.StatusFound)
}

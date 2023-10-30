package logincontroller

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/kataras/go-sessions"
	"github.com/najuwa28/web_organized/entities"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) != 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if r.Method != "POST" {

		tmpl, err := template.ParseFiles("views/home/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println(username)
	fmt.Println(password)

	users := entities.QueryUser(username)
	admins := entities.QueryAdmin(username)

	var passwordMatch bool
	var role int

	if (entities.User{}) != users {
		passwordMatch = (bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(password)) == nil)
		role = users.Role
	} else if (entities.Admin{}) != admins {
		passwordMatch = (bcrypt.CompareHashAndPassword([]byte(admins.Password), []byte(password)) == nil)
		role = 1
	}

	if passwordMatch {
		// Login success
		session := sessions.Start(w, r)
		session.Set("username", username)
		session.Set("role", role)
		if role == 1 {
			// Admin login
			http.Redirect(w, r, "/home_admin", http.StatusSeeOther)
			fmt.Println("Sukses")
		} else {
			// User login
			http.Redirect(w, r, "/home_user", http.StatusSeeOther)
			fmt.Println("Sukses")
		}
	} else {
		// Login failed
		fmt.Println("Gagal, username atau password salah")
		fmt.Fprint(w, "Gagal")
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

package admincontroller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/kataras/go-sessions"
	"github.com/najuwa28/web_organized/config"
	"github.com/najuwa28/web_organized/entities"
	"github.com/najuwa28/web_organized/util"
	"golang.org/x/crypto/bcrypt"
)

func HomeAdmin(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/views/admin/home_admin", http.StatusMovedPermanently)
	}

	var data = map[string]string{
		"username": session.GetString("username"),
	}
	var t, err = template.ParseFiles("views/admin/home_admin.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t.Execute(w, data)
}

func RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	db := config.Connect_DB()

	if r.Method == "GET" {
		temp, err := template.ParseFiles("views/admin/register_admin.html")
		if err != nil {
			panic(err)
		}

		temp.Execute(w, nil)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	nim := r.FormValue("nim")
	nama := r.FormValue("nama")
	asal_instansi := r.FormValue("asal_instansi")
	mulai_pkl := r.FormValue("mulai_pkl")
	selesai_pkl := r.FormValue("selesai_pkl")
	upload_file := r.FormValue("upload_file")
	role := r.FormValue("role")
	status := r.FormValue("status")
	respon := r.FormValue("respon")
	sertifikat := r.FormValue("sertifikat")

	users := entities.QueryUser(username)

	if (entities.User{}) == users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if len(hashedPassword) != 0 && util.CheckErr(w, r, err) {
			stmt, err := db.Prepare("INSERT INTO users SET username=?, password=?, nim=?, nama=?, asal_instansi=?, mulai_pkl=?, selesai_pkl=?, upload_file=?, role=?, status=?, respon=?, sertifikat=?")
			if err == nil {
				_, err := stmt.Exec(&username, &hashedPassword, &nim, &nama, &asal_instansi, &mulai_pkl, &selesai_pkl, &upload_file, &role, &status, &respon, &sertifikat)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				http.Redirect(w, r, "/data_peserta", http.StatusSeeOther)
				fmt.Println("Berhasil mendaftarkan akun")
				return
			}
		}
	} else {
		http.Redirect(w, r, "/register_admin", http.StatusSeeOther)
		fmt.Println("Gagal mendaftarkan akun (username sudah digunakan)")
		return
	}
}

func ProfileAdmin(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	username := session.GetString("username")

	if len(username) == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	u, err := GetUserByUsername(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data = map[string]interface{}{
		"username": u.Username,
		"password": u.Password,
		"nim":      u.Nim,
		"nama":     u.Nama,
	}

	var t *template.Template
	t, err = template.ParseFiles("views/admin/profile_admin.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)
}

func EditProfileAdmin(w http.ResponseWriter, r *http.Request) {
	db := config.Connect_DB()
	session := sessions.Start(w, r)
	username := session.GetString("username")

	if len(username) == 0 {
		http.Redirect(w, r, "/profile_admin", http.StatusFound)
		return
	}

	u, err := GetUserByUsername(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		newUsername := r.FormValue("new_username")
		newPassword := r.FormValue("new_password")

		if newUsername != "" {
			_, err := db.Exec("UPDATE users SET username = ? WHERE username = ?", newUsername, username)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			username = newUsername
			session.Set("username", newUsername)
		}

		if newPassword != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = db.Exec("UPDATE users SET password = ? WHERE username = ?", hashedPassword, username)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	var data = map[string]interface{}{
		"username": u.Username,
		// "password":     u.Password,
		"nim":  u.Nim,
		"nama": u.Nama,
	}

	var t *template.Template
	t, err = template.ParseFiles("views/admin/edit_profile_admin.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)
}

func GetUserByUsername(username string) (entities.User, error) {
	db := config.Connect_DB()
	var u entities.User
	err := db.QueryRow("SELECT * FROM users WHERE username = ?", username).
		Scan(&u.ID, &u.Username, &u.Password, &u.Nim, &u.Nama, &u.AsalInstansi,
			&u.MulaiPkl, &u.SelesaiPkl, &u.UploadFile, &u.Role, &u.Status, &u.Respon, &u.Sertifikat)
	if err != nil {
		return u, err
	}
	return u, nil
}

type StatusOption struct {
	Value string
	Label string
}

func GetDataPeserta(w http.ResponseWriter, r *http.Request) {
	db := config.Connect_DB()
	var u entities.User
	statusFilter := r.FormValue("status")

	query := "SELECT * FROM users WHERE role = 2"
	if statusFilter != "" && statusFilter != "all" {
		query += " AND status = " + statusFilter
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	userx := []entities.User{}

	for rows.Next() {
		err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.Nim, &u.Nama, &u.AsalInstansi,
			&u.MulaiPkl, &u.SelesaiPkl, &u.UploadFile, &u.Role, &u.Status, &u.Respon, &u.Sertifikat)
		if err != nil {
			log.Fatal(err)
		}
		userx = append(userx, u)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Define the status options
	statusOptions := []StatusOption{
		{Value: "all", Label: "All"},
		{Value: "0", Label: "Tidak Aktif"},
		{Value: "1", Label: "Aktif"},
		{Value: "2", Label: "Selesai"},
	}

	tmpl, err := template.ParseFiles("views/admin/data_peserta.html")
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		DtEmp         []entities.User
		StatusFilter  string
		StatusOptions []StatusOption
	}{
		DtEmp:         userx,
		StatusFilter:  statusFilter,
		StatusOptions: statusOptions,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}
func EditUser(w http.ResponseWriter, r *http.Request) {
	db := config.Connect_DB()
	username := r.URL.Query().Get("username")

	if len(username) == 0 {
		http.Redirect(w, r, "/data_peserta", http.StatusFound)
		fmt.Println("Gagal mengubah data")
		return
	}

	u, err := GetUserByUsername(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		newUsername := r.FormValue("new_username")

		if newUsername != "" {
			_, err := db.Exec("UPDATE users SET username = ? WHERE username = ?", newUsername, username)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}
	}

	var data = map[string]interface{}{
		"username": u.Username,
		// "password":      u.Password,
		"nim":           u.Nim,
		"nama":          u.Nama,
		"asal_instansi": u.AsalInstansi,
		"mulai_pkl":     u.MulaiPkl,
		"selesai_pkl":   u.SelesaiPkl,
		// "upload_file":   u.UploadFile,
		"role":       u.Role,
		"status":     u.Status,
		"respon":     u.Respon,
		"sertifikat": u.Sertifikat,
	}

	var t *template.Template
	t, err = template.ParseFiles("views/admin/edit_user.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func EditAdmin(w http.ResponseWriter, r *http.Request) {
	db := config.Connect_DB()
	username := r.URL.Query().Get("username")

	if len(username) == 0 {
		http.Redirect(w, r, "/data_peserta", http.StatusFound)
		fmt.Println("Gagal mengubah data")
		return
	}

	u, err := GetUserByUsername(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		newUsername := r.FormValue("new_username")
		newPassword := r.FormValue("new_password")
		// newStatus := r.FormValue("new_status")

		if newUsername != "" {
			_, err := db.Exec("UPDATE users SET username = ? WHERE username = ?", newUsername, username)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}

		if newPassword != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = db.Exec("UPDATE users SET password = ? WHERE username = ?", hashedPassword, username)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	var data = map[string]interface{}{
		"username": u.Username,
		// "password":     u.Password,
		"nim":  u.Nim,
		"nama": u.Nama,
	}

	var t *template.Template
	t, err = template.ParseFiles("views/admin/edit_admin.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func UploadSertif(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseMultipartForm(10 << 20) // 10 MB limit for file size

		username := r.URL.Query().Get("username")
		if len(username) == 0 {
			http.Redirect(w, r, "/data_peserta", http.StatusFound)
			fmt.Println("Gagal mengubah data")
			return
		}

		_, err := GetUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Retrieve the uploaded file
		file, handler, err := r.FormFile("sertifikat")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// You can save the uploaded file to a location on your server
		// For example, you can save it to a 'uploads' folder
		// Make sure the 'uploads' folder exists on your server
		f, err := os.Create("sertif/" + handler.Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		// Update the sertifikat file path in the database
		db := config.Connect_DB()
		_, err = db.Exec("UPDATE users SET sertifikat = ? WHERE username = ?", "sertif/"+handler.Filename, username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/data_peserta", http.StatusFound)
		return
	}

	username := r.URL.Query().Get("username")

	if len(username) == 0 {
		http.Redirect(w, r, "/data_peserta", http.StatusFound)
		fmt.Println("Gagal mengubah data")
		return
	}

	u, err := GetUserByUsername(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data = map[string]interface{}{
		"username":   u.Username,
		"sertifikat": u.Sertifikat,
	}

	var t *template.Template
	t, err = template.ParseFiles("views/admin/upload_sertif.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db := config.Connect_DB()

	username := r.URL.Query().Get("username")

	_, err := db.Exec("DELETE FROM users WHERE username = ?", username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/data_peserta", http.StatusSeeOther)
}

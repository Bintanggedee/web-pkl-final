package usercontroller

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/kataras/go-sessions"
	"github.com/najuwa28/web_organized/config"
	"github.com/najuwa28/web_organized/entities"
	"golang.org/x/crypto/bcrypt"
)

func HomeUser(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/home_user", http.StatusMovedPermanently)
	}

	var data = map[string]string{
		"username": session.GetString("username"),
	}
	var t, err = template.ParseFiles("views/user/home_user.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t.Execute(w, data)
}

func ProfileUser(w http.ResponseWriter, r *http.Request) {
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
		"username":      u.Username,
		"password":      u.Password,
		"nim":           u.Nim,
		"nama":          u.Nama,
		"asal_instansi": u.AsalInstansi,
		"mulai_pkl":     u.MulaiPkl,
		"selesai_pkl":   u.SelesaiPkl,
		"upload_file":   u.UploadFile,
		"role":          u.Role,
		"status":        u.Status,
		"respon":        u.Respon,
		"sertifikat":    u.Sertifikat,
	}

	var t *template.Template
	t, err = template.ParseFiles("views/user/profile_user.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)
}

func EditProfileUser(w http.ResponseWriter, r *http.Request) {
	db := config.Connect_DB()
	session := sessions.Start(w, r)
	username := session.GetString("username")

	if len(username) == 0 {
		http.Redirect(w, r, "/profile_user", http.StatusFound)
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

		// //upload file
		// // Mengecek apakah folder pengguna sudah ada, jika tidak, maka membuatnya
		// userFolder := fmt.Sprintf("uploads/%s", username)
		// if _, err := os.Stat(userFolder); os.IsNotExist(err) {
		// 	os.Mkdir(userFolder, os.ModePerm)
		// }

		// // Menerima file yang diunggah
		// file, handler, err := r.FormFile("upload_file")
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		// defer file.Close()

		// // Menyimpan file ke folder pengguna dengan nama yang unik
		// fileName := filepath.Join(userFolder, fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename))
		// f, err := os.Create(fileName)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		// defer f.Close()

		// _, err = io.Copy(f, file)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// // Simpan path file ke dalam database
		// if fileName != "" {
		// 	_, err = db.Exec("UPDATE users SET upload_file = ? WHERE username = ?", fileName, username)
		// 	if err != nil {
		// 		http.Error(w, err.Error(), http.StatusInternalServerError)
		// 		return
		// 	}

		// }
		// //upload file

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
		//"password":      u.Password,w
		"nim":           u.Nim,
		"nama":          u.Nama,
		"asal_instansi": u.AsalInstansi,
		"mulai_pkl":     u.MulaiPkl,
		"selesai_pkl":   u.SelesaiPkl,
		"upload_file":   u.UploadFile,
		"role":          u.Role,
		"status":        u.Status,
		"respon":        u.Respon,
		// "sertifikat":    u.Sertifikat,
	}

	var t *template.Template
	t, err = template.ParseFiles("views/user/edit_profile_user.html")
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

package usermodels

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/kataras/go-sessions"
	"github.com/najuwa28/web_organized/config"
	"golang.org/x/crypto/bcrypt"
)

func GetExistingFilePath(username string) (string, error) {
	db := config.Connect_DB()

	var existingFilePath string
	err := db.QueryRow("SELECT upload_file FROM users WHERE username = ?", username).Scan(&existingFilePath)
	if err != nil {
		return "", err
	}

	return existingFilePath, nil
}

func SaveProfileUser(w http.ResponseWriter, r *http.Request) {
	db := config.Connect_DB()
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/edit_profile_user", http.StatusFound)
		return
	}

	session := sessions.Start(w, r)
	username := session.GetString("username")

	if len(username) == 0 {
		http.Redirect(w, r, "/home_user", http.StatusFound)
		return
	}

	Nim := r.FormValue("nim")
	Nama := r.FormValue("nama")
	AsalInstansi := r.FormValue("asal_instansi")
	MulaiPkl := r.FormValue("mulai_pkl")
	SelesaiPkl := r.FormValue("selesai_pkl")
	// fileName := r.FormValue("upload_file")
	Role := r.FormValue("role")
	Status := r.FormValue("status")
	Respon := r.FormValue("respon")
	// Sertifikat := r.FormValue("sertifikat")
	// Status, err := strconv.Atoi(StatusStr)

	layout := "2006-01-02"
	mulaiPkl, err := time.Parse(layout, MulaiPkl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	selesaiPkl, err := time.Parse(layout, SelesaiPkl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newUsername := r.FormValue("new_username")
	newPassword := r.FormValue("new_password")

	//upload file3
	userFolder := fmt.Sprintf("uploads/%s", username)
	if err := os.MkdirAll(userFolder, os.ModePerm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//upload file3
	var fileName string
	file, handler, err := r.FormFile("upload_file")
	if err != nil {
		// Tidak ada file baru yang diunggah, gunakan file yang sudah ada
		existingFilePath, err := GetExistingFilePath(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fileName = existingFilePath
	} else {
		defer file.Close()

		fileName = filepath.Join(userFolder, fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename))
		f, err := os.Create(fileName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	//upload file3

	if newUsername != "" || newPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`
			UPDATE users
			SET username=?, password=?, nim=?, nama=?, asal_instansi=?, mulai_pkl=?, selesai_pkl=?, upload_file=?, role=?, status=?, respon=?
			WHERE username=?
		`,
			newUsername,
			hashedPassword,
			Nim,
			Nama,
			AsalInstansi,
			mulaiPkl,
			selesaiPkl,
			fileName,
			Role,
			Status,
			Respon,
			// Sertifikat,
			username,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Clear()
		sessions.Destroy(w, r)
		http.Redirect(w, r, "/profile_user", http.StatusSeeOther)
		fmt.Println("Data berhasil diubah")
		return
	}

	// affectedRows, _ := result.RowsAffected()
	// if affectedRows == 0 {
	// 	http.Error(w, "Tidak ada perubahan", http.StatusInternalServerError)
	// 	return
	// }
	http.Redirect(w, r, "/edit_profile_user", http.StatusSeeOther)
	fmt.Print("Gagal menambahkan data")

}

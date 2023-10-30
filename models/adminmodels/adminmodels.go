package adminmodels

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/kataras/go-sessions"
	"github.com/najuwa28/web_organized/config"
	"golang.org/x/crypto/bcrypt"
)

func SaveProfileAdmin(w http.ResponseWriter, r *http.Request) {
	db := config.Connect_DB()
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/edit_profileAdmin", http.StatusFound)
		return
	}

	session := sessions.Start(w, r)
	username := session.GetString("username")

	if len(username) == 0 {
		http.Redirect(w, r, "/home_admin", http.StatusFound)
		return
	}

	Nim := r.FormValue("nim")

	newUsername := r.FormValue("new_username")
	newPassword := r.FormValue("new_password")

	if newUsername != "" || newPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`
			UPDATE users
			SET username=?, password=?, nim=?
			WHERE username=?
		`,
			newUsername,
			hashedPassword,
			Nim,
			username,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Clear()
		sessions.Destroy(w, r)
		http.Redirect(w, r, "/profile_admin", http.StatusSeeOther)
		return
	}

	// affectedRows, _ := result.RowsAffected()
	// if affectedRows == 0 {
	// 	http.Error(w, "Tidak ada perubahan", http.StatusInternalServerError)
	// 	return
	// }
	http.Redirect(w, r, "/profile_admin", http.StatusSeeOther)
}

func SaveUser(w http.ResponseWriter, r *http.Request) {
	db := config.Connect_DB()
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/edit_user", http.StatusFound)
		return
	}

	username := r.URL.Query().Get("username")

	if len(username) == 0 {
		http.Redirect(w, r, "/data_peserta", http.StatusFound)
		fmt.Println("Gagal mengubah data")
		return
	}

	Nim := r.FormValue("nim")
	Nama := r.FormValue("nama")
	AsalInstansi := r.FormValue("asal_instansi")
	MulaiPkl := r.FormValue("mulai_pkl")
	SelesaiPkl := r.FormValue("selesai_pkl")
	// UploadFile := r.FormValue("upload_file")
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
	// newPassword := r.FormValue("new_password")

	// if newUsername != "" {
	// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	_, err = db.Exec(`
			UPDATE users
			SET username=?, nim=?, nama=?, asal_instansi=?, mulai_pkl=?, selesai_pkl=?, role=?, status=?, respon=?
			WHERE username=?
		`,
		newUsername,
		// hashedPassword,
		Nim,
		Nama,
		AsalInstansi,
		mulaiPkl,
		selesaiPkl,
		// UploadFile,
		Role,
		Status,
		Respon,
		username,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Sukses")
	log.Println("Data berhasil diedit")
	http.Redirect(w, r, "/data_peserta", http.StatusSeeOther)
}

func SaveAdmin(w http.ResponseWriter, r *http.Request) {
	db := config.Connect_DB()
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/edit_admin", http.StatusFound)
		return
	}

	username := r.URL.Query().Get("username")

	if len(username) == 0 {
		http.Redirect(w, r, "/data_peserta", http.StatusFound)
		fmt.Println("Gagal mengubah data")
		return
	}

	Nim := r.FormValue("nim")
	Nama := r.FormValue("nama")

	newUsername := r.FormValue("new_username")
	newPassword := r.FormValue("new_password")

	if newUsername != "" || newPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`
			UPDATE users
			SET username=?, password=?, nim=?, nama=?
			WHERE username=?
		`,
			newUsername,
			hashedPassword,
			Nim,
			Nama,
			username,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	fmt.Println("Sukses")
	log.Println("Data berhasil diedit")
	http.Redirect(w, r, "/data_peserta", http.StatusSeeOther)
}

func GetExistingFilePaths(username string) (string, error) {
	db := config.Connect_DB()

	var existingFilePaths string
	err := db.QueryRow("SELECT sertifikat FROM users WHERE username = ?", username).Scan(&existingFilePaths)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // No rows found, return nil without an error
		}
		return "", err
	}

	return existingFilePaths, nil
}

func SaveSertif(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/upload_sertif", http.StatusFound)
		return
	}

	r.ParseMultipartForm(10 << 20) // 10 MB limit for file size

	username := r.URL.Query().Get("username")
	if len(username) == 0 {
		http.Redirect(w, r, "/data_peserta", http.StatusFound)
		fmt.Println("Gagal mengubah data")
		return
	}

	userFolder := fmt.Sprintf("sertif/%s", username)
	if err := os.MkdirAll(userFolder, os.ModePerm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var sertifName string
	sertif, handler, err := r.FormFile("sertif")
	if err != nil {
		// Tidak ada file baru yang diunggah, gunakan file yang sudah ada
		existingFilePaths, err := GetExistingFilePaths(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sertifName = existingFilePaths
	} else {
		defer sertif.Close()

		sertifName = filepath.Join(userFolder, fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename))
		f, err := os.Create(sertifName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		_, err = io.Copy(f, sertif)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update the sertifikat file path in the database
		db := config.Connect_DB()
		_, err = db.Exec("UPDATE users SET sertifikat = ? WHERE username = ?", sertifName, username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	fmt.Println("Sukses")
	log.Println("Data berhasil diedit")
	http.Redirect(w, r, "/data_peserta", http.StatusSeeOther)
}

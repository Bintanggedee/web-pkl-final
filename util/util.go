package util

import (
	"net/http"
)

func CheckErr(w http.ResponseWriter, r *http.Request, err error) bool {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	return true
}

// func QueryUser(username string) User {
// 	db := config.Connect_DB()
// 	var users = User{}

// 	_ = db.QueryRow(`
// 		SELECT id,
// 		username,
// 		password,
// 		nim,
// 		nama,
// 		asal_instansi,
// 		mulai_pkl,
// 		selesai_pkl,
// 		upload_file,
// 		role,
// 		status
// 		FROM users WHERE username=?
// 		`, username).
// 		Scan(
// 			&users.ID,
// 			&users.Username,
// 			&users.Password,
// 			&users.Nim,
// 			&users.Nama,
// 			&users.AsalInstansi,
// 			&users.MulaiPkl,
// 			&users.SelesaiPkl,
// 			&users.UploadFile,
// 			&users.Role,
// 			&users.Status,
// 		)
// 	return users
// }

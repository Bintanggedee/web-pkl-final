package entities

import (
	"time"

	"github.com/najuwa28/web_organized/config"
)

type User struct {
	ID           int
	Username     string
	Password     string
	Nim          string
	Nama         string
	AsalInstansi string
	MulaiPkl     time.Time
	SelesaiPkl   time.Time
	UploadFile   string
	Role         int
	Status       int
	Respon       int
	Sertifikat   string
}

func QueryUser(username string) User {
	var users = User{}
	db := config.Connect_DB()
	_ = db.QueryRow(`
		SELECT id, 
		username, 
		password,
		nim,
		nama,
		asal_instansi,
		mulai_pkl,
		selesai_pkl,
		upload_file,
		role,
		status,
		respon,
		sertifikat
		FROM users WHERE username=?
		`, username).
		Scan(
			&users.ID,
			&users.Username,
			&users.Password,
			&users.Nim,
			&users.Nama,
			&users.AsalInstansi,
			&users.MulaiPkl,
			&users.SelesaiPkl,
			&users.UploadFile,
			&users.Role,
			&users.Status,
			&users.Respon,
			&users.Sertifikat,
		)

	return users
}

package entities

import "github.com/najuwa28/web_organized/config"

type Admin struct {
	Nim      int
	Username string
	Password string
}

func QueryAdmin(username string) Admin {
	var admins = Admin{}
	db := config.Connect_DB()
	_ = db.QueryRow(`
		SELECT nim,
		username,
		password,
		FROM admin WHERE username=?
		`, username).
		Scan(
			&admins.Username,
			&admins.Password,
			&admins.Nim,
		)
	return admins
}

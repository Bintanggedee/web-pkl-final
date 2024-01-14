package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/najuwa28/web_organized/config"
	"github.com/najuwa28/web_organized/controller/admincontroller"
	"github.com/najuwa28/web_organized/controller/homecontroller"
	"github.com/najuwa28/web_organized/controller/logincontroller"
	"github.com/najuwa28/web_organized/controller/logoutcontroller"
	"github.com/najuwa28/web_organized/controller/usercontroller"
	"github.com/najuwa28/web_organized/models/adminmodels"
	"github.com/najuwa28/web_organized/models/usermodels"
)

func main() {
	// 	Database connection
	config.Connect_DB()

	http.Handle("/sertif/", http.StripPrefix("/sertif/", http.FileServer(http.Dir("sertif"))))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("fonts"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.Handle("/video/", http.StripPrefix("/video/", http.FileServer(http.Dir("video"))))

	// Routes
	// 1.Home
	http.HandleFunc("/home", homecontroller.Home)

	// 2. Admin
	http.HandleFunc("/home_admin", admincontroller.HomeAdmin)
	http.HandleFunc("/data_peserta", admincontroller.GetDataPeserta)
	http.HandleFunc("/edit_user", admincontroller.EditUser)
	http.HandleFunc("/edit_admin", admincontroller.EditAdmin)
	http.HandleFunc("/register_admin", admincontroller.RegisterAdmin)
	http.HandleFunc("/profile_admin", admincontroller.ProfileAdmin)
	http.HandleFunc("/edit_profile_admin", admincontroller.EditProfileAdmin)
	http.HandleFunc("/upload_sertif", admincontroller.UploadSertif)
	http.HandleFunc("/delete_user", admincontroller.DeleteUser)

	http.HandleFunc("/save_profile_admin", adminmodels.SaveProfileAdmin)
	http.HandleFunc("/save_user", adminmodels.SaveUser)
	http.HandleFunc("/save_admin", adminmodels.SaveAdmin)
	http.HandleFunc("/save_sertif", adminmodels.SaveSertif)

	// 3. User

	http.HandleFunc("/home_user", usercontroller.HomeUser)
	http.HandleFunc("/profile_user", usercontroller.ProfileUser)
	http.HandleFunc("/edit_profile_user", usercontroller.EditProfileUser)

	http.HandleFunc("/save_profile_user", usermodels.SaveProfileUser)

	// 4. Logout
	http.HandleFunc("/logout", logoutcontroller.Logout)

	// 5.Login
	http.HandleFunc("/login", logincontroller.Login)

	// Run server
	log.Println("Server running on port: 8524")
	log.Fatal(http.ListenAndServe(":8524", nil))
}

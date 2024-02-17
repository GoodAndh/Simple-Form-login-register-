package path

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"webpart2/entities"
	"webpart2/models"

	structtomap "github.com/Klathmon/StructToMap"
)

type Userinput struct {
	Username, Password, Email string
}

var userModel = models.Newusermodel()

func Index(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("example")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	data := map[string]interface{}{
		"Username": username,
	}
	temp, err := template.ParseFiles("views/index.html")
	if err != nil {
		panic(err)
	}
	temp.Execute(w, data)

}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp, err := template.ParseFiles("views/login.html")
		if err != nil {
			panic(err)
		}
		temp.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		r.ParseForm()

		Masukan := &Userinput{
			Username: r.Form.Get("username"),
			Password: r.Form.Get("password"),
		}

		var user entities.User
		userModel.Where(&user, "username", Masukan.Username)

		var message error
		if user.Username == "" {
			message = errors.New("Username atau Password salah")
		} else {
			if user.Password != Masukan.Password {
				message = errors.New("Username atau Password salah")
			}
		}
		if message != nil {
			data := map[string]interface{}{
				"error": message,
			}
			temp, err := template.ParseFiles("views/login.html")
			if err != nil {
				panic(err)
			}
			temp.Execute(w, data)
		} else {
			cookie := http.Cookie{
				Name:     "example",
				Value:    user.NamaLengkap,
				MaxAge:   360,
				Path:     "/",
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("example")
	if err != nil {

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		cookie.MaxAge = -1

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/login", http.StatusSeeOther)

	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		temp, err := template.ParseFiles("views/register.html")
		if err != nil {
			panic(err)
		}
		temp.Execute(w, nil)

	} else if r.Method == http.MethodPost {

		//mengambil inputan form
		r.ParseForm()

		user := entities.User{
			NamaLengkap: r.Form.Get("nama_lengkap"),
			Email:       r.Form.Get("Email"),
			Username:    r.Form.Get("Username"),
			Password:    r.Form.Get("Password"),
			Cpassword:   r.Form.Get("Cpassword"),
		}

		errorMessage := map[string]interface{}{}
		//pengecekan form tidak boleh kosong
		fieldName, _ := structtomap.Convert(user)

		for Fieldnm, FieldValue := range fieldName {

			r.Form.Get(Fieldnm)

			if FieldValue == "" {

				errorMessage[Fieldnm] = "Form ini wajib di isi"
			} else {
				//pengecekan konfirmasi password harus sama dengan password
				if user.Cpassword != user.Password {
					errorMessage["Cpassword"] = "Konfirmasi password harus sama dengan password"
				}
			}
		}
		// pengecekan apa username dan email sudah terdaftar
		result := userModel.SameChecker(&user, user.Username, user.Email)
		if result.Field == "Email" {
			errorMessage["Email"] = "Email ini sudah terdaftar"
		} else if result.Field == "Username" {
			errorMessage["Username"] = "Username ini sudah terdaftar"
		}

		if len(errorMessage) > 0 {

			data := map[string]interface{}{
				"dataError": errorMessage,
			}

			temp, err := template.ParseFiles("views/register.html")
			if err != nil {
				panic(err)
			}
			temp.Execute(w, data)

		} else {

			//no error lanjut
			//insert daatabase
			newCopy := *&user
			userModel.Create(&newCopy)
			fmt.Println(newCopy)

			data := map[string]interface{}{
				"pesan": "Registrasi berhasil",
			}

			temp, err := template.ParseFiles("views/register.html")
			if err != nil {
				panic(err)
			}
			temp.Execute(w, data)

		}
	}
}

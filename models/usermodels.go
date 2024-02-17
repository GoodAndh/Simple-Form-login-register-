package models

import (
	"database/sql"
	"webpart2/config"
	"webpart2/entities"
)

type Usermodel struct {
	db *sql.DB
}

type SameCheckerIsi struct {
	Field string
	Err   error
}


func Newusermodel() *Usermodel {
	conn := config.DBsonek()
	return &Usermodel{
		db: conn,
	}
}

//func untk apa ada data login ada di db(database)
func (u Usermodel) Where(user *entities.User, fieldName, fieldValue string) error {
	row, err := u.db.Query("select id, nama_lengkap ,email, username, password from users where "+fieldName+"= ? limit 1", fieldValue)
	// row, err := u.db.Query("select id,nama_lengkap,email,username,password from users where "+fieldName+" = ?  limit 1", fieldValue)
	if err != nil {
		panic(err)
	}
	defer row.Close()
	for row.Next() {
		row.Scan(&user.Id, &user.NamaLengkap, &user.Email, &user.Username, &user.Password)
	}
	return nil
}

//func untuk insert
func (u Usermodel) Create(user *entities.User) (int64, error) {
	result, err := u.db.Exec("insert into users(nama_lengkap,email,username,password) values(?,?,?,?)", user.NamaLengkap, user.Email, user.Username, user.Password)
	if err != nil {
		return 0, err
	}
	Lastinsert, _ := result.LastInsertId()
	return Lastinsert, nil
}


// func buat mengecek apa ada username dan email yg sama
func (u Usermodel) SameChecker(user *entities.User, kunci, sandi string) SameCheckerIsi {
	roww, err := u.db.Query("select username,email from users where username = ? or email =? limit 1;", kunci, sandi)
	if err != nil {
		return SameCheckerIsi{Err: err}
	}
	defer roww.Close()
	for roww.Next() {
		roww.Scan(&user.Username, &user.Email)
	}
	var field string
	if user.Username == "" {
		field = ""
	} else if user.Username == kunci {
		field = "Username"
	} else if user.Email == sandi {
		field = "Email"
	}

	return SameCheckerIsi{Err: nil, Field: field}
}

package auth

import (
	"app/registry"
	"app/crypto"
	"time"
	"fmt"
	//"encoding/json"
	//"errors"
)

type Auth struct {
	User *User
	SSID string
	Lang string
}

type RegisterForm struct {
	Login string `weaselform:"text" formLabel:"login"`
	Email string `weaselform:"text" formLabel:"Email"`
	Password string `weaselform:"password" formLabel:"Пароль"`
	Password2 string `weaselform:"password" formLabel:"Повторите пароль"`
	UserLastName string `weaselform:"text" formLabel:"Фамилия"`
	UserFirstName  string `weaselform:"text" formLabel:"Имя"`
	UserMiddleName string `weaselform:"text" formLabel:"Отчество"`
	OrganizationName string `weaselform:"text" formLabel:"Организация"`
	OrganizationINN string `weaselform:"text" formLabel:"ИНН"`
	OrganizationKPP string `weaselform:"text" formLabel:"КПП"`
	OrganizationOKOPF uint `weaselform:"text" formLabel:"ОКОПФ"`
}

type User struct {
	UserLastName   string `json:"ul" db:"user_lastname"`
	UserFirstName  string `json:"uf" db:"user_firstname"`
	UserMiddleName string `json:"um" db:"user_middlename"`
	OrganizationId uint   `json:"oi" db:"organization_id"`
	UserID         uint   `json:"i" db:"user_id"`
	IsActive       bool   `json:"a" db:"is_active"`
	Login          string `json:"l" db:"user_login"`
	Email          string `json:"e" db:"user_email"`
	IsAdmin        bool   `json:"adm" db:"is_admin"`
	SessionID      string `json:"-" db:"-"`
	OAuthToken     string `json:"-" db:"-"`
	OAuthRToken    string `json:"-" db:"-"`
}

func AuthUser(login, password string) (*User, error) {

	u := User{}

	salt := ""

	if err := registry.Registry.Connect.SQLX().Get(&salt, "select salt from weasel_auth.users where user_login=$1 and is_active = true"); err != nil {

		time.Sleep(2000 * time.Millisecond)

		return &User{}, err
	}

	p2 := crypto.Encrypt(password, salt)

	if err := registry.Registry.Connect.SQLX().Get(&u, `select user_lastname, user_firstname, user_middlename, user_id, is_active, user_login, user_email, is_admin, organization_id
	from weasel_auth.users where user_login=$1 and user_password=$2 and is_active = true`,
		login,
		p2,
	); err != nil {

		time.Sleep(2000 * time.Millisecond)

		return &User{}, err

	}

	return &User{
		UserLastName : u.UserLastName,
		UserFirstName : u.UserFirstName,
		UserMiddleName : u.UserMiddleName,
		OrganizationId : u.OrganizationId,
		UserID : u.UserID,
		IsActive : u.IsActive,
		Login : u.Login,
		Email : u.Email,
		IsAdmin : u.IsAdmin,
	}, nil

	//return &r, nil

}

func AddUser(r RegisterForm) (uint, error) {

	res := 0

	fmt.Println(r)

	password := crypto.Encrypt(r.Password, "")

	if err := registry.Registry.Connect.SQLX().Get(&res, `select * from weasel_auth.add_user($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		1,
		r.UserFirstName,
		r.UserLastName,
		r.UserMiddleName,
		"job_title",
		"",
		r.Login,
		password,
		1,
		1,
	); err != nil {

		return 0, err

	}

	return uint(res), nil

}

//func (u *User) Scan(src interface {}) error {
//
//	var source []byte
//
//	switch src.(type) {
//
//	case string:
//
//		source = []byte(src.(string))
//
//	case []byte:
//
//		source = src.([]byte)
//
//	default:
//
//		return errors.New("Incompatible type for auth.User")
//	}
//
//	if err := json.Unmarshal(source, &u); err != nil {
//
//		return err
//	}
//
//	return nil
//}
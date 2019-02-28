package common

import (
	"encoding/hex"
	"math/rand"
	"lib/auth"
	"app/mailer"
	"errors"
	"strings"
	"app/crypto"
	"app/registry"
	"strconv"
)

func RegisterUser(nu *auth.RegisterForm) (string, error) {

	u := make([]byte, 24)

	_, err := rand.Read(u)
	if err != nil {
		return "", err
	}

	nu.Login = strings.TrimSpace(nu.Login)

	if len(nu.Login) < 3 {
		return "", errors.New("Login must be at least 4 symbols")
	}

	if len(nu.Password) < 7 {
		return "", errors.New("Password must be at least 8 symbols")
	}

	if nu.Password != nu.Password2 {
		return "", errors.New("Passwords don't match")
	}

	sa := make([]byte, 10)

	_, err = rand.Read(sa)
	if err != nil {
		return "", err
	}

	sa[2] = byte(nu.Login[1])
	sa[7] = byte(nu.Login[0])
	sa[5] = byte(nu.Login[2])

	salt := hex.EncodeToString(sa)

	var ni int

	if err := registry.Registry.Connect.SQLX().Get(&ni, "select * from weasel_auth.add_user($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
		0,
		nu.UserFirstName,
		nu.UserLastName,
		nu.UserMiddleName,
		"",
		"",
		nu.Email,
		nu.Login,
		crypto.Encrypt(nu.Password, salt),
		false,
		1,
		salt); err != nil {
			return "", err
	}

	mailer.SendQueued("/mail/registration.html", map[string]interface{}{
		"userName" : nu.UserFirstName,
		"tmpLink": hex.EncodeToString(u),
	})

	registry.Registry.KVS.Put([]byte("new_user_reg_"+strconv.Itoa(ni)), hex.EncodeToString(u))

	return hex.EncodeToString(u), nil

}

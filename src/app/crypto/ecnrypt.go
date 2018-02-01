package crypto

import (
	"encoding/base64"
	"crypto/sha1"
	"strings"
	"strconv"
	"math/rand"
	"time"
	"fmt"
)

const salt = "o0d*0sfJFMxWea2kd#sel#fajBee"
const key = "0j8j4na__osffs99.03l4"

func Encrypt(s, ls string) string {

	return fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%s%s%s", salt, s, ls))))

}

func Unique() string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%s%s%d", salt, time.Now(), rand.Int()))))
}

func GenSessionId(i uint, l_salt string) string {

	return fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%s%d%s%s%d", salt, i, l_salt, time.Now(), rand.Int()))))

}

func EncryptUrl(o uint) string {
	return strings.TrimRight(EncryptUint(o, key), "=")
}

func DecryptUrl(crypt string) (uint, error) {
	return DecryptUint(crypt, key)
}


func EncryptUint(o uint, key string) string {

	return EncryptB64(fmt.Sprintf("%d", o), key)
}

func DecryptUint(crypt, key string) (uint, error) {

	dec, err := DecryptB64(crypt, key)
	if err != nil {

		return 0, err
	}

	r, err := strconv.ParseUint(dec, 10, 64)
	if err != nil {

		return 0, err
	}

	return uint(r), nil
}

func EncryptB64(original, key string) string {

	kl := len([]byte(key))

	var buffer []byte

	for i := 0; i < len([]byte(original)); i++ {

		buffer = append(buffer, original[i]^key[i%kl])
	}

	return base64.URLEncoding.EncodeToString(buffer)
}

func DecryptB64(crypt, key string) (string, error) {

	kl := len([]byte(key))

	if fix := 4 - (len(crypt) % 4); fix != 4 {

		crypt = fmt.Sprintf("%s%s", crypt, strings.Repeat("=", fix))
	}

	if message, err := base64.URLEncoding.DecodeString(crypt); err == nil {

		var buffer []byte

		for i := 0; i < len([]byte(message)); i++ {

			buffer = append(buffer, message[i]^key[i%kl])
		}

		return string(buffer), nil

	} else {

		return "", err
	}
}

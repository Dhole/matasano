package main

import (
	"./modes"
	"./pkcs7"
	"bytes"
	"crypto/aes"
	crand "crypto/rand"
	"errors"
	"fmt"
	"strings"
)

var key []byte

var blockSize int

func ParseCookie(cookie string) (data map[string]string, err error) {
	data = make(map[string]string)
	pairs := strings.Split(cookie, "&")
	for _, p := range pairs {
		if len(p) < 2 {
			return nil, errors.New("wrong cookie")
		}
		kv := strings.Split(p, "=")
		data[kv[0]] = kv[1]
	}
	return data, nil
}

func GenCookie(data map[string]string) (cookie string) {
	var buffer bytes.Buffer
	for k, v := range data {
		buffer.WriteString(fmt.Sprintf("%s=%s&", k, v))
	}
	if buffer.Len() > 0 {
		buffer.Truncate(buffer.Len() - 1)
	}

	return buffer.String()
}

/*
func ProfileFor(email string) (cookie string) {
	data := make(map[string]string)

	role := "user"
	uid := "10"
	email = strings.Replace(email, "=", "", -1)
	email = strings.Replace(email, "&", "", -1)

	data["role"] = role
	data["uid"] = uid
	data["email"] = email

	return GenCookie(data)
}
*/

func ProfileFor(email string) (cookie string) {
	return fmt.Sprintf("email=%s&uid=10&role=user", email)
}

func GenKey() {
	key = make([]byte, blockSize)
	_, err := crand.Read(key)
	if err != nil {
		panic(err)
	}
	fmt.Println("key set")
}

func EncryptCookie(cookie string) (ciphtxt []byte) {
	clrtxt, err := pkcs7.Pad([]byte(cookie), blockSize)
	if err != nil {
		panic(err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	ciphtxt = make([]byte, len(clrtxt))
	modes.NewECBEncrypter(block).CryptBlocks(ciphtxt, clrtxt)
	return ciphtxt
}

func DecryptCookie(ciphtxt []byte) (cookie string) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	clrtxt := make([]byte, len(ciphtxt))
	modes.NewECBDecrypter(block).CryptBlocks(clrtxt, ciphtxt)
	clrtxt, err = pkcs7.Unpad(clrtxt)
	if err != nil {
		panic(err)
	}
	return string(clrtxt)
}

func GenCiphProfileFor(email string) (ciphtxt []byte) {
	cookie := ProfileFor(email)
	ciphtxt = EncryptCookie(cookie)
	return ciphtxt
}

func main() {
	//cookie := "foo=bar&baz=qux&zap=zazzle"
	//data, _ := ParseCookie(cookie)
	//fmt.Println(data)
	//fmt.Println(GenCookie(data))
	//fmt.Println(ProfileFor("foo@bar.com"))
	// we chose an email such that the plaintext is as follows:
	// b0 = "email=foo12@bar."
	// b1 = "admin           " <- padded
	// b2 = "com&uid=10&role="
	// b3 ="user"
	// we then create a ciphertext with the ciphers of b0+b2+b1
	blockSize = 16

	padded_role, _ := pkcs7.Pad([]byte("admin"), blockSize)
	email := "foo12@bar.com" // Needs to be of size 13
	email_mod := email[:10] + string(padded_role) + email[10:13]
	//fmt.Println("email: ", email)
	GenKey()
	ciphtxt := GenCiphProfileFor(email_mod)

	ciphtxt_mod := make([]byte, blockSize*3)
	copy(ciphtxt_mod[blockSize*0:blockSize*1], ciphtxt[blockSize*0:blockSize*1])
	copy(ciphtxt_mod[blockSize*1:blockSize*2], ciphtxt[blockSize*2:blockSize*3])
	copy(ciphtxt_mod[blockSize*2:blockSize*3], ciphtxt[blockSize*1:blockSize*2])

	fmt.Println("Decription: ", DecryptCookie(ciphtxt_mod))

}

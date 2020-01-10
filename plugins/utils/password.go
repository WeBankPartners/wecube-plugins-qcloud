package utils

import (
	"math/rand"
	"strings"
	"time"
)

const (
	PASSWORD_LEN  = 12
	DEFALT_CIPHER = "CIPHER_A"
)

var CIPHER_MAP = map[string]string{
	"CIPHER_A": "{cipher_a}",
}

func CreateRandomPassword() string {
	digitals := "0123456789"
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(letters)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < PASSWORD_LEN-4; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	bytes = []byte(digitals)
	for i := 0; i < 4; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	return string(result)
}

func AesEnPassword(guid, seed, password, cipher string) (string, error) {
	if cipher == "" {
		cipher = DEFALT_CIPHER
	}
	md5sum := Md5Encode(guid + seed)
	enPassword, err := AesEncode(md5sum[0:16], password)
	if err != nil {
		return "", err
	}
	return CIPHER_MAP[cipher] + enPassword, nil
}

func AesDePassword(guid, seed, password string) (string, error) {
	var cipher string
	for _, _cipher := range CIPHER_MAP {
		if strings.HasPrefix(password, _cipher) {
			cipher = _cipher
			break
		}
	}
	if cipher == "" {
		return password, nil
	}
	password = password[len(cipher)-1:]

	md5sum := Md5Encode(guid + seed)
	dePassword, err := AesDecode(md5sum[0:16], password)
	if err != nil {
		return "", err
	}
	return dePassword, nil
}

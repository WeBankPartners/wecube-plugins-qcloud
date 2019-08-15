package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func Md5Encode(rawData string) string {
	data := []byte(rawData)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length > unpadding {
	  return origData[:(length - unpadding)]
	}
	return  []byte{}
}

func AesEncode(key string, rawData string) (string, error) {
	bytesRawKey := []byte(key)
	block, err := aes.NewCipher(bytesRawKey)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	origData := PKCS7Padding([]byte(rawData), blockSize)
	blockMode := cipher.NewCBCEncrypter(block, bytesRawKey[:blockSize])
	crypted := make([]byte, len([]byte(origData)))
	blockMode.CryptBlocks(crypted, origData)
	return hex.EncodeToString(crypted), nil
}

func AesDecode(key string, encryptData string) (password string, err error) {
	defer func(){
			 if r:=recover();r!=nil{
					err=fmt.Errorf("%v",r)
			}
	}()

    bytesRawKey := []byte(key)
    bytesRawData, _ := hex.DecodeString(encryptData)
    block, err := aes.NewCipher(bytesRawKey)
    if err != nil {
       return 
    }
    blockSize := block.BlockSize()
    blockMode := cipher.NewCBCDecrypter(block, bytesRawKey[:blockSize])
    origData := make([]byte, len(bytesRawData))
    blockMode.CryptBlocks(origData, bytesRawData)
	
	origData = PKCS7UnPadding(origData)
    if len(origData) == 0 {
	   err=fmt.Errorf("password wrong")
	   return 
    }

   password=string(origData)
   return 
}

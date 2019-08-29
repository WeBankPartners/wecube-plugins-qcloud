package main

import (
        "bufio"
        "bytes"
        "crypto/aes"
        "crypto/cipher"
        "crypto/md5"
        "encoding/hex"
        "fmt"
        "log"
        "os"
        "strings"
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
        return origData[:(length - unpadding)]
}

func AesEncode(key string, rawData string) (string, string) {
        bytesRawKey := []byte(key)
        block, err := aes.NewCipher(bytesRawKey)
        if err != nil {
                return "", ""
        }
        blockSize := block.BlockSize()
        origData := PKCS7Padding([]byte(rawData), blockSize)
        blockMode := cipher.NewCBCEncrypter(block, bytesRawKey[:blockSize])
        crypted := make([]byte, len([]byte(origData)))
        blockMode.CryptBlocks(crypted, origData)
        return hex.EncodeToString(crypted), ""
}

func AesDecode(key string, encryptData string) (string, error) {
        bytesRawKey := []byte(key)
        bytesRawData, _ := hex.DecodeString(encryptData)
        block, err := aes.NewCipher(bytesRawKey)
        if err != nil {
                return "", err
        }
        blockSize := block.BlockSize()
        blockMode := cipher.NewCBCDecrypter(block, bytesRawKey[:blockSize])
        origData := make([]byte, len(bytesRawData))
        blockMode.CryptBlocks(origData, bytesRawData)
        origData = PKCS7UnPadding(origData)
        return string(origData), nil
}

func encode(guid string, seed string, password string) string {
        md5sum := Md5Encode(guid + seed)
        encoded, _ := AesEncode(md5sum[0:16], password);
        return encoded
}

func decode(guid string, seed string, encoded string) (string, error){
        md5sum := Md5Encode(guid+seed)
        decode,err := AesDecode(md5sum[0:16], encoded)
        if err != nil {
                log.Println("AesDecode meet error(%v)", err)
                return decode , err
        }
        return decode,nil
}

func main() {
        reader := bufio.NewReader(os.Stdin)
        fmt.Println("AES algorithm:")
        fmt.Println("---------------------")
        for {
                fmt.Print("\nEnc/Dec method: [E/D]-> ")
                method, _ := reader.ReadString('\n')
                method = strings.Replace(method, "\n", "", -1)
                fmt.Print("guid-> ")
                guid, _ := reader.ReadString('\n')
                guid = strings.Replace(guid, "\n", "", -1)
                fmt.Print("seed-> ")
                seed, _ := reader.ReadString('\n')
                seed = strings.Replace(seed, "\n", "", -1)
                if strings.EqualFold(method, "E")  {
                        fmt.Print("password-> ")
                        password, _ := reader.ReadString('\n')
                        password = strings.Replace(password, "\n", "", -1)
                        enc := encode (guid, seed, password)
                        fmt.Print("Encode result: ", enc)
                } else if strings.EqualFold(method, "D") {
                        fmt.Print("password-> ")
                        encoded, _ := reader.ReadString('\n')
                        encoded = strings.Replace(encoded, "\n", "", -1)
                        dec,_ := decode(guid, seed, encoded)
                        fmt.Print("Decode result: ", dec)
                }else {
                        fmt.Print("Please input E or D!")
                }
        }
}



package encrypx

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
)

// AESCBCEncrypt AES_CBC+PKCS5PADDING
func AESCBCEncrypt(data, secret, secretIV string) (string, error) {
	// Checks if secret is a valid secret
	if !isValidSecret(secret) {
		return "", errors.New("secret key must be 16, 24, or 32 bytes in length")
	}
	// Check secretIV's length
	if len(secretIV) != 16 {
		return "", errors.New("secretIV must be 16 bytes in length")
	}
	// Creates a new cipher block
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}
	// Pad data
	dataByte := padPKCS5([]byte(data), block.BlockSize())
	// Creates a new CBC encryptor block
	blockMode := cipher.NewCBCEncrypter(block, []byte(secretIV))
	cipherText := make([]byte, len(dataByte))
	blockMode.CryptBlocks(cipherText, dataByte)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// AESCBCDecrypt AES_CBC+PKCS5PADDING
func AESCBCDecrypt(encryptData, secret, secretIV string) (string, error) {
	// Checks if secret is a valid secret
	if !isValidSecret(secret) {
		return "", errors.New("secret key must be 16, 24, or 32 bytes in length")
	}
	// Check secretIV's length
	if len(secretIV) != 16 {
		return "", errors.New("secretIV must be 16 bytes in length")
	}
	// Creates a new cipher block
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}

	encryptDataByte, err := base64.StdEncoding.DecodeString(encryptData)
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(block, []byte(secretIV))
	tmpData := make([]byte, len(encryptDataByte))
	blockMode.CryptBlocks(tmpData, encryptDataByte)

	tmpData = unPadPKCS5(tmpData)
	// Unquote the escaped string
	oriData, err := unquoteStr(string(tmpData))
	if err != nil {
		return "", err
	}
	return oriData, nil
}

// PKCS5PADDING pad
func padPKCS5(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// PKCS5PADDING unPad
func unPadPKCS5(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

// secret's length must be 16, 24 or 32 bytes
func isValidSecret(secret string) bool {
	sLen := len(secret)
	switch sLen {
	default:
		return false
	case 16, 24, 32:
		return true
	}
}

// Unquote the escaped string
func unquoteStr(str string) (string, error) {
	if strings.HasPrefix(str, `"`) && strings.HasSuffix(str, `"`) {
		return strconv.Unquote(str)
	}
	return str, nil
}

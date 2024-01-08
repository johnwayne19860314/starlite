package encrypx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	Data     = "myTestData"
	Secret   = "1234567890secret"
	SecretIV = "testSecretIv0123"
)

func TestAESCBCEncrypt(t *testing.T) {
	encrypt, err := AESCBCEncrypt(Data, Secret, SecretIV)
	assert.Nil(t, err)
	assert.Equal(t, "vwG4HFhgAksO7dBdLmy2Bw==", encrypt)

	// illegal Secret
	illegalSecret := "illegalSecret"
	_, err = AESCBCEncrypt(Data, illegalSecret, SecretIV)
	assert.NotNil(t, err)

	// illegal SecretIV
	illegalSecretIV := "illegalSecretIV"
	_, err = AESCBCEncrypt(Data, Secret, illegalSecretIV)
	assert.NotNil(t, err)
}

func TestAESCBCDecrypt(t *testing.T) {
	encryptedData, _ := AESCBCEncrypt(Data, Secret, SecretIV)

	decryptedData, err := AESCBCDecrypt(encryptedData, Secret, SecretIV)
	assert.Nil(t, err)
	assert.Equal(t, Data, decryptedData)

	// illegal Secret
	illegalSecret := "illegalSecret"
	_, err = AESCBCDecrypt(encryptedData, illegalSecret, SecretIV)
	assert.NotNil(t, err)

	// illegal SecretIV
	illegalSecretIV := "illegalSecretIV"
	_, err = AESCBCDecrypt(encryptedData, Secret, illegalSecretIV)
	assert.NotNil(t, err)

	// illegal encryptedData
	illegalEncryptedData := "illegal encrypted Data"
	_, err = AESCBCDecrypt(illegalEncryptedData, Secret, SecretIV)
	assert.NotNil(t, err)
}

func TestIsValidSecret(t *testing.T) {
	// Testing valid Secret
	validSecret := "1234567890secret"
	isValid := isValidSecret(validSecret)
	assert.True(t, isValid)

	// Testing invalid Secret
	invalidSecret := "thissecretislongerthan32bits"
	isInvalid := isValidSecret(invalidSecret)
	assert.False(t, isInvalid)
}

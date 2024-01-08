package signature

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSign_success(t *testing.T) {
	sigSecret := "sigSecret"
	info := "info"
	expectedSignature := "324C02D2AF84DC0F37CEAB9991F10144"

	signature, err := Sign(sigSecret, info)
	assert.Nil(t, err, "Expected no error, but got: %v", err)
	assert.Equal(t, expectedSignature, signature, "The generated signature doesn't match the expected signature")
}

func TestSign_fail(t *testing.T) {
	// sigSecret is empty
	_, err := Sign("", "info")
	assert.NotNil(t, err, "Expected error but got nil")
	assert.Equal(t, errors.New("both sigSecret and sigData must be non-empty"),
		err, "Expected 'both sigSecret and sigData must be non-empty' error")

	// sigData is empty
	_, err = Sign("sigSecret", "")
	assert.NotNil(t, err, "Expected error but got nil")
	assert.Equal(t, errors.New("both sigSecret and sigData must be non-empty"), err,
		"Expected 'both sigSecret and sigData must be non-empty' error")

}

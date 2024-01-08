package stringutil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasText(t *testing.T) {
	assertive := assert.New(t)
	assertive.Equal(false, HasText(""), "An empty string should return false")
	assertive.Equal(false, HasText("   "), "A string with only spaces should return false")
	assertive.Equal(true, HasText("Hello"), "A string with none-space-characters should return true")
	assertive.Equal(true, HasText(" Hello  "), "A string with spaces at both ends and none-space-characters should return true")
}

func TestAppendError(t *testing.T) {
	assertive := assert.New(t)
	errMsg := "An error occurred"
	err := errors.New("fatal error")
	assertive.Equal("Some message", AppendError("Some message", nil), "Without error append, message should stay the same")
	assertive.Equal(errMsg+"; error: "+err.Error(), AppendError(errMsg, err), "With error append, message should append error message")
}

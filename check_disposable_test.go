package mailck

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckDisposable(t *testing.T) {
	assert.False(t, CheckDisposable("sebastian@mancke.net"))
	assert.True(t, CheckDisposable("foo@mailinator.com"))
}

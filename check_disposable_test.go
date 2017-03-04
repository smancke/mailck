package mailck

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCheckDisposable(t *testing.T) {
	start := time.Now()

	assert.False(t, CheckDisposable("sebastian@mancke.net"))
	assert.True(t, CheckDisposable("foo@mailinator.com"))

	fmt.Printf("check for %v\n", time.Since(start))
}

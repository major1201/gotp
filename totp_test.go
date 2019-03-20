package gotp

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTotp_Generate(t *testing.T) {
	ta := assert.New(t)
	totp := NewDefaultTotp("4S62BZNFXXSZLCRO")
	val, _ := totp.GenerateWithRemainingSeconds(time.Unix(1524485781, 0))
	ta.Equal(uint32(179394), val)
}

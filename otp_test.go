package gotp

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewOtpFromURI(t *testing.T) {
	ta := assert.New(t)
	totp, err := NewOtpFromURI("otpauth://totp/github:major1201?secret=4S62BZNFXXSZLCRO&issuer=github")
	ta.Nil(err)
	val, _ := totp.(*Totp).GenerateWithRemainingSeconds(time.Unix(1524485781, 0))
	ta.Equal(uint32(179394), val)
}

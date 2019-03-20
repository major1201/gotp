package gotp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHotp_Generate(t *testing.T) {
	ta := assert.New(t)
	hotp := NewDefaultHotp("4S62BZNFXXSZLCRO", 12345)
	ta.Equal(uint32(194001), hotp.Generate())
}

package gotp

import (
	"crypto/sha1"
	"fmt"
	"hash"
	"strings"
	"time"
)

// Totp defines an TOTP object
type Totp struct {
	otp
	period uint64
}

// NewDefaultTotp creates an Totp object with default arguments
func NewDefaultTotp(secret string) *Totp {
	return NewTotp("", "", secret, sha1.New, "sha1", 6, 30)
}

// NewTotp creates an Totp object with detailed arguments
func NewTotp(accountName, issuer, secret string, algorithm func() hash.Hash, algoString string, digits uint8, period uint64) *Totp {
	return &Totp{
		otp: otp{
			accountName: accountName,
			issuer:      issuer,
			secret:      secret,
			algorithm:   algorithm,
			algoString:  algoString,
			digits:      digits,
		},
		period: period,
	}
}

// Type returns "totp"
func (o *Totp) Type() string {
	return "totp"
}

// GenerateWithRemainingSeconds generates an TOTP value with the remaining seconds
func (o *Totp) GenerateWithRemainingSeconds(t time.Time) (uint32, uint64) {
	timeUnix := uint64(t.Unix())
	return o.At(timeUnix / o.period), o.period - (timeUnix % o.period)
}

// Generate an OTP value
func (o *Totp) Generate() uint32 {
	result, _ := o.GenerateWithRemainingSeconds(time.Now())
	return result
}

// URI reassembles the otp uri
func (o *Totp) URI() string {
	label, parameters := o.uriFragments()
	parameters = append(parameters, fmt.Sprintf("period=%v", o.period))
	return fmt.Sprintf("otpauth://%s/%s?%s", o.Type(), label, strings.Join(parameters, "&"))
}

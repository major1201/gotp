package gotp

import (
	"crypto/sha1"
	"fmt"
	"hash"
	"strings"
)

// Hotp defines an HOTP object
type Hotp struct {
	otp
	counter uint64
}

// NewDefaultHotp creates an Hotp object with default arguments
func NewDefaultHotp(secret string, counter uint64) *Hotp {
	return NewHotp("", "", secret, sha1.New, "sha1", 6, counter)
}

// NewHotp creates an Hotp object with detailed arguments
func NewHotp(accountName, issuer, secret string, algorithm func() hash.Hash, algoString string, digits uint8, counter uint64) *Hotp {
	return &Hotp{
		otp: otp{
			accountName: accountName,
			issuer:      issuer,
			secret:      secret,
			algorithm:   algorithm,
			algoString:  algoString,
			digits:      digits,
		},
		counter: counter,
	}
}

// Type returns "hotp"
func (o *Hotp) Type() string {
	return "hotp"
}

// Generate an OTP value
func (o *Hotp) Generate() uint32 {
	v := o.At(o.counter)
	o.counter++
	return v
}

// URI reassembles the otp uri
func (o *Hotp) URI() string {
	label, parameters := o.uriFragments()
	parameters = append(parameters, fmt.Sprintf("counter=%v", o.counter))
	return fmt.Sprintf("otpauth://%s/%s?%s", o.Type(), label, strings.Join(parameters, "&"))
}

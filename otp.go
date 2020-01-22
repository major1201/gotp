package gotp

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/major1201/goutils"
	"hash"
	"math"
	"net/url"
	"strconv"
	"strings"
)

// Otp defines an OTP interface
type Otp interface {
	ID() uint64
	SetID(id uint64)
	Type() string
	AccountName() string
	Issuer() string
	Digits() uint8
	URI() string

	At(uint64) uint32
	Generate() uint32
	ToString(uint32) string
}

type otp struct {
	id          uint64
	accountName string
	issuer      string
	secret      string
	algorithm   func() hash.Hash
	algoString  string
	digits      uint8
}

// ID returns the OTP id
func (o *otp) ID() uint64 {
	return o.id
}

// SetID sets the OTP id
func (o *otp) SetID(id uint64) {
	o.id = id
}

// AccountName returns the OTP account name
func (o *otp) AccountName() string {
	return o.accountName
}

// Issuer returns the OTP issuer
func (o *otp) Issuer() string {
	return o.issuer
}

// Issuer returns the OTP digits
func (o *otp) Digits() uint8 {
	return o.digits
}

// At actually generate the OTP value
func (o *otp) At(c uint64) uint32 {
	hm := hmac.New(o.algorithm, o.byteSecret())
	buf := make([]byte, 8) // uint64 length is 8
	binary.BigEndian.PutUint64(buf, c)
	hm.Write(Itob(c))
	hmacHash := hm.Sum(nil)

	offset := hmacHash[len(hmacHash)-1] & 0x0f
	dbc1 := hmacHash[offset : offset+4]
	dbc1[0] = dbc1[0] & 0x7f
	dbc2 := binary.BigEndian.Uint32(dbc1)
	return uint32(dbc2 % uint32(math.Pow10(int(o.digits))))
}

// ToString convert the generated OTP value to string
func (o *otp) ToString(d uint32) string {
	return fmt.Sprintf(fmt.Sprintf("%%0%dd", o.digits), d)
}

func (o *otp) byteSecret() []byte {
	missingPadding := len(o.secret) % 8
	if missingPadding != 0 {
		o.secret = o.secret + strings.Repeat("=", 8-missingPadding)
	}
	bytes, err := base32.StdEncoding.DecodeString(o.secret)
	if err != nil {
		panic("decode secret failed")
	}
	return bytes
}

func (o *otp) uriFragments() (label string, parameters []string) {
	// assemble label
	if o.issuer == "" {
		label = url.QueryEscape(o.accountName)
	} else {
		label = fmt.Sprintf("%s:%s", url.QueryEscape(o.issuer), url.QueryEscape(o.accountName))
	}

	// assemble parameters
	parameters = append(parameters,
		fmt.Sprintf("secret=%s", o.secret),
		fmt.Sprintf("digits=%v", o.digits),
		fmt.Sprintf("algorithm=%s", o.algoString),
	)
	if o.issuer != "" {
		parameters = append(parameters, fmt.Sprintf("issuer=%s", url.QueryEscape(o.issuer)))
	}

	return
}

// NewOtpFromURI makes an OTP object from a URI, ref. <https://github.com/google/google-authenticator/wiki/Key-Uri-Format>
func NewOtpFromURI(uri string) (Otp, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "otpauth" {
		return nil, errors.New("uri scheme should start with otpauth://")
	}

	path := strings.TrimLeft(u.Path, "/")
	labels := strings.SplitN(path, ":", 2)
	var accountName, issuer string
	if len(labels) == 1 {
		accountName = labels[0]
	} else {
		issuer = labels[0]
		accountName = strings.TrimLeft(labels[1], " ")
	}

	values := u.Query()

	secret := strings.ToUpper(values.Get("secret"))
	if secret == "" {
		return nil, errors.New("secret is empty")
	}

	if values.Get("issuer") != "" {
		issuer = values.Get("issuer")
	}

	var algorithm func() hash.Hash
	algoString := strings.ToLower(values.Get("algorithm"))
	switch algoString {
	case "":
		algoString = "sha1"
		algorithm = sha1.New
	case "sha1":
		algorithm = sha1.New
	case "sha256":
		algorithm = sha256.New
	case "sha512":
		algorithm = sha512.New
	default:
		return nil, errors.New("algorithm should be in sha1, sha256, sha512, or leave empty as default(sha1)")
	}

	digits := uint8(goutils.ToIntDv(values.Get("digits"), 6))

	switch u.Host {
	case "hotp":
		var counter uint64
		if values.Get("counter") == "" {
			counter = 0 // default value
		} else {
			counter, err = strconv.ParseUint(values.Get("counter"), 10, 64)
			if err != nil {
				return nil, errors.New("error parsing argument: counter")
			}
		}
		return NewHotp(accountName, issuer, secret, algorithm, algoString, digits, counter), nil
	case "totp":
		var period uint64
		if values.Get("period") == "" {
			period = 30 // default value
		} else {
			period, err = strconv.ParseUint(values.Get("period"), 10, 64)
			if err != nil {
				return nil, errors.New("error parsing argument: period")
			}
		}
		return NewTotp(accountName, issuer, secret, algorithm, algoString, digits, period), nil
	default:
		return nil, errors.New("uri type should be hotp or totp")
	}
}

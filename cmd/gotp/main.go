package main

import (
	"errors"
	"fmt"
	"github.com/major1201/gotp"
	"github.com/major1201/goutils"
	"github.com/urfave/cli"
	"os"
	"strconv"
	"strings"
	"time"
)

// AppVer means the project's version
const AppVer = "0.1.0"

func runApp(c *cli.Context) error {
	s, err := NewStore(c.String("database"))
	if err != nil {
		return err
	}
	defer s.Close()

	otps, err := s.List()
	if err != nil {
		return err
	}

	now := time.Now() // use unified time
	for _, otp := range otps {
		var fmtStrArr []string
		var args []interface{}
		if c.Bool("id") {
			fmtStrArr = append(fmtStrArr, "<%d>: ")
			args = append(args, otp.ID())
		}
		fmtStrArr = append(fmtStrArr, "%s(%s): %s")
		args = append(args, otp.Issuer(), otp.AccountName())
		if otp.Type() == "totp" {
			fmtStrArr = append(fmtStrArr, " (%ds remaining)")
			val, remaining := otp.(*gotp.Totp).GenerateWithRemainingSeconds(now)
			args = append(args, otp.ToString(val), remaining)
		} else {
			args = append(args, otp.ToString(otp.Generate()))
		}
		fmtStrArr = append(fmtStrArr, "\n")
		fmt.Printf(strings.Join(fmtStrArr, ""), args...)
	}
	return nil
}

func runSubcommandAdd(c *cli.Context) error {
	s, err := NewStore(c.Parent().String("database"))
	if err != nil {
		return err
	}
	defer s.Close()

	if c.IsSet("uri") {
		otp, err := gotp.NewOtpFromURI(c.String("uri"))
		if err != nil {
			return err
		}

		return s.Add(otp)
	}

	accountName := c.String("accountname")
	if goutils.IsBlank(accountName) {
		return errors.New("accountname should not be blank")
	}

	secret := c.String("secret")
	if goutils.IsBlank(secret) {
		return errors.New("secret should not be blank")
	}

	algo, err := gotp.ConvAlgoString(c.String("algorithm"))
	if err != nil {
		return err
	}

	switch c.String("type") {
	case "totp":
		return s.Add(gotp.NewTotp(
			accountName, c.String("issuer"), secret, algo,
			c.String("algorithm"), uint8(c.Uint("digits")), c.Uint64("period")),
		)
	case "hotp":
		return s.Add(gotp.NewHotp(
			accountName, c.String("issuer"), secret, algo,
			c.String("algorithm"), uint8(c.Uint("digits")), c.Uint64("counter")),
		)
	default:
		return errors.New("unknown otp type: " + c.String("type"))
	}
}

func runSubcommandDelete(c *cli.Context) error {
	s, err := NewStore(c.Parent().String("database"))
	if err != nil {
		return err
	}
	defer s.Close()

	var failIds []string
	for _, id := range c.Args() {
		uid, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			failIds = append(failIds, id)
		}
		err = s.Delete(uid)
		if err != nil {
			failIds = append(failIds, id)
		}
	}
	if len(failIds) > 0 {
		fmt.Printf("%d record(s) delete failed: %s\n", len(failIds), strings.Join(failIds, ", "))
	}
	return nil
}

func runSubcommandExport(c *cli.Context) error {
	s, err := NewStore(c.Parent().String("database"))
	if err != nil {
		return err
	}
	defer s.Close()

	otps, err := s.List()
	if err != nil {
		return err
	}

	for _, otp := range otps {
		fmt.Println(otp.URI())
	}
	return nil
}

func main() {
	// parse flags
	if err := getApp().Run(os.Args); err != nil {
		panic(err)
	}
}

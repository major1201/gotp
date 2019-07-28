package main

import (
	"github.com/urfave/cli"
)

func getApp() *cli.App {
	app := cli.NewApp()
	app.Name = Name
	app.HelpName = app.Name
	app.Usage = app.Name
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "help, h",
			Usage: "show help",
		},
		cli.VersionFlag,
		cli.StringFlag{
			Name:   "database, d",
			Usage:  "specify database file",
			Value:  "gotp.db",
			EnvVar: "GOTP_DBFILE",
		},
		cli.BoolFlag{
			Name:  "id",
			Usage: "show item id or not",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add an otp object",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "uri",
					Usage: "(RECOMMENDED) add otp with URI",
				},
				cli.StringFlag{
					Name:  "type",
					Usage: "should be one of totp(default), hotp",
					Value: "totp",
				},
				cli.StringFlag{
					Name:  "issuer",
					Usage: "indicating the provider or service this account is associated with",
				},
				cli.StringFlag{
					Name:  "accountname",
					Usage: "identifying the provider or service managing that account",
				},
				cli.StringFlag{
					Name:  "secret",
					Usage: "an arbitrary key value encoded in Base32 according to RFC 3548",
				},
				cli.StringFlag{
					Name:  "algorithm",
					Usage: "should be one of sha1, sha256, sha512",
					Value: "sha1",
				},
				cli.UintFlag{
					Name:  "digits",
					Usage: "determines how long of a one-time passcode to display to the user",
					Value: 6,
				},
				cli.Uint64Flag{
					Name:  "counter",
					Usage: "REQUIRED if type is hotp: The counter parameter is required when provisioning a key for use with HOTP. It will set the initial counter value",
					Value: 0,
				},
				cli.Uint64Flag{
					Name:  "period",
					Usage: "OPTIONAL only if type is totp: The period parameter defines a period that a TOTP code will be valid for, in seconds",
					Value: 30,
				},
			},
			Action: runSubcommandAdd,
		},
		{
			Name:    "delete",
			Aliases: []string{"del", "d"},
			Usage:   "delete OTP item by id",
			Action:  runSubcommandDelete,
		},
		{
			Name:   "export",
			Usage:  "export all OTP object to URI",
			Action: runSubcommandExport,
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.Bool("help") {
			cli.ShowAppHelpAndExit(c, 0)
		}
		return runApp(c)
	}
	app.HideHelp = true
	return app
}

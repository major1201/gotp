# gotp

Golang OTP library and Google Authenticator CLI 

[![GoDoc](https://godoc.org/github.com/major1201/gotp?status.svg)](https://godoc.org/github.com/major1201/dohproxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/major1201/gotp)](https://goreportcard.com/report/github.com/major1201/dohproxy)

## Installation

Download from the latest [release](https://github.com/major1201/gotp/releases) page,

or install from source.

```bash
$ go get -u github.com/major1201/gotp/cmd/gotp
```

## Use as command line

Set DB file

```bash
# set DB file with envvar
export GOTP_DBFILE=/var/lib/gotp/default.db

# or you can specify db file path every time you execute gotp
gotp --database /var/lib/gotp/default.db
```

Add an OTP with URI provided

```bash
gotp add --uri otpauth://totp/ACME%20Co:john.doe@email.com?secret=HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ&issuer=ACME%20Co&algorithm=SHA1&digits=6&period=30
```

Add an OTP with detailed arguments

```bash
gotp add --issuer "ACME Co" --accountname "john.doe@email.com" --secret HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ
```

Display all otp objects and generate values

```bash
gotp

# with ID
gotp --id
```

Delete otp objects

```bash
gotp delete 11 13
```

Export all otp objects

```bash
gotp export
```

## Contributing

Just fork the repository and open a pull request with your changes.

## Licence

MIT

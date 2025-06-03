package main

import (
	"errors"
	"flag"

	"rmm23/src/l"
)

func main() {
	l.Name.Set("rmm23")
	l.CLI.Set()

	l.Z{l.M: "main", "daemon": l.Name.String()}.Debug()
	defer l.Z{l.M: "exit", "daemon": l.Name.String()}.Debug()

	var (
		err       error
		xmlConfig = new(xmlConf)
	)

	switch err = xmlConfig.load(); {
	case errors.Is(err, l.ENOCONF):
		flag.PrintDefaults()
		l.Z{l.E: err}.Critical()
	case err != nil:
		flag.PrintDefaults()
		l.Z{l.E: err}.Critical()
	}

	switch err = xmlConfig.LDAP.Fetch(); {
	case err != nil:
		l.Z{l.E: err}.Critical()
	}
}

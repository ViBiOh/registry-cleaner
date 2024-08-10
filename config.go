package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/flags"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
)

type configuration struct {
	logger *logger.Config

	url      *string
	username *string
	owner    *string
	password *string
	image    *string
	grep     *string
	last     *bool
	invert   *bool
	delete   *bool
	list     *bool
}

func newConfig() configuration {
	fs := flag.NewFlagSet("kitten", flag.ExitOnError)
	fs.Usage = flags.Usage(fs)

	config := configuration{
		logger: logger.Flags(fs, "logger"),

		url:      flags.New("URL", "Registry URL").DocPrefix("registry").String(fs, dockerHub, nil),
		username: flags.New("Username", "Registry username").DocPrefix("registry").String(fs, "", nil),
		owner:    flags.New("Owner", "For Docker Hub, fallback to username if not defined").DocPrefix("registry").String(fs, "", nil),
		password: flags.New("Password", "Registry password").DocPrefix("registry").String(fs, "", nil),
		image:    flags.New("Image", "Image name").DocPrefix("registry").String(fs, "", nil),
		grep:     flags.New("Grep", "Matching tags regexp, the capturing name tagBucket determine the bucket for getting the last").DocPrefix("cleaner").String(fs, "", nil),
		last:     flags.New("Last", "Keep only last tag found, in alphabetic order").DocPrefix("cleaner").Bool(fs, false, nil),
		invert:   flags.New("Invert", "Invert alphabetic order").DocPrefix("cleaner").Bool(fs, false, nil),
		delete:   flags.New("Delete", "Perform delete").DocPrefix("cleaner").Bool(fs, false, nil),
		list:     flags.New("List", "List repositories and doesn't do anything else").DocPrefix("cleaner").Bool(fs, false, nil),
	}

	_ = fs.Parse(os.Args[1:])

	return config
}

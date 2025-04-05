package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ViBiOh/flags"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
)

type configuration struct {
	image    *string
	url      *string
	username *string
	owner    *string
	password *string
	logger   *logger.Config
	last     *bool
	invert   *bool
	dryRun   *bool
	grep     *string
	source   *string
	target   *string
	command  string
}

func newConfig() configuration {
	fs := flag.NewFlagSet("registry-cleaner", flag.ExitOnError)
	fs.Usage = flags.Usage(fs)

	config := configuration{
		logger: logger.Flags(fs, "logger"),

		url:      flags.New("Url", "Registry URL").DocPrefix("registry").String(fs, dockerHub, nil),
		username: flags.New("Username", "Registry username").DocPrefix("registry").String(fs, "", nil),
		owner:    flags.New("Owner", "For Docker Hub, fallback to username if not defined").DocPrefix("registry").String(fs, "", nil),
		password: flags.New("Password", "Registry password").DocPrefix("registry").String(fs, "", nil),

		image: flags.New("Image", "Image name").DocPrefix("registry").String(fs, "", nil),

		last:   flags.New("Last", "Keep only last tag found, in alphabetic order").DocPrefix("delete").Bool(fs, false, nil),
		invert: flags.New("Invert", "Invert alphabetic order").DocPrefix("delete").Bool(fs, false, nil),
		dryRun: flags.New("DryRun", "Don't perform delete").DocPrefix("delete").Bool(fs, false, nil),
		grep:   flags.New("Grep", "Matching tags regexp, the capturing name tagBucket determine the bucket for getting the last").DocPrefix("delete").String(fs, "", nil),

		source: flags.New("Source", "Source tag").DocPrefix("promote").String(fs, "", nil),
		target: flags.New("Target", "Target tag").DocPrefix("promote").String(fs, "", nil),
	}

	_ = fs.Parse(os.Args[1:])

	args := fs.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Action (list, promote, delete) is required\n")
		os.Exit(1)
	}

	config.command = args[0]

	return config
}

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/ViBiOh/flags"
	"github.com/ViBiOh/httputils/v4/pkg/concurrent"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/registry-cleaner/pkg/hub"
	"github.com/ViBiOh/registry-cleaner/pkg/registry"
)

const (
	dockerHub = "https://registry-1.docker.io/"
)

// RegistryService definition
type RegistryService interface {
	Repositories(context.Context) ([]string, error)
	Tags(context.Context, string, func(string)) error
	Delete(context.Context, string, string) error
}

func main() {
	fs := flag.NewFlagSet("registry-cleaner", flag.ExitOnError)
	fs.Usage = flags.Usage(fs)

	loggerConfig := logger.Flags(fs, "logger")

	url := flags.New("URL", "Registry URL").DocPrefix("registry").String(fs, dockerHub, nil)
	username := flags.New("Username", "Registry username").DocPrefix("registry").String(fs, "", nil)
	owner := flags.New("Owner", "For Docker Hub, fallback to username if not defined").DocPrefix("registry").String(fs, "", nil)
	password := flags.New("Password", "Registry password").DocPrefix("registry").String(fs, "", nil)
	image := flags.New("Image", "Image name").DocPrefix("registry").String(fs, "", nil)
	grep := flags.New("Grep", "Matching tags regexp").DocPrefix("cleaner").String(fs, "", nil)
	last := flags.New("Last", "Keep only last tag found, in alphabetic order").DocPrefix("cleaner").Bool(fs, false, nil)
	invert := flags.New("Invert", "Invert alphabetic order").DocPrefix("cleaner").Bool(fs, false, nil)
	delete := flags.New("Delete", "Perform delete").DocPrefix("cleaner").Bool(fs, false, nil)
	list := flags.New("List", "List repositories and doesn't do anything else").DocPrefix("cleaner").Bool(fs, false, nil)

	logger.Fatal(fs.Parse(os.Args[1:]))

	logger.Global(logger.New(loggerConfig))
	defer logger.Close()

	ctx := context.Background()

	registryURL, imageName, matcher := checkParam(*url, *image, *grep, *list)

	var service RegistryService
	var err error

	if registryURL == dockerHub {
		service, err = hub.New(ctx, *username, *password, *owner)
	} else {
		service, err = registry.New(registryURL, *username, *password)
	}

	if err != nil {
		logger.Fatal(fmt.Errorf("create registry client: %w", err))
	}

	if *list {
		repositories, err := service.Repositories(ctx)
		if err != nil {
			logger.Fatal(fmt.Errorf("list repositories: %w", err))
		}

		fmt.Printf("%s", strings.Join(repositories, "\n"))
		return
	}

	var lastTag string
	var handled bool

	limiter := concurrent.NewLimited(runtime.NumCPU())

	err = service.Tags(ctx, imageName, func(tag string) {
		if !matcher.MatchString(tag) {
			return
		}

		if *last {
			if lastTag, handled = lastHandler(service, *invert, *delete, imageName, tag, lastTag); handled {
				return
			}
		}

		limiter.Go(func() {
			tagHandler(service, *delete, imageName, tag)
		})
	})

	if err != nil {
		logger.Fatal(fmt.Errorf("list tags: %w", err))
	}

	limiter.Wait()
}

func lastHandler(service RegistryService, invert, delete bool, image, tag, lastTag string) (string, bool) {
	if len(lastTag) == 0 {
		return tag, true
	}

	if (invert && tag < lastTag) || tag > lastTag {
		tagHandler(service, delete, image, lastTag)
		return tag, true
	}

	return lastTag, false
}

func tagHandler(service RegistryService, delete bool, image, tag string) {
	if !delete {
		logger.Warn("%s:%s is eligible to deletion", image, tag)
		return
	}

	if err := service.Delete(context.Background(), image, tag); err != nil {
		logger.Fatal(fmt.Errorf("delete `%s:%s`: %w", image, tag, err))
	}
	logger.Info("%s:%s deleted!", image, tag)
}

func checkParam(url, image, grep string, list bool) (string, string, *regexp.Regexp) {
	registryURL := strings.TrimSpace(url)
	if len(registryURL) == 0 {
		logger.Fatal(errors.New("url is required"))
	}

	if list {
		return registryURL, "", nil
	}

	imageName := strings.ToLower(strings.TrimSpace(image))
	if len(imageName) == 0 {
		logger.Fatal(errors.New("image is required"))
	}

	grepValue := strings.TrimSpace(grep)
	if len(grepValue) == 0 {
		logger.Fatal(errors.New("grep pattern is required"))
	}

	matcher, err := regexp.Compile(grepValue)
	if err != nil {
		logger.Fatal(fmt.Errorf("compile grep regexp: %w", err))
	}

	if err != nil {
		logger.Fatal(fmt.Errorf("create registry client: %w", err))
	}

	return registryURL, imageName, matcher
}

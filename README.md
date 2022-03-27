# registry-cleaner

[![Build](https://github.com/ViBiOh/registry-cleaner/workflows/Build/badge.svg)](https://github.com/ViBiOh/registry-cleaner/actions)
[![codecov](https://codecov.io/gh/ViBiOh/registry-cleaner/branch/main/graph/badge.svg)](https://codecov.io/gh/ViBiOh/registry-cleaner)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ViBiOh_registry-cleaner&metric=alert_status)](https://sonarcloud.io/dashboard?id=ViBiOh_registry-cleaner)

## Getting started

Golang binary is built with static link. You can download it directly from the [Github Release page](https://github.com/ViBiOh/registry-cleaner/releases) or build it by yourself by cloning this repo and running `make`.

A Docker image is available for `amd64`, `arm` and `arm64` platforms on Docker Hub: [vibioh/registry-cleaner](https://hub.docker.com/r/vibioh/registry-cleaner/tags).

You can configure app by passing CLI args or environment variables (cf. [Usage](#usage) section). CLI override environment variables.

## Usage

The application can be configured by passing CLI args described below or their equivalent as environment variable. CLI values take precedence over environments variables.

Be careful when using the CLI values, if someone list the processes on the system, they will appear in plain-text. Pass secrets by environment variables: it's less easily visible.

```bash
Usage of registry-cleaner:
  -delete
        [cleaner] Perform delete {REGISTRY_CLEANER_DELETE}
  -grep string
        [cleaner] Matching tags regexp {REGISTRY_CLEANER_GREP}
  -image string
        [registry] Image name {REGISTRY_CLEANER_IMAGE}
  -invert
        [cleaner] Invert alphabetic order {REGISTRY_CLEANER_INVERT}
  -last
        [cleaner] Keep only last tag found, in alphabetic order {REGISTRY_CLEANER_LAST}
  -loggerJson
        [logger] Log format as JSON {REGISTRY_CLEANER_LOGGER_JSON}
  -loggerLevel string
        [logger] Logger level {REGISTRY_CLEANER_LOGGER_LEVEL} (default "INFO")
  -loggerLevelKey string
        [logger] Key for level in JSON {REGISTRY_CLEANER_LOGGER_LEVEL_KEY} (default "level")
  -loggerMessageKey string
        [logger] Key for message in JSON {REGISTRY_CLEANER_LOGGER_MESSAGE_KEY} (default "message")
  -loggerTimeKey string
        [logger] Key for timestamp in JSON {REGISTRY_CLEANER_LOGGER_TIME_KEY} (default "time")
  -password string
        [registry] Registry password {REGISTRY_CLEANER_PASSWORD}
  -uRL string
        [registry] Registry URL {REGISTRY_CLEANER_URL} (default "https://registry-1.docker.io/")
  -username string
        [registry] Registry username {REGISTRY_CLEANER_USERNAME}
```

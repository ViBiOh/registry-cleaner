# registry-cleaner

[![Build](https://github.com/ViBiOh/registry-cleaner/workflows/Build/badge.svg)](https://github.com/ViBiOh/registry-cleaner/actions)

## Getting started

```bash
go install github.com/ViBiOh/registry-cleaner@latest

for repo in $(registry-cleaner -username "vibioh" -password "secret" -list); do
  registry-cleaner -username "vibioh" -password "secret" -image "${repo}" -grep ".*"
done

registry-cleaner -image "vibioh/fibr" -username "vibioh" -password "secret" -grep "[0-9]{12}" -last # keep only the last image that match regexp (which is a timestamp)
```

Golang binary is built with static link. You can download it directly from the [GitHub Release page](https://github.com/ViBiOh/registry-cleaner/releases), use the above command line or build it by yourself by cloning this repo and running `make`.

You can configure app by passing CLI args or environment variables (cf. [Usage](#usage) section). CLI override environment variables.

## Usage

The application can be configured by passing CLI args described below or their equivalent as environment variable. CLI values take precedence over environments variables.

Be careful when using the CLI values, if someone list the processes on the system, they will appear in plain-text. Pass secrets by environment variables: it's less easily visible.

```bash
Usage of registry-cleaner:
  --delete                    [cleaner] Perform delete ${REGISTRY_CLEANER_DELETE} (default false)
  --grep              string  [cleaner] Matching tags regexp ${REGISTRY_CLEANER_GREP}
  --image             string  [registry] Image name ${REGISTRY_CLEANER_IMAGE}
  --invert                    [cleaner] Invert alphabetic order ${REGISTRY_CLEANER_INVERT} (default false)
  --last                      [cleaner] Keep only last tag found, in alphabetic order ${REGISTRY_CLEANER_LAST} (default false)
  --list                      [cleaner] List repositories and doesn't do anything else ${REGISTRY_CLEANER_LIST} (default false)
  --loggerJson                [logger] Log format as JSON ${REGISTRY_CLEANER_LOGGER_JSON} (default false)
  --loggerLevel       string  [logger] Logger level ${REGISTRY_CLEANER_LOGGER_LEVEL} (default "INFO")
  --loggerLevelKey    string  [logger] Key for level in JSON ${REGISTRY_CLEANER_LOGGER_LEVEL_KEY} (default "level")
  --loggerMessageKey  string  [logger] Key for message in JSON ${REGISTRY_CLEANER_LOGGER_MESSAGE_KEY} (default "msg")
  --loggerTimeKey     string  [logger] Key for timestamp in JSON ${REGISTRY_CLEANER_LOGGER_TIME_KEY} (default "time")
  --owner             string  [registry] For Docker Hub, fallback to username if not defined ${REGISTRY_CLEANER_OWNER}
  --password          string  [registry] Registry password ${REGISTRY_CLEANER_PASSWORD}
  --uRL               string  [registry] Registry URL ${REGISTRY_CLEANER_URL} (default "https://registry-1.docker.io/")
  --username          string  [registry] Registry username ${REGISTRY_CLEANER_USERNAME}
```

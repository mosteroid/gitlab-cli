gitlab-cli
=========
[![Build Status](https://travis-ci.org/mosteroid/gitlab-cli.svg?branch=master)](https://travis-ci.org/mosteroid/gitlab-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/mosteroid/gitlab-cli)](https://goreportcard.com/report/github.com/mosteroid/gitlab-cli)

The command line interface for gitlab.


```
Usage:
  gitlab-cli [command]

Available Commands:
  help        Help about any command
  pipelines   Manage the pipelines
  projects    Manage the projects

Flags:
      --accessToken string   Set the user access token
      --baseUrl string       Set the gitlab base url
      --config string        Set the config file (default is $HOME/.gitlab-cli.yaml)
  -h, --help                 help for gitlab-cli
  -k, --insecure             Allow connections to SSL sites without certs

```

## Configuration File
Create in your home directory the file **.gitlab-cli.yaml**.

The following table lists the configurable parameters and theirs default values.

|             Parameter               |            Description                       |                    Default                |
|-------------------------------------|----------------------------------------------|-------------------------------------------|
| `gitlab.baseUrl`                    | The gitlab instance base URL                 | `nil`                                     |
| `gitlab.accessToken`                | The access token                             | `nil`                                     |
| `gitlab.insecure`                   | Allow connections to SSL sites without certs | `false`                                   |

For generating the **access token** follow the steps described [here](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html).

An example of a configuration file is the follow:

```yaml
  gitlab:
    baseUrl: https://gitlab.com
    accesstoken: XXX
```


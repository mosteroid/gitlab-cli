gitlabctl
=========
[![Build Status](https://travis-ci.org/mosteroid/gitlabctl.svg?branch=master)](https://travis-ci.org/mosteroid/gitlabctl)
[![Go Report Card](https://goreportcard.com/badge/github.com/mosteroid/gitlabctl)](https://goreportcard.com/report/github.com/mosteroid/gitlabctl)

`gitlabct` is a command line interface for GitLab.


```
Usage:
  gitlabctl [command]

Available Commands:
  config      Modify the configuration file
  help        Help about any command
  job         Manage jobs
  pipeline    Manage pipelines
  project     Manage projects

Flags:
      --accessToken string   Set the user access token
      --baseUrl string       Set the gitlab base url
      --config string        Set the config file (default is $HOME/.gitlabctl.yaml)
  -h, --help                 help for gitlabctl
  -k, --insecure             Allow connections to SSL sites without certs

Use "gitlabctl [command] --help" for more information about a command.
```

## Configuration File
Create in your home directory the file **.gitlabctl.yaml**.

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


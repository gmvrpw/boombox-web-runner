# Boombox - modular, open and easily customizable music player designed with respect for the ideas of clean architecture.

> Boombox Web Runner **is not a standalone piece of code** and it does not work by itself.
> This repository contains only runner features and does not know how to response on and store user requests.
> You can find a fully working examples in the [corresponding repository].

## Environment Variables
Environment variables are used only for passing secrets, and the path to the configuration file, but not the program configuration itself.

All secrets have a paired environment variable with the suffix `_FILE`,
which specifies the path to the file containing the secret.
The secret file is less specific than the environment variable
and will not be used if the secret's environment variable exists.

| Secret                               | Description                                |
|--------------------------------------|--------------------------------------------|
| `BOOMBOX_WEB_RUNNER_CONFIG_FILE`     | Config file path                           |

## Configuration
The parameters specified in the config file are the most specific and if they are present, the program will not use environment variables.
*- Yes, this means that you can specify secrets in the config file, but we strongly discourage you from doing so.*

| Field                           | Type     | Description                                                                    |
|---------------------------------|----------|--------------------------------------------------------------------------------|
| `entrypoints.http.port`         | `int`    | HTTP entrypoint port                                                           |
| `modules[n].name`               | `string` | Name of module                                                                 |
| `modules[n].test`               | `string` | Discord front-end token                                                        |
| `modules[n].auth`               | `object` | (Optional) Authorization cookie                                                |
| `modules[n].auth.name`          | `string` | Authorization cookie name                                                      |
| `modules[n].auth.value`         | `string` | Authorization cookie value                                                     |
| `modules[n].auth.domain`        | `string` | Authorization cookie domain                                                    |
| `modules[n].playback.selector`  | `string` | (Optional[^1]) CSS selector of element containing track timecode               |
| `modules[n].duration.selector`  | `string` | (Optional[^1]) CSS selector of element containing track duration               |
| `modules[n].remaining.selector` | `string` | (Optional[^1]) CSS selector of the element containing the remaining track time |
| `modules[n].play.selector`      | `string` | (Optional) CSS selector of the element to be clicked when playback starts      |

[^1]: One of the playback and duration or remaining sets must be configured. The remaining field is more specific.

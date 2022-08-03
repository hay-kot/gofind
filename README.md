<!-- PROJECT LOGO -->
<br />
<div align="center">
    <h1 align="center">GoFind</h1>
    <p align="center"><i>An fzf knockoff</i></p>
</div>

- [About](#about)
- [Installation](#installation)
- [Config](#config)
  - [Commands](#commands)
  - [Configuration Example](#configuration-example)
- [Help](#help)
- [Examples](#examples)

## About

GoFind is a small cli program for quickly finding and searching directories using the producer/consumer pattern. It uses the `filepath.Match` function to look for matching children in root directories, and if the child is found, it will stop return the parent directory and stop recursively looking into that parents children.

The primary use case for this library is to tie into other scripts to quickly and easily navigate (or other action) to a file path. For example you can use an alias to quickly search through all your git repositories and quickly search, find, and navigate in your terminal. See the [examples](#examples) for more details

https://user-images.githubusercontent.com/64056131/182485271-0c906802-c44e-4059-8079-37d6ea86e005.mp4

## Installation

You can install GoFind by running `go install github.com/hay-kot/gofind@latest` or download the latest release from [GitHub](https://github.com/hay-kot/gofind/releases)

After installing you must run `gofind setup` to initialize the configuration file and setup the cache

## Config

GoFind uses a json file in `~/.config/gofind.json` to store the configuration for the search entries and the default search. It also has a cache file that's used to store the results of a search so they aren't computed every time. You can set this path in the configuration file. If you ever need to re-cache the results, you can run `gofind cache` to recompute all caches, however they expire after 24 hours.


| Key             | Type        | Description                                                   |
| --------------- | ----------- | ------------------------------------------------------------- |
| `default`       | string      | The default command to run when calling `gofind find`         |
| `commands`      | object      | An object containing all of the registered commands available |
| `cache`         | string/path | Path to the cache directory                                   |
| `max_recursion` | number      | Max level of recursion from the root directory                |


### Commands

The commands key is an object where the `key` is the name of the command and the value is an object with the following keys:

| Key     | Type         | Description                                            |
| ------- | ------------ | ------------------------------------------------------ |
| `roots` | string array | The root directories to search inside of for the match |
| `match` | string       | The file match pattern used to find matches            |

In the example configuration there are two commands registered

- **docker:** will search the ~/docker directory of and present you with a list of all subdirectories that contain a docker-compose* file match. This is useful for quickly navigating to a docker-compose file directory to manage a container stack
- **repos:** will search the "~Code/OtherRepos" and "~/Code/Repos" directories for a match of the `.git` directory (or file). This is useful for quickly navigating to a git repository

### Configuration Example

```json
{
  "default": "repos",
  "commands": {
    "docker": {
      "roots": [
        "~/docker"
      ],
      "match": "docker-compose*"
    },
    "repos": {
      "roots": [
        "~/Code/OtherRepos",
        "~/Code/Repos"
      ],
      "match": ".git"
    }
  },
  "cache": "~/.cache/gofind/"
}
```

## Help

```shell
NAME:
   gofind - an interactive search for directories using the filepath.Match function

USAGE:
   gofind [global options] command [command options] [arguments...]

VERSION:
   0.1.3

COMMANDS:
   cache, c   cache all config entries
   find, f    run interactive finder for entry
   setup      first time setup
   config, c  add, remove, or list configuration entries
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

## Examples

**Change into Directory**

```shell
repos() {
    # Navigate to repos director and open target directory is specified
    if [ -z "$1" ]; then
        cd "`gofind find repos`"
        return
    fi
}
```

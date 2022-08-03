# Gofind

GoFind is a small cli program for quickly finding and searching directories using the producer/consumer pattern. It uses the `filepath.Match` function to look for matching children in root directories, and if the child is found, it will stop return the parent directory and stop recursively looking into that parents children.

**Why though?** The primary use case for this library is to tie into other scripts to quickly and easily navigate (or other action) to a file path. For example you can use an alias to quickly search through all your git repositories and quickly search, find, and navigate in your terminal. See the [examples](#examples) for more details

https://user-images.githubusercontent.com/64056131/182485271-0c906802-c44e-4059-8079-37d6ea86e005.mp4

## Config
GoFind uses a json file in `~/.config/gofind.json` to store the configuration for the search entries and the default search. It also uses this file to cache results so that the search is faster on subsequent runs. The config file example here has two jobs,

**repos:** which will recursively search the ~/Code directory for any directory that matches `.git`
**compose:** which will recursively search the ~/Docker directory for any directory that matches `docker-compose*`

I use these to either quickly find a repository I forgot the name of or where it exists, or quickly find a docker stack location and navigate to it.

```json
{
    "default": "repos",
    "commands": {
        "compose": {
            "root": ["~/Docker"],
            "match": "docker-compose*"
        },
        "repos": {
            "root": ["~/Code"],
            "match": ".git"
        }
    },
    "cache": {}
}
```

## Help

```shell
NAME:
   gofind - an interactive search for directories using the filepath.Match function

USAGE:
   gofind [global options] command [command options] [arguments...]

VERSION:
   0.1.0

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


**Open VSCode**

```shell
alias fcode="code \`gofind find repos\`"
```

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

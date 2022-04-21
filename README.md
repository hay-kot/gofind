# Gofind

GoFind is a tiny and poorly written tool that I use to quickly search through my file system to quickly navigate between repositories or any directory that has a specific child. It uses the `filepath.Match` function to look for matching children in root directories, and if the child is found, it will stop return the parent directory and stop recursively looking into that parents children. 

Once a match is found it's then piped to `fzf` to allow the user to select a directory and then it's spit out to the terminal which can then be used to change into that director or open it in vscode or whatever. **IMPORTANT** fzf must be installed and in the path for this to work, additionally the way I'm piping the output to fzf is probably not good practice and will likely not work on other machines YMMV.

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
            "root": "~/Docker",
            "match": "docker-compose*"
        },
        "repos": {
            "root": "~/Code",
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
   gofind [config-entry string] e.g. `gofind repos`

COMMANDS:
   cache, c  cache all config entries
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```
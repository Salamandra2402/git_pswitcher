# README #

git_pswitcher is a utility to easily keep and switch between your Git profiles. Only Linux is currently supported.

## Usage
```
USAGE:
   git_pswitcher [global options] command [command options] [arguments...]

VERSION:
   0.1.2

DESCRIPTION:
   git_pswitcher is an utility to keep git profiles and easily switch between them

COMMANDS:
     list, l    get all stored profiles
     remove, r  remove a profile from db by name. Specify the name with --name flag
     add, a     add a profile to the list. Use --name and --email to set values
     switch, s  switch to a profile. Just set the name with --name
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --name value   name of a git user
   --email value  email of a git user
   --help, -h     show help
   --version, -v  print the version
```

## Where is the database stored?

The database with your profiles can be found in ```/home/$USER/.git_pswitcher
It can be easily moved to another machine.

## TODO

* Add Windows support (only Linux is supported currently)
* Support SQLite databases to store git profiles (optional)
* Use git wrapper instead of using os.Exec calls



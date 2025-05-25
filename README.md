# toggle

**toggle** is a go tool that acts as a toggle for any program by using a lockfile.

you run it with a program like `toggle myscript.sh` - it creates a lockfile (by hash or name with `--name`), starts the program, and saves its PID to a lockfile

next time you run the same toggle, it finds the lockfile, reads the PID, and sends it a signal (default SIGTERM, or custom with `--signal`)

good for hotkeys that launch infinite loop scripts (like autoclickers) - press once to start, press again to terminate

```c
go install github.com/notwithering/toggle@latest
```

## licenses

this project uses the following dependancies with the licenses as noted:

- [github.com/alecthomas/kong](https://github.com/alecthomas/kong) - MIT License

each dependancy retains its respective license. for more details refer to their official documentation or source code

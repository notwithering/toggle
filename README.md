# toggle

**toggle** is a tiny go tool that acts as a toggle for any program using a lockfile.

you run it with a program like `toggle myscript.sh` - it creates a lockfile (by hash or name with `--name`), starts the program, and saves its PID.

next time you run the same toggle, it finds the lockfile, reads the PID, and sends it a signal (default SIGTERM, or custom with `--signal`).

handy for hotkeys that launch infinite loop scripts (like autoclickers) - press once to start, press again to kill.

```c
go install github.com/notwithering/toggle@latest
```

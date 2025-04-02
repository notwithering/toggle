# toggle

toggle is a go program that acts as a toggle by lockfile for any program

it works by accepting a program in the arguments and creaing a lockfile for it (by hash or by specified name). it starts the program and places the PID in the lockfile. when another instance of toggle tries to open the same program (program with same hash or same specified name) it will read the lockfile and send a signal to it specified by --signal or default SIGTERM

this is useful if you have a hotkey that starts an inifinite loop script (like an autoclicker) that needs to be killed with the same hotkey

in the hotkey you can specifiy it to run `toggle myscript.sh` or `toggle a.out`. the first time you press the hotkey it will start, the next time you press the hotkey it will terminate

```c
go install github.com/notwithering/toggle@latest
```

package main

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"syscall"

	"github.com/alecthomas/kingpin/v2"
)

var (
	nameFlag = kingpin.Flag("name", "Specify custom lockfile name.").Short('n').String()
	name     string

	exeArg = kingpin.Arg("exe", "File to execute.").Required().ExistingFile()
	exe    string

	argsArg = kingpin.Arg("args", "Arguments to pass to the file.").Strings()
	args    []string

	signalFlag = kingpin.Flag("signal", "Specify a signal number to send when stopping the process.").Short('s').Default("15").Int()
	signalInt  int
)

func main() {
	kingpin.Parse()
	name = *nameFlag
	exe = *exeArg
	args = *argsArg
	signalInt = *signalFlag

	exe = abs(exe)
	cmd := makeCommand(exe, args)
	name = getName(exe, name)
	lockPath := getLockPath(name)
	lock, lockExists := openLock(lockPath)
	if lockExists {
		defer lock.Close()
		pid := getPID(lock)
		sendSignal(pid, signalInt, lock)
	} else {
		lock = makeLock(lockPath)
		startCommand(cmd, lock)
		go forwardSignals(cmd)
		writePID(lock, cmd)
		cmd.Wait()
		os.Remove(lockPath)
	}
}

func abs(file string) string {
	file, err := filepath.Abs(exe)
	if err != nil {
		kingpin.Fatalf("error finding absolute path: %v", err)
	}
	return file
}

func makeCommand(exe string, args []string) *exec.Cmd {
	var cmd *exec.Cmd

	if len(args) > 0 {
		cmd = exec.Command(exe, args...)
	} else {
		cmd = exec.Command(exe)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd
}

func getName(exe, name string) string {
	if name != "" {
		return name
	}

	hash := md5.New()

	file, err := os.Open(exe)
	if err != nil {
		kingpin.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	if _, err := io.Copy(hash, file); err != nil {
		kingpin.Fatalf("error hashing file: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}

func getLockPath(name string) string {
	return filepath.Join(os.TempDir(), "toggle", name+".lock")
}

func openLock(lockPath string) (*os.File, bool) {
	if err := os.MkdirAll(filepath.Join(os.TempDir(), "toggle"), 0o755); err != nil {
		kingpin.Fatalf("error creating directory: %v", err)
	}

	file, err := os.OpenFile(lockPath, os.O_RDONLY, 0o644)
	var noexist bool
	if err != nil {
		if err.(*os.PathError).Err.Error() == "no such file or directory" {
			noexist = true
		} else {
			kingpin.Fatalf("error opening lock file: %v", err)
		}
	}

	return file, !noexist
}

func getPID(lock *os.File) int {
	b, err := io.ReadAll(lock)
	if err != nil {
		kingpin.Fatalf("error reading lock file: %v", err)
	}

	pid, err := strconv.Atoi(string(b))
	if err != nil {
		if !regexp.MustCompile("\\d+").Match(b) {
			os.Remove(lock.Name())
			kingpin.Fatalf("invalid PID format; removed lock file")
		}
		kingpin.Fatalf("error parsing PID: %v", err)
	}

	return pid
}

func sendSignal(pid, signal int, lock *os.File) {
	process, err := os.FindProcess(pid)
	if err != nil {
		os.Remove(lock.Name())
		kingpin.Fatalf("error finding process: %v; removed lock file", err)
	}

	if err := process.Signal(syscall.Signal(signalInt)); err != nil && !errors.Is(err, os.ErrProcessDone) {
		kingpin.Fatalf("error sending signal: %v", err)
	}
}

func makeLock(lockPath string) *os.File {
	file, err := os.Create(lockPath)
	if err != nil {
		kingpin.Fatalf("error creating lock file: %v", err)
	}
	return file
}

func startCommand(cmd *exec.Cmd, lock *os.File) {
	if err := cmd.Start(); err != nil {
		os.Remove(lock.Name())
		kingpin.Fatalf("error starting command: %v", err)
	}
}

func forwardSignals(cmd *exec.Cmd) {
	c := make(chan os.Signal, 1)
	signal.Notify(c)

	for {
		sig := <-c
		if err := cmd.Process.Signal(sig); err != nil {
			if errors.Is(err, os.ErrProcessDone) {
				return
			}
			kingpin.Fatalf("error forwarding signal: %v", err)
		}
	}
}

func writePID(lock *os.File, cmd *exec.Cmd) {
	if _, err := fmt.Fprint(lock, cmd.Process.Pid); err != nil {
		os.Remove(lock.Name())
		kingpin.Fatalf("error writing to lockfile: %v", err)
	}
}

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
	"strconv"
	"syscall"

	"github.com/alecthomas/kingpin/v2"
)

var (
	nameFlag = kingpin.Flag("name", "Specify custom lockfile name.").Short('n').String()
	name     string

	// execFlag = kingpin.Flag("exec", "Execute a command inline.").Short('e').Strings()
	// execArgs []string

	fileArg  = kingpin.Arg("file", "File to execute.").Required().ExistingFile()
	fileName string

	argsArg = kingpin.Arg("args", "Arguments to pass to the file.").Strings()
	args    []string

	signalFlag = kingpin.Flag("signal", "Specify a signal number to send when stopping the process.").Short('s').Default("15").Int()
	signalInt  int
)

func main() {
	kingpin.Parse()
	name = *nameFlag
	// execArgs = *execFlag
	fileName = *fileArg
	args = *argsArg
	signalInt = *signalFlag

	var lockfileName string = name
	var cmd *exec.Cmd

	hash := md5.New()

	// if len(execArgs) > 0 {
	// 	var command string
	// 	if len(execArgs) > 1 {
	// 		cmd = exec.Command(execArgs[0], execArgs[1:]...)
	// 		command = execArgs[0] + " " + strings.Join(execArgs[1:], " ")
	// 	} else {
	// 		cmd = exec.Command(execArgs[0])
	// 		command = execArgs[0]
	// 	}

	// 	if lockfileName == "" {
	// 		hash.Write([]byte(command))
	// 		lockfileName = base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
	// 	}
	// } else {

	filePath, err := filepath.Abs(fileName)
	if err != nil {
		exitWithError("error getting absolute path: %v", err)
	}

	if len(args) > 0 {
		cmd = exec.Command(filePath, args...)
	} else {
		cmd = exec.Command(filePath)
	}

	if lockfileName == "" {
		file, err := os.Open(filePath)
		if err != nil {
			exitWithError("error opening file: %v", err)
		}
		defer file.Close()

		if _, err := io.Copy(hash, file); err != nil {
			exitWithError("error hashing file: %v", err)
		}
		lockfileName = base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
	}
	// }

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := os.MkdirAll(filepath.Join(os.TempDir(), "toggle"), 0o755); err != nil {
		exitWithError("error creating directory: %v", err)
	}

	lockfilePath := filepath.Join(os.TempDir(), "toggle", lockfileName+".lock")
	defer os.Remove(lockfilePath)

	file, err := os.OpenFile(lockfilePath, os.O_RDONLY, 0o644)
	var noexist bool
	if err != nil {
		if err.(*os.PathError).Err.Error() == "no such file or directory" {
			noexist = true
		} else {
			exitWithError("error opening lock file: %v", err)
		}
	}

	defer file.Close()

	if noexist {
		file, err = os.Create(lockfilePath)
		if err != nil {
			exitWithError("error creating lock file: %v", err)
		}
		defer file.Close()

		err = cmd.Start()
		if err != nil {
			exitWithError("error starting command: %v", err)
		}
		fmt.Fprint(file, cmd.Process.Pid)

		c := make(chan os.Signal, 1)
		signal.Notify(c)

		go func() {
			for {
				sig := <-c
				if err := cmd.Process.Signal(sig); err != nil {
					if errors.Is(err, os.ErrProcessDone) {
						return
					}
					exitWithError("error forwarding signal: %v", err)
				}
			}
		}()

		cmd.Wait()
	} else {
		b, err := io.ReadAll(file)
		if err != nil {
			exitWithError("error reading lock file: %v", err)
		}

		pid, err := strconv.Atoi(string(b))
		if err != nil {
			exitWithError("error parsing PID: %v", err)
		}

		process, err := os.FindProcess(pid)
		if err != nil {
			exitWithError("error finding process: %v", err)
		}

		if err := process.Signal(syscall.Signal(signalInt)); err != nil && !errors.Is(err, os.ErrProcessDone) {
			exitWithError("error sending signal: %v", err)
		}
	}
}

func exitWithError(format string, args ...interface{}) {
	kingpin.Errorf(format, args...)
	os.Exit(1)
}

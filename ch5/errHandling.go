package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
)

//MyError ...
type MyError struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

func wrapError(err error, messagef string, msgArgs ...interface{}) MyError {
	return MyError{
		Inner:      err, //<1>
		Message:    fmt.Sprintf(messagef, msgArgs...),
		StackTrace: string(debug.Stack()),        // <2>
		Misc:       make(map[string]interface{}), // <3>
	}
}

func (err MyError) Error() string {
	return err.Message
}

// LowLevelErr ...
type LowLevelErr struct {
	error
}

func isGloballyExec(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, LowLevelErr{(wrapError(err, err.Error()))} // <1>
	}
	return info.Mode().Perm()&0100 == 0100, nil
}

//IntermediateErr ...
type IntermediateErr struct {
	error
}

func runJob(id string) error {
	const jobBinPath = "/bad/job/binary"
	isExecutable, err := isGloballyExec(jobBinPath)
	if err != nil {
		return IntermediateErr{wrapError(
			err,
			"cannot run job %q: requisite binaries not availble",
			id,
		)} // <1>
	} else if isExecutable == false {
		return wrapError(nil, "cannot run job %q: requisite binaries are not executable", id)
	}

	return exec.Command(jobBinPath, "--id="+id).Run() // <1>
}

func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v]: ", key))
	log.Printf("%#v", err) // <3>
	fmt.Printf("[%v] %v", key, message)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	err := runJob("1")
	if err != nil {
		msg := "There was an unexpected issue; please report this as a bug."
		if _, ok := err.(IntermediateErr); ok { // <1>
			msg = err.Error()
		}
		handleError(1, err, msg) // <2>
	}
}

package errs

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/cmarkh/go-mail"
)

// SetupLog creates a log file and sets up program for logging with my desired parameters
func SetupLog() (file *os.File, logsPath string, err error) {
	t := time.Now().Format("2006-01-02_15-04-05")

	logsPath, err = os.Executable()
	if err != nil {
		return
	}
	logsPath = filepath.Dir(logsPath)

	name, err := ProgramName()
	if err != nil {
		return
	}

	logsPath = filepath.Join(logsPath, ".logs", name)
	err = os.MkdirAll(logsPath, 0700)
	if err != nil {
		return
	}

	logsPath = filepath.Join(logsPath, name+"_"+t+"_log.txt")

	//setup logging file
	//log.SetFlags(log.Lshortfile | log.LstdFlags)
	file, err = os.OpenFile(logsPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600) //O_WRONLY is open for writing only
	if err != nil {
		return
	}

	log.SetOutput(io.MultiWriter(file, os.Stdout)) //set output of logs to f and to the console
	return
}

// SetLogFlags sets my default log flags
func SetLogFlags() {
	log.SetFlags(log.Llongfile | log.LstdFlags) //Turns on line numbers for errors
}

func Log(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1) //1 = 1 level up the stack
	log.Printf(":\n\t[%s][%d]:\n\t%v", file, line, v)
}

func WrapErr(err error, addInfo ...string) error {
	_, file, line, _ := runtime.Caller(1) //1 = 1 level up the stack
	if len(addInfo) > 0 {
		return fmt.Errorf("[%s][%d]:\n\t%s\n\t%w", file, line, strings.Join(addInfo, "\n\t"), err)
	}
	return fmt.Errorf("[%s][%d]:\n\t%w", file, line, err)
}

// ProgramName returns the name of the running program
func ProgramName() (string, error) {
	name, err := os.Executable()
	if err != nil {
		log.Print(err)
		return "", err
	}
	return filepath.Base(name), nil
}

// Email sends error in my standard format
func Email(err error, account mail.Account, to ...string) error {
	name, e := ProgramName()
	if e != nil {
		name = fmt.Sprint("couldn't get program name: ", e)
	}

	e = account.Send("Program Error - "+name, fmt.Sprintln(err), to...)
	if err != nil {
		log.Println(err)
	}

	return e
}

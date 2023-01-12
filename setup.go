/*
This is a way to start the main function
Sets up logging
Resolves panics in main cleanly
*/

package errs

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

// DefaultSetup serves as a sort of default main function with logging, catch and reocver, and email errors
// Note: err won't be emailed if the send account was not previously specified
func DefaultSetup(fn func() error) {
	file, path, err := SetupLog()
	if err != nil {
		Email(err)
		return
	}

	//NOTE: defers execute in reverse line order (ie last first)
	defer DeleteEmptyLogFile(path) //no need to store empty logs
	defer file.Close()             //close log file
	defer CatchAndRecover()        //Handle any panics in main to exit cleanly

	//begin will run actual program functions and catch any errors
	if err = fn(); err != nil {
		Email(err)
		return
	}
}

func DeleteEmptyLogFile(path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	b = bytes.TrimSpace(b)

	if len(b) == 0 {
		err := os.Remove(path)
		if err != nil {
			Email(err)
			return
		}
	}
}

// PauseExit prevents the terminal from closing so user can read the error messages
func PauseExit(err error) {
	if err == nil {
		return
	}

	var exit string
	for {
		fmt.Println()
		fmt.Println("Warning - An error has occurred. Type 'exit' to end program:")
		fmt.Scanln(&exit)
		if strings.ToLower(exit) == "exit" {
			break
		}
	}
}

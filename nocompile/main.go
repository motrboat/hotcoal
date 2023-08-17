package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

var command = []string{"go", "test", "-gcflags=-e", "nocompile/nocompile.go"}

const expectedExitCode = 1

const expected = "# command-line-arguments\nnocompile/nocompile.go:7:22: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to hotcoal.Wrap\nnocompile/nocompile.go:9:19: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to hotcoal.W\nnocompile/nocompile.go:23:19: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to append\nnocompile/nocompile.go:29:22: cannot use []string{} (value of type []string) as []hotcoal.hotcoalString value in argument to hotcoal.Join\nnocompile/nocompile.go:31:25: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to hotcoal.Join\nnocompile/nocompile.go:35:19: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to y.Replace\nnocompile/nocompile.go:37:22: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to y.Replace\nnocompile/nocompile.go:41:22: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to y.ReplaceAll\nnocompile/nocompile.go:43:25: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to y.ReplaceAll\nnocompile/nocompile.go:47:17: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to b.Write\nnocompile/nocompile.go:53:19: cannot use b.String() (value of type string) as hotcoal.hotcoalString value in argument to hotcoal.W\nnocompile/nocompile.go:55:27: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to hotcoal.Allowlist\nnocompile/nocompile.go:59:30: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to hotcoal.Allowlist\nnocompile/nocompile.go:63:33: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to hotcoal.Allowlist\nFAIL\n"

func main() {
	fmt.Println("Running nocompile test")

	outputBytes, err := exec.
		Command(command[0], command[1:]...).
		CombinedOutput()

	output := string(outputBytes)

	if err != nil {
		var unwrappedError *exec.ExitError

		if !errors.As(err, &unwrappedError) {
			log.Fatalf("Cannot execute command: %#v, cannot unwrap error: %#v", command, err)
		}

		if unwrappedError.ExitCode() != expectedExitCode {
			log.Fatalf("Cannot execute command: %#v, received error: %#v", command, output)
		}
	}

	if expected != output {
		logFatalAndExit(expected, output)
	}

	fmt.Println("PASS\n")
}

func logFatalAndExit(expected, output string) {
	const (
		downArrows = "vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv"
		upArrows   = "^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^"
	)

	log.Fatalf(
		strings.Join(
			[]string{
				"NOCOMPILE TEST FAILED",
				"",
				"EXPECTED:",
				downArrows,
				"%s",
				upArrows,
				"",
				"GOT:",
				downArrows,
				"%s",
				upArrows,
				"",
				"REPR FOR SNAPSHOT TESTING:",
				downArrows,
				"%#v",
				upArrows,
				"\n",
			},
			"\n",
		),
		expected,
		output,
		output,
	)
}

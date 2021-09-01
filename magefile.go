// +build mage

package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var (
	Default             = Build
	GolangciLintVersion = "v1.41.1"
	GotestsumVersion    = "1.6.4"
	LDFLAGS             = "-ldflags='-s -w'"
)

func Clean() error {
	fmt.Println("CLEAN")
	return sh.RunV("rm", "-rf", "bin/*")
}

//Build Compile and lint the cli
func Build() error {
	files, err := ioutil.ReadDir("cmd")
	if err != nil {
		return err
	}

	clean := false
	prep := func() {
		mg.Deps(Format)
		mg.Deps(Lint)
		clean = true
	}

	for _, f := range files {
		notExists, err := target.Dir("bin/"+f.Name(), "cmd/"+f.Name(), "pkg")
		if err != nil {
			return err
		}

		if !notExists {
			continue
		}

		if !clean {
			prep()
		}

		mg.Deps(mg.F(Compile, f.Name()))
	}

	return sh.RunV("go", "build", "./...")
}

func Compile(bin string) error {
	fmt.Println(fmt.Sprintf("GO BUILD %s", strings.ToUpper(bin)))
	cmd := fmt.Sprintf("go build %s -o bin/%s cmd/%s/*.go", LDFLAGS, bin, bin)
	return sh.RunV("sh", "-c", cmd)
}

//Lint Nit the hell outta this code
func Lint() error {
	mg.Deps(golangciLint)

	fmt.Println("GO LINT")
	return sh.RunV("bin/golangci-lint", "run", "--color", "always")
}

type Test mg.Namespace

//Run Runs tests
func (Test) Run() error {
	mg.Deps(gotestsum)

	fmt.Println("GO TEST")
	return sh.RunV("bin/gotestsum")
}

//Watch Runs tests continously on file change
func (Test) Watch() error {
	mg.Deps(gotestsum)

	fmt.Println("GO TEST WATCH")
	return sh.RunV("bin/gotestsum", "--watch")
}

func Format() error {
	fmt.Println("GO FMT")
	return sh.RunV("go", "fmt", "./...")
}

func golangciLint() error {
	notExists, err := target.Path("bin/golangci-lint")
	if err != nil || !notExists {
		return err
	}

	fmt.Println("DEP GOLANGCI-LINT")
	cmd := fmt.Sprintf("curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b bin %s", GolangciLintVersion)
	return sh.RunV("sh", "-c", cmd)
}

func gotestsum() error {
	notExists, err := target.Path("bin/gotestsum")
	if err != nil || !notExists {
		return err
	}

	fmt.Println("DEP GOTESTSUM")
	cmd := fmt.Sprintf("curl -sSfL https://github.com/gotestyourself/gotestsum/releases/download/v%s/gotestsum_%s_linux_amd64.tar.gz | tar zx -C bin/ gotestsum", GotestsumVersion, GotestsumVersion)
	return sh.RunV("sh", "-c", cmd)
}

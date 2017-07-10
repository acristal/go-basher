// +build linux darwin
// +build !arm

package basher

import (
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
)

// Application sets up a common entrypoint for a Bash application that
// uses exported Go functions. It uses the DEBUG environment variable
// to set debug on the Context, and SHELL for the Bash binary if it
// includes the string "bash". You can pass a loader function to use
// for the sourced files, and a boolean for whether or not the
// environment should be copied into the Context process.
func Application(
	funcs map[string]func([]string),
	scripts []string,
	loader func(string) ([]byte, error),
	copyEnv bool) {

	bashDir, err := homedir.Expand("~/.basher")
	if err != nil {
		log.Fatal(err, "1")
	}

	bashPath := bashDir + "/bash"
	if _, err := os.Stat(bashPath); os.IsNotExist(err) {
		err = RestoreAsset(bashDir, "bash")
		if err != nil {
			log.Fatal(err, "1")
		}
	}
	bash, err := NewContext(bashPath, os.Getenv("DEBUG") != "")
	if err != nil {
		log.Fatal(err)
	}
	for name, fn := range funcs {
		bash.ExportFunc(name, fn)
	}
	if bash.HandleFuncs(os.Args) {
		os.Exit(0)
	}

	for _, script := range scripts {
		bash.Source(script, loader)
	}
	if copyEnv {
		bash.CopyEnv()
	}
	status, err := bash.Run("main", os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(status)
}

package version

import (
	"fmt"
	"os"
)

// Version of the package
var Version = "0.0.1"

// PrintVersion displays version in stdout
func PrintVersion() {
	fmt.Fprintln(os.Stdout, os.Args[0], Version)
}

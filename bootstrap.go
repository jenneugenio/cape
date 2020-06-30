// +build ignore

// This file contains the bootstrap enabling a user to bootstrap their local
// environment without needing to explicitly download and install Mage.
//
// For more details see here: https://magefile.org/zeroinstall/

package main

import (
	"os"

	"github.com/magefile/mage/mage"
)

func main() { os.Exit(mage.Main()) }

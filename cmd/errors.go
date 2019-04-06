package cmd

import (
	"fmt"
	"os"
)

func exitOnError(e error) {
	if e != nil {
		fmt.Println(e)
		defer os.Exit(1)
	}
}

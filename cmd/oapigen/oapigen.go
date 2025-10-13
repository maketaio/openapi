package main

import (
	"fmt"
	"os"

	"github.com/maketaio/api/internal/cli"
)

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		fmt.Printf("failed to execute codegen: %v", err)
		os.Exit(1)
	}
}

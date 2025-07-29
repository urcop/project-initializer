package main

import (
	"fmt"
	"os"

	"github.com/urcop/project-initializer/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка выполнения команды: %v\n", err)
		os.Exit(1)
	}
}

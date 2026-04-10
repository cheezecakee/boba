package main

import (
	"fmt"
	"os"

	"github.com/cheezecakee/boba"
	"github.com/cheezecakee/boba/screens"
)

func main() {
	if err := boba.NewApp(screens.NewMenu).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

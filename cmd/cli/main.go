package main

import (
	"fmt"
	"os"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// commands hold clis command and hook function
var commands = map[string]func() error{
	"migrate": migration,
	"seed":    seed,
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("command needed [migrate, seed]")
		os.Exit(0)
	}

	if err := commands[os.Args[1]](); err != nil {
		fmt.Println(err)
	}

}

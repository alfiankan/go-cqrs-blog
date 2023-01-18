package main

import (
	"fmt"
	"os"

	"github.com/alfiankan/go-cqrs-blog/common"
)

// commands hold clis command and hook function
var commands = map[string]func(wd string) error{
	"migrate": common.Migration,
	"seed":    common.Seed,
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("command needed [migrate, seed]")
		os.Exit(0)
	}

	if err := commands[os.Args[1]]("."); err != nil {
		fmt.Println(err)
	}

}

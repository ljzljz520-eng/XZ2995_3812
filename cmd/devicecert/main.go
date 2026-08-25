package main

import (
	"fmt"
	"os"

	"devicecert/cli"
	"devicecert/safelog"
	"devicecert/service"
	"devicecert/store"
)

func main() {
	path := "devicecert.db"
	db, err := store.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()
	app := cli.New(service.NewManager(db, safelog.New(os.Stdout)))
	output, err := app.Execute(os.Args[1:])
	if output != "" {
		fmt.Println(output)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

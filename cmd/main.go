package main

import (
	"log"
	"os"

	"github.com/buurzx/in-mem-kvdb/cmd/kvcli"
	"github.com/buurzx/in-mem-kvdb/cmd/server"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "key value database",
		Usage: "A simple in-memory KV database",
		Commands: []*cli.Command{
			kvcli.BuildCmd(),
			server.BuildCmd(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

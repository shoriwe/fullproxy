package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/shoriwe/fullproxy/v4/compose"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

const JsonFlag = "json"

var app = &cli.App{
	Commands: []*cli.Command{
		{
			Name:  "compose",
			Usage: "start fullproxy services based on the defined compose contract",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  JsonFlag,
					Usage: "decode compose as JSON",
				},
			},
			Action: func(ctx *cli.Context) error {
				composeFile := "fullproxy-compose.yaml"
				if ctx.Args().Len() == 1 {
					composeFile = ctx.Args().First()
				}
				composeContents, err := os.ReadFile(composeFile)
				if err != nil {
					return err
				}
				var c compose.Compose
				if ctx.Bool(JsonFlag) {
					err = json.NewDecoder(bytes.NewReader(composeContents)).Decode(&c)
					if err != nil {
						return err
					}
				} else {
					err = yaml.NewDecoder(bytes.NewReader(composeContents)).Decode(&c)
					if err != nil {
						return err
					}
				}
				err = c.Start()
				return err
			},
		},
	},
}

func init() {

}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

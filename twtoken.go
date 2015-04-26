package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = os.Args[0]
	app.Usage = "Get twitter access token from command line."
	app.Author = "pocke"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "ck, consumer-key",
			Value:  "",
			Usage:  "Consumer key of Twitter app",
			EnvVar: "CONSUMER_KEY",
		},
		cli.StringFlag{
			Name:   "cs, consumer-secret",
			Value:  "",
			Usage:  "Consumer secret of Twitter app",
			EnvVar: "CONSUMER_SECRET",
		},
	}

	app.Action = func(c *cli.Context) {
		ck := c.String("ck")
		cs := c.String("cs")
		if ck == "" || cs == "" {
			fmt.Fprintln(os.Stderr, "Consumer Key and Consumer Secret are required")
			os.Exit(1)
		}

		token := NewToken(ck, cs)
		url := token.URL()
		fmt.Printf("Access > %s\n", url)
		fmt.Println()

		aToken := token.AccessToken()
		fmt.Printf("Access Token:        %s\n", aToken.Token)
		fmt.Printf("Access Token Secret: %s\n", aToken.Secret)
	}

	app.Run(os.Args)
}

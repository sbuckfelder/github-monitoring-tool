/*
   Copyright awslabs Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

const usage = `GitHub Monitoring Tool`

func main() {
	app := App()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "github-monitoring-tool (ghmt): %s\n", err)
		os.Exit(1)
	}
}

func App() *cli.App {
	app := cli.NewApp()
	app.Name = "ghmt"
	app.Version = "0.0.0"
	app.Usage = usage
	app.Description = `
Simple CLI tool to help monitor and manage projects hosted on GitHub.`
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "testing,t",
			Usage: "checking passing in values",
			Value: "HACKY TEST STRING",
		},
	}
	app.Commands = []cli.Command{
		eventsCommand,
	}
	app.Action = func(context *cli.Context) error {
		testingString := context.GlobalString("testing")
		fmt.Printf("Testing String: %s\n", testingString)
		return nil
	}
	return app
}

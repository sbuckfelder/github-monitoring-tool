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
	app.Commands = []cli.Command{
		eventsCommand,
		prCommand,
	}
	return app
}

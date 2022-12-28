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
	"github.com/urfave/cli"
	"github.com/sbuckfelder/github-monitoring-tool/proxy"
)

var eventsCommand = cli.Command{
	Name:  "events",
	Usage: "list events for a repo",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "repo",
			Usage: "github repo to list events for",
		},
	},
	Action: func(ctx *cli.Context) error {
		_, err := proxy.NewProxy()
		if err != nil {
			panic("Failed to create GitHub client")
		}
		return nil
	},
}

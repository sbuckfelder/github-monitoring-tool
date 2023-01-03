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
	"github.com/sbuckfelder/github-monitoring-tool/proxy"
	"github.com/urfave/cli"
)

var prCommand = cli.Command{
	Name:  "pr",
	Usage: "List open pull requests events for a github org/repo.  Meant to identify old/stale ps's for triage.",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     ORG_NAME,
			Usage:    "github org repo belongs to",
			Required: true,
		},
		cli.StringFlag{
			Name:     REPO_NAME,
			Usage:    "github repo to list events for, cannot be used with repoall flag",
			Required: true,
		},
	},
	Action: func(ctx *cli.Context) error {
		orgInput := ctx.String(ORG_NAME)
		repoInput := ctx.String(REPO_NAME)

		ghProxy, err := proxy.NewProxy()
		if err != nil {
			panic("Failed to create GitHub client")
		}

		ghProxy.GetPullRequests(orgInput, repoInput)

		return nil
	},
}

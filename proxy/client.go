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

package proxy

import (
	"context"
	"io/ioutil"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

var (
	tokenFile = ".ghmt"
)

type GithubProxy struct {
	client github.Client
}

func NewProxy() (GithubProxy, error) {
	ctx := context.Background()
	token, err := getToken()	
	if err != nil {
		panicMsg := fmt.Sprintf("Failed to get token from %s: %v", getTokenFileName(), err)
		panic(panicMsg)
	}	
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token},
	)
	tcpClient := oauth2.NewClient(ctx, tokenSource)
	client := github.NewClient(tcpClient)
	return GithubProxy{
		client: *client}, nil
}

func getTokenFileName() string {
	homeDir, _  := os.UserHomeDir()
	return homeDir + "/" + tokenFile
}

func getToken() (string, error) {
	token, err := ioutil.ReadFile(getTokenFileName())
	if err != nil {
		return "", err
	}
	cleanToken := strings.Replace(string(token), "\n", "", -1)
	return cleanToken, nil
}

// Copyright 2025 Baptiste Gouhoury <baptiste.gouhoury@scalend.fr>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"io"
	"net/http"
	"os"

	"github.com/baptistegh/go-lakekeeper/cmd/lkctl/commands"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(commands.NewCommand(), "./docs/user-guide/commands")
	if err != nil {
		log.Fatal(err)
	}

	// get lakekeeper assets
	resp, err := http.Get("https://raw.githubusercontent.com/lakekeeper/lakekeeper/refs/heads/main/site/docs/assets/bear.svg")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	file, err := os.OpenFile("docs/assets/bear.svg", os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}

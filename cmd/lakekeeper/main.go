package main

import (
	"fmt"

	"github.com/baptistegh/go-lakekeeper/pkg/version"
)

func main() {
	fmt.Printf("version=%s, commit=%s, date=%s\n", version.Version, version.Commit, version.Date)
}

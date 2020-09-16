package main

import (
	"os"

	"github.com/rancher/wrangler/pkg/cleanup"
	"github.com/xshrim/gol"
)

func main() {
	if err := cleanup.Cleanup("./pkg/apis"); err != nil {
		gol.Fatal(err)
	}
	if err := os.RemoveAll("./pkg/generated"); err != nil {
		gol.Fatal(err)
	}
}

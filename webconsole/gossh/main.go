package main

import (
	"fmt"
	"os"

	"./utils"
)

func main() {
	cli := utils.New("127.0.0.1", "xshrim", " ", 22)
	output, err := cli.Run("free -h")
	fmt.Printf("%v\n%v", output, err)

	cli.RunTerminal("vim test.sh", os.Stdout, os.Stdin)
}

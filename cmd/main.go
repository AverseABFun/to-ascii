package main

import (
	"fmt"
	"os"

	to_ascii "github.com/averseabfun/to-ascii"
)

func main() {
	var file, err = os.Open("probably-ai.jpg")
	if err != nil {
		panic(err)
	}
	fmt.Println(to_ascii.Convert(file, to_ascii.ConvertConfig{OutputWidth: 197, OutputHeight: 28}))
}

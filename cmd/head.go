package cmd

import (
	"fmt"
)

func HandelHead(args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Error: Please specify a file to add.")
		return
	}
	number := args[1]

	fmt.Printf("going back %s commits \n", number)
}

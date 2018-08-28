package main

import (
	"bufio"
	"fmt"
	"github.com/stdi0/calc"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	in := strings.NewReader(text)
	_, err := calc.Calc(in, os.Stdout)

	if err != nil {
		fmt.Println(err)
		return
	}
}

package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	"github.com/stdi0/calc"
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
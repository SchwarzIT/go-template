package main

import (
	"fmt"

	"github.com/schwarzit/go-template/v2/modules/base"
)

func main() {
	fmt.Println("Hello, World!")

	m, err := base.New()
	if err != nil {
		panic(err)
	}

	fmt.Println(m.TemplateFiles)
}

package main

import (
	"github.com/schwarzit/go-template/v2/gotemplate"
	"github.com/schwarzit/go-template/v2/modules/base"
	"github.com/schwarzit/go-template/v2/view/bubble"
)

func main() {
	m, err := base.New()
	if err != nil {
		panic(err)
	}

	t := gotemplate.NewTemplate(
		bubble.NewTeaView(),
	)

	t.AddModules([]gotemplate.Module{
		m,
	})

	err = t.ExecuteWizard()
	if err != nil {
		panic(err)
	}
}

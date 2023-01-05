package goincv

import (
	"fmt"
	"log"

	"github.com/AlecAivazis/survey"
)

type cmdUI struct {
}

var CmdUI *cmdUI = &cmdUI{}

func (c *cmdUI) MultipleSelection(title string, v ...string) ([]int, []string) {

	var result []int
	var results []string

	prompt := &survey.MultiSelect{
		Message: title + "[空格选择 回车确定]",
		Options: v,
	}
	survey.AskOne(prompt, &result)

	for e := range result {
		results = append(results, v[result[e]])
	}

	if len(result) == 0 {
		return c.MultipleSelection(title, v...)
	}
	return result, results
}

func (c *cmdUI) Confirm(message string) bool {
	name := false
	prompt := &survey.Confirm{
		Message: message,
	}
	survey.AskOne(prompt, &name)
	return name
}

func (c *cmdUI) Select(title string, v ...string) (int, string) {
	result := -1
	prompt := &survey.Select{
		Message: title,
		Options: v,
	}

	err := survey.AskOne(prompt, &result)
	if err != nil {
		log.Fatal(err)
	}

	if result == -1 {
		return c.Select(title, v...)
	}
	return result, v[result]
}

func (c *cmdUI) Text(title string) string {

	result := ""

	prompt := &survey.Input{
		Message: title,
	}
	survey.AskOne(prompt, &result)

	if result == "" {
		fmt.Println("不允许为空")
		return c.Text(title)
	}
	return result
}

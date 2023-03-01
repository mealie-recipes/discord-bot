package main

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed config.yml
var cmds []byte

func ReadCommands() *Commands {
	var c Commands
	err := yaml.Unmarshal(cmds, &c)
	if err != nil {
		panic(err)
	}
	return &c
}

type Commands struct {
	Footer   string    `yaml:"footer"`
	Commands []Command `yaml:"commands"`
}

type Command struct {
	Command     string `yaml:"command"`
	Description string `yaml:"description"`
	Content     string `yaml:"content"`
}

func (c Commands) WrapMessage(msg string) string {
	return fmt.Sprintf("%s\n%s", msg, c.Footer)
}

var BotDebug = `
**Bot Details**

Version %s
`

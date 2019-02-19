package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func readConf(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}
	v := &Config{}
	if err := yaml.Unmarshal(b, v); err != nil {
		return nil, fmt.Errorf("failed to load file as yaml: %s", err)
	}
	return v, nil
}

type Config struct {
	Shell *ShellConfig
}

type ShellConfig struct {
	Type     string
	Prompt   string
	Listen   string
	Commands Commands
}

type Commands map[string]*Command

func (c Commands) Match(in string) *Command {
	cmd, ok := map[string]*Command(c)[in]
	if !ok {
		return nil
	}
	return cmd
}

type Command struct {
	ChangePrompt *string `yaml:"change_prompt"`
	Output       string
	Continues    []string
	NoPrompt     bool `yaml:"no_prompt"`
}

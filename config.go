package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

type Pane struct {
	Root     string   `yaml:"root,omitempty"`
	Type     string   `yaml:"type,omitempty"`
	Commands []string `yaml:"commands"`
}

type Window struct {
	Name        string   `yaml:"name"`
	Root        string   `yaml:"root,omitempty"`
	BeforeStart []string `yaml:"before_start"`
	Panes       []Pane   `yaml:"panes"`
	Commands    []string `yaml:"commands"`
	Layout      string   `yaml:"layout"`
	Manual      bool     `yaml:"manual,omitempty"`
}

type Config struct {
	Session                   string   `yaml:"session"`
	Root                      string   `yaml:"root"`
	BeforeStart               []string `yaml:"before_start"`
	Stop                      []string `yaml:"stop"`
	Windows                   []Window `yaml:"windows"`
	RebalanceWindowsThreshold int      `yaml:"rebalance_panes_after"`
}

func EditConfig(path string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func GetConfig(path string, settings map[string]string) (Config, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	config := string(f)
	config = os.Expand(config, func(v string) string {
		if val, ok := settings[v]; ok {
			return val
		}

		return v
	})

	return ParseConfig(config)
}

func ParseConfig(data string) (Config, error) {
	c := Config{}

	err := yaml.Unmarshal([]byte(data), &c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}

func ListConfigs(dir string) ([]string, error) {
	var result []string
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		return result, err
	}

	for _, file := range files {
		result = append(result, strings.TrimSuffix(file.Name(), path.Ext(file.Name())))
	}

	return result, nil
}

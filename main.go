package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ktr0731/go-fuzzyfinder"

	"gopkg.in/yaml.v2"
)

func main() {
	var (
		cur string
	)
	flag.StringVar(&cur, "cur", "", "current dir path")
	flag.Parse()
	os.Exit(run(cur))
}

type Rule struct {
	Names []string
}

type Menu struct {
	options []string
}

func (m *Menu) load(path string) error {
	buf, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var rule Rule
	if err := yaml.Unmarshal(buf, &rule); err != nil {
		return err
	}
	m.options = rule.Names
	return nil
}

func (m Menu) pick() (string, error) {

	idx, err := fuzzyfinder.Find(m.options, func(i int) string {
		return m.options[i]
	})
	if err != nil {
		return "", err
	}
	return m.options[idx], nil
}

func (m Menu) getName() (string, error) {
	n, err := m.pick()
	if err != nil {
		return "", err
	}
	now := time.Now()
	ts := now.Format("20060102")
	return fmt.Sprintf("%s_%s_", ts, n), nil
}

func run(path string) int {
	y := filepath.Join(path, "rule.yml")
	if _, err := os.Stat(y); err != nil {
		return 1
	}
	var menu Menu
	err := menu.load(y)
	if err != nil {
		return 1
	}
	n, err := menu.getName()
	if err != nil {
		return 1
	}
	p := filepath.Join(path, n)
	if _, err := os.Stat(p); err == nil {
		return 1
	}
	os.Mkdir(p, os.ModePerm)
	return 0
}

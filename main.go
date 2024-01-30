package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	Prefix      string
	Description string
}

type Menu struct {
	options []Rule
}

func (m *Menu) load(path string) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return
	}
	rules := []Rule{}
	if err := yaml.Unmarshal(buf, &rules); err != nil {
		return
	}
	m.options = rules
}

func (m Menu) pick() (string, error) {
	if len(m.options) < 1 {
		return "", fmt.Errorf("no options to pick")
	}
	idx, err := fuzzyfinder.Find(m.options, func(i int) string {
		o := m.options[i]
		return fmt.Sprintf("%s - %s", o.Prefix, o.Description)
	})
	if err != nil {
		return "", err
	}
	return m.options[idx].Prefix, nil
}

func (m Menu) getPrefix() string {
	now := time.Now()
	ts := now.Format("20060102")
	n, err := m.pick()
	if err != nil {
		return ts
	}
	return fmt.Sprintf("%s_%s", ts, n)
}

func (m Menu) getName() (string, error) {
	p := m.getPrefix()
	fmt.Printf("Enter after '%s_': ", p)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	n := scanner.Text()
	n = strings.TrimSpace(n)
	if len(n) < 1 {
		if len(p) == 8 {
			return p, nil
		}
		return "", fmt.Errorf("input cancelled")
	}
	return fmt.Sprintf("%s_%s", p, n), nil
}

func run(path string) int {
	y := filepath.Join(path, "rule.yml")
	var menu Menu
	menu.load(y)
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

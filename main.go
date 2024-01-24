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

func (m *Menu) load(path string) error {
	buf, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	rules := []Rule{}
	if err := yaml.Unmarshal(buf, &rules); err != nil {
		return err
	}
	m.options = rules
	return nil
}

func (m Menu) pick() (string, error) {

	idx, err := fuzzyfinder.Find(m.options, func(i int) string {
		o := m.options[i]
		return fmt.Sprintf("%s - %s", o.Prefix, o.Description)
	})
	if err != nil {
		return "", err
	}
	return m.options[idx].Prefix, nil
}

func (m Menu) getPrefix() (string, error) {
	n, err := m.pick()
	if err != nil {
		return "", err
	}
	now := time.Now()
	ts := now.Format("20060102")
	return fmt.Sprintf("%s_%s_", ts, n), nil
}

func (m Menu) getName() (string, error) {
	p, err := m.getPrefix()
	if err != nil {
		return "", err
	}
	fmt.Printf("Enter after '%s': ", p)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	n := scanner.Text()
	n = strings.TrimPrefix(n, "_")
	n = strings.TrimSpace(n)
	if len(n) < 1 {
		return "", fmt.Errorf("input cancelled")
	}
	return (p + n), nil
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

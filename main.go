package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func init() {
	log.SetFlags(log.LUTC | log.Lmicroseconds)
	log.SetPrefix(filepath.Base(os.Args[0] + ": "))
}

func main() {
	layouts, err := validLayouts()
	switch {
	case err != nil:
		log.Printf("cannot retrieve valid layouts: %s", err)
		return
	case len(layouts) == 0:
		return
	}

	var rotation []string
	for _, l := range os.Args[1:] {
		i := sort.SearchStrings(layouts, l)

		if !(i < len(layouts) && layouts[i] == l) {
			continue // missing
		}

		rotation = append(rotation, l)
	}
	sort.Strings(rotation)

	if len(rotation) < 2 {
		return // nothing to rotate
	}

	current, err := currentLayout()
	if err != nil {
		log.Printf("cannot get current layouts: %s", err)
	}

	index := sort.SearchStrings(rotation, current)
	if index >= len(rotation) || rotation[index] != current {
		// missing
		return
	}

	if index == len(rotation)-1 {
		index = 0
	} else {
		index++
	}

	newLayout := rotation[index]
	if err := setLayout(newLayout); err != nil {
		log.Printf("cannot switch to %q: %s", newLayout, err)
	}
}

func validLayouts() ([]string, error) {
	buf := &bytes.Buffer{}

	cmd := exec.Command("localectl", "--no-pager", "list-x11-keymap-layouts")
	cmd.Stdout = buf
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var ret []string
loop:
	for {
		token, err := buf.ReadString('\n')

		switch err {
		case nil:
			ret = append(ret, strings.TrimSpace(token))
		case io.EOF:
			break loop
		default:
			return nil, err
		}
	}
	sort.Strings(ret)

	return ret, nil
}

func currentLayout() (string, error) {
	const prefix = "layout:"

	buf := &bytes.Buffer{}

	cmd := exec.Command("setxkbmap", "-query")
	cmd.Stdout = buf
	if err := cmd.Run(); err != nil {
		return "", err
	}

	ret := ""

loop:
	for {
		line, err := buf.ReadString('\n')
		switch err {
		case nil:
			break
		case io.EOF:
			break loop
		default:
			return "", err
		}

		if !strings.HasPrefix(line, prefix) {
			continue
		}

		tokens := strings.Split(strings.TrimSpace(line[len(prefix):]), ",")
		if len(tokens) > 0 {
			ret = tokens[0]
			break
		}
	}

	if ret == "" {
		return ret, errors.New("no current layout determined")
	}

	return ret, nil
}

func setLayout(l string) error {
	return exec.Command("setxkbmap", "-layout", l).Run()
}

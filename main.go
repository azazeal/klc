package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/azazeal/exit"
)

func init() {
	log.SetFlags(log.LUTC | log.Lmicroseconds | log.Lmsgprefix)
	log.SetPrefix(filepath.Base(os.Args[0] + ": "))
}

func main() {
	exit.With(run())
}

func run() (err error) {
	var keyboards []*keyboard
	if keyboards, err = detect(); err != nil {
		log.Printf("failed detecting keyboards: %v", err)

		return
	}

	// rotate
	k := keyboards[0]
	copy(keyboards, keyboards[1:])
	keyboards[len(keyboards)-1] = k

	if err = set(keyboards); err != nil {
		log.Printf("failed setting keyboards: %v", err)
	}

	return
}

type keyboard struct {
	layout  string
	variant string
}

const (
	_ = iota
	ecDetect
	ecNoKeyboards
	ecSwitch
)

var errNoKeyboards = exit.Wrap(ecNoKeyboards, errors.New("no keyboards detected"))

func detect() (keyboards []*keyboard, err error) {
	defer func() {
		switch {
		case err != nil:
			err = exit.Wrap(ecDetect, err)
		case len(keyboards) == 0:
			err = errNoKeyboards
		}
	}()

	var buf *bytes.Buffer
	if buf, err = capture("setxkbmap", "-query"); err != nil {
		return
	}

	for {
		var line string
		if line, err = buf.ReadString('\n'); err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			}

			break
		}

		switch {
		case strings.HasPrefix(line, "layout:"):
			line = strings.TrimSpace(strings.TrimPrefix(line, "layout:"))

			for _, l := range strings.Split(line, ",") {
				keyboards = append(keyboards, &keyboard{
					layout: l,
				})
			}
		case strings.HasPrefix(line, "variant:"):
			line = strings.TrimSpace(strings.TrimPrefix(line, "variant:"))

			for i, variant := range strings.Split(line, ",") {
				keyboards[i].variant = variant
			}
		}
	}

	return
}

func set(keyboards []*keyboard) (err error) {
	defer func() {
		if err != nil {
			err = exit.Wrap(ecSwitch, err)
		}
	}()

	// setxkbmap -layout us,gr -variant ,simple -option "grp:alt_space_toggle"

	var layouts, variants []string
	for _, k := range keyboards {
		layouts = append(layouts, k.layout)
		variants = append(variants, k.variant)
	}

	err = exec.Command("setxkbmap", //nolint:gosec // args are provided by the user
		"-layout", strings.Join(layouts, ","),
		"-variant", strings.Join(variants, ","),
	).Run()

	return
}

func capture(name string, args ...string) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}

	cmd := exec.Command(name, args...)
	cmd.Stdout = buf

	return buf, cmd.Run()
}

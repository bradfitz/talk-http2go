// +build ignore

package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	up    = string([]byte{27, 91, 65})
	down  = string([]byte{27, 91, 66})
	left  = string([]byte{27, 91, 68})
	right = string([]byte{27, 91, 67})
	keyn  = "n"
	keyN  = "N"
	keyp  = "p"
	keyP  = "P"
	keyq  = "q"
	ctrlC = string([]byte{3})
	ctrlZ = string([]byte{26})
)

func main() {
	slides, err := readSlidesDat()
	if err != nil {
		log.Fatal(err)
	}
	if len(slides) == 0 {
		slides = []string{"No slides\n"}
	}

	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(0, oldState)

	clear, err := exec.Command("clear").Output()
	if err != nil {
		log.Fatal(err)
	}
	keyBuf := make([]byte, 10)
	n := 0
	for {
		os.Stdout.Write(clear)
		io.WriteString(os.Stdout, "\n ")
		io.WriteString(os.Stdout, strings.TrimSuffix(strings.Replace(slides[n], "\n", "\n ", -1), " "))
		got, err := os.Stdin.Read(keyBuf)
		if err != nil {
			log.Printf("Read key: %v", err)
			return
		}
		key := string(keyBuf[:got])
		switch key {
		case ctrlC, ctrlZ, keyq:
			return
		case up, left, keyp:
			n--
		case down, right, keyn:
			n++
		case keyP:
			n = 0
		case keyN:
			n = len(slides) - 1
		}
		if n < 0 {
			n = 0
		}
		if n >= len(slides) {
			n = len(slides) - 1
		}
	}

}

func readSlidesDat() (slides []string, err error) {
	f, err := os.Open("slides.dat")
	if err != nil {
		return
	}
	defer f.Close()
	var buf bytes.Buffer
	flush := func() {
		s := strings.TrimSpace(buf.String())
		if s != "" {
			slides = append(slides, s+"\n")
		}
		buf.Reset()
	}
	bs := bufio.NewScanner(f)
	for bs.Scan() {
		if bs.Text() == "--" {
			flush()
			continue
		}
		buf.WriteString(bs.Text())
		buf.WriteByte('\n')
	}
	flush()
	return slides, bs.Err()
}

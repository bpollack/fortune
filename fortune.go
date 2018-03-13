package main // import "bitbucket.org/bpollack/fortune"

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	normalFileName    = "freebsd.fortunes"
	offensiveFileName = "offensive.fortunes"
)

func readFortunes(filename string) (fortunes []string, err error) {
	fortuneFile, err := os.Open(filename)
	if err != nil {
		return
	}
	defer fortuneFile.Close()

	reader := bufio.NewReader(fortuneFile)
	fortunes = make([]string, 0, 100)
	currentFortune := ""
	for line, err := reader.ReadString('\n'); err == nil; line, err = reader.ReadString('\n') {
		if line[0] == '%' {
			fortunes = append(fortunes, currentFortune)
			currentFortune = ""
		} else {
			currentFortune += line
		}
	}

	if len(currentFortune) > 0 {
		fortunes = append(fortunes, currentFortune)
	}

	if err == io.EOF {
		err = nil
	}

	return
}

func executablePath() string {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		panic(err)
	}
	path, err = filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return filepath.Dir(path)
}

func main() {
	offensive := flag.Bool("o", false, "be offensive")
	flag.Parse()
	fortunes, err := readFortunes(filepath.Join(executablePath(), "fortunes", normalFileName))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to load normal fortunes: %v\n", err)
		os.Exit(1)
	}
	if *offensive {
		offensiveFortunes, err := readFortunes(filepath.Join(executablePath(), offensiveFileName))
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to load offensive fortunes: %v\n", err)
			os.Exit(1)
		}
		fortunes = append(fortunes, offensiveFortunes...)
	}

	fmt.Print(fortunes[rand.Intn(len(fortunes))])
}

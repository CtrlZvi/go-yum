package yum

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func Provides(pattern string) (map[string][]Package, error) {
	out, err := exec.Command("yum", "provides", "-q", fmt.Sprintf("%v", pattern)).Output()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		emptyLineCount := 0
		for emptyLineCount < 3 && advance < len(data) {
			lineAdvance, lineToken, err2 := bufio.ScanLines(data[advance:], atEOF)
			if err2 != nil {
				err = err2
				return
			}

			if string(lineToken) == "" {
				emptyLineCount += 1
			} else {
				emptyLineCount = 0

				for _, b := range lineToken {
					token = append(token, b)
				}
				token = append(token, '\n')
			}
			advance += lineAdvance
		}

		return
	})

	pkgs := make(map[string][]Package)
	for scanner.Scan() {
		result := scanner.Text()
		scanner := bufio.NewScanner(strings.NewReader(result))

		pkg := Package{}
		for scanner.Scan() {
			line := scanner.Text()
			var filename string

			splitLine := strings.SplitN(line, ":", 3)
			switch strings.Trim(splitLine[0], " ") {
			case "Repo":
				pkg.Repository = strings.Trim(splitLine[1], " ")
			case "Matched from":
			case "Filename":
				filename = strings.Trim(splitLine[1], " ")
			default:
				if len(splitLine) == 3 {
					pkg.Epoch, err = strconv.ParseInt(splitLine[0], 10, 32)
					if err != nil {
						splitLine = []string{splitLine[0], strings.Join([]string{splitLine[1], splitLine[2]}, ":")}
					} else {
						splitLine = splitLine[1:]
					}
				}
				tokens := strings.Split(strings.Trim(splitLine[0], " "), "-")
				pkg.Name = tokens[0]
				pkg.Version = tokens[1]
				components := strings.Split(tokens[2], ".")
				pkg.Release = strings.Join([]string{components[0], components[1]}, ".")
				pkg.Architecture = components[2]
				pkg.Summary = strings.Trim(splitLine[1], " ")
			}

			if pkg.Repository[0] != '@' {
				oldRepository := pkg.Repository
				pkg.Repository = pkg.Repository[1:]
				alreadyExists := false
				for _, p := range pkgs[filename] {
					if p == pkg {
						alreadyExists = true
					}
				}

				if !alreadyExists {
					pkg.Repository = oldRepository
					pkgs[filename] = append(pkgs[filename], pkg)
				}
			}
		}
	}
	return pkgs, nil
}

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

			splitLine := strings.SplitN(line, ":", 2)
			switch strings.Trim(splitLine[0], " ") {
			case "Repo":
				pkg.Repository = strings.Trim(splitLine[1], " ")
			case "Matched from":
			case "Filename":
				filename := strings.Trim(splitLine[1], " ")
				pkgs[filename] = append(pkgs[filename], pkg)
			default:
				if len(splitLine) == 3 {
					pkg.Epoch, err = strconv.ParseInt(splitLine[0], 10, 32)
					if err != nil {
						return nil, err
					}
					splitLine = splitLine[1:]
				}
				tokens := strings.Split(strings.Trim(splitLine[0], " "), "-")
				pkg.Name = tokens[0]
				pkg.Version = tokens[1]
				components := strings.Split(tokens[2], ".")
				pkg.Release = strings.Join([]string{components[0], components[1]}, ".")
				pkg.Architecture = components[2]
				pkg.Summary = strings.Trim(splitLine[1], " ")
			}
		}
	}
	return pkgs, nil
}

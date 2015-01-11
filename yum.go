package yum

import (
	"bufio"
	"bytes"
	"fmt"
	// "net/url"
	"os/exec"
	// "strconv"
	"strings"
)

/*
func List(pattern string) (installed []Package, available []Package, err error) {
  out, err := exec.Command(fmt.Sprintf("yum list %v", pattern)).Output()
  return
}

func Info(pattern string) ([]Package, error) {
  out, err := exec.Command(fmt.Sprintf("yum info -q %v", pattern)).Output()
  if err != nil {
    return nil, err
  }

  scanner := bufio.NewScanner(bytes.NewReader(out))
  pkg := Package{}
  for scanner.Scan() {
    line := scanner.Text()
    if line == "Available Packages" ||
      line == "Installed Packages" ||
      line == "" {
      continue
    }

    splitLine := strings.SplitN(line, ":", 2)
    switch splitLine[0] {
    case "Name        ":
      pkg.Name = splitLine[1]
    case "Arch        ":
      pkg.Architecture = splitLine[1]
    case "Epoch       ":
      pkg.Epoch = splitLine[1]
    case "Version     ":
      pkg.Version = splitLine[1]
    case "Release     ":
      pkg.Release = splitLine[1]
    case "Size        ":
      pieces := strings.Split(splitLine[1], " ")
      pkg.Size, err = strconv.ParseInt(pieces[0], 10, 64)
      switch pieces[1] {
      case "k":
        pkg.Size *= 1024
      case "M":
        pkg.Size *= 1024 * 1024
      }
    case "Repo        ":
      pkg.Repository = splitLine[1]
    case "Summary     ":
      pkg.Summary = splitLine[1]
    case "URL         ":
      pkg.URL, err = url.Parse(splitLine[1])
    case "License     ":
      pkg.License = splitLine[1]
    case "Description ":
      pkg.Description = splitLine[1]
    }
  }
  return nil, nil
}
*/

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
				tokens := strings.Split(strings.Trim(splitLine[0], " "), "-")
				pkg.Name = tokens[0]
				components := strings.Split(tokens[1], ":")
				if len(components) == 2 {
					pkg.Epoch = components[0]
					pkg.Version = components[1]
				} else {
					pkg.Version = components[0]
				}
				components = strings.Split(tokens[2], ".")
				pkg.Release = strings.Join([]string{components[0], components[1]}, ".")
				pkg.Architecture = components[2]
				pkg.Summary = strings.Trim(splitLine[1], " ")
			}
		}
	}
	return pkgs, nil
}

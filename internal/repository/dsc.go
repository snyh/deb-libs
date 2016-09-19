package packages

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type DSCFile map[string]string

func NewDSCFile(r io.Reader) (DSCFile, error) {
	splitFn := func(data []byte, atEOF bool) (advance int, toke []byte, err error) {
		l := len(data)
		for i, c := range data {
			if c == '\n' {
				if i+1 < l && (data[i+1] != ' ' && data[i+1] != '\t') {
					return i + 1, data[:i], nil
				}

				if i+1 == l && atEOF {
					return i + 1, data[:i], nil
				}
			}
		}
		if !atEOF {
			return 0, nil, nil
		}
		if atEOF && l != 0 {
			return l, data, nil
		}

		return l, data, fmt.Errorf("End of file")
	}

	f := make(DSCFile)
	s := bufio.NewScanner(r)
	s.Split(splitFn)
	for s.Scan() {
		line := s.Text()
		d := strings.SplitN(line, ":", 2)
		if len(d) != 2 {
			continue
			return nil, fmt.Errorf("NewDSCFile there has %d separators at:%q", len(d), line)
		}
		f[strings.ToLower(d[0])] = strings.Trim(d[1], " \n")
	}
	return f, nil
}

func (d DSCFile) GetString(key string) string {
	return d[strings.ToLower(key)]
}

func (d DSCFile) GetArrayString(key string) []string {
	var r []string
	for _, c := range strings.Split(d.GetString(key), " ") {
		r = append(r, c)
	}
	return r
}

func (d DSCFile) GetMultiline(key string) []string {
	return strings.Split(d[strings.ToLower(key)], "\n")
}

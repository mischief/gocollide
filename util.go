package main

import (
  "bufio"
  "os"
  "io"
  "log"
)

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func LoadDict(dict string) ([]string, error) {
  f, err := os.Open(dict)
  if err != nil {
    return nil, err
  }
  defer f.Close()

  bufreader := bufio.NewReader(f)

  lines := 1
  out := make([]string, 10)

  for {
    line, err := Readln(bufreader)

    if err == io.EOF && line == "" {
      // end of dictionary
      break
    }

    if line == "" {
      log.Printf("Empty line at line %d", lines)
    }

    if err != nil {
      // some other kind of error
      return nil, err
    }

    out = append(out, line)
    lines++
  }

  return out, nil
}

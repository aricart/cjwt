package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	flag.Usage = usage
	flag.Parse()

	switch flag.NArg() {
	case 0:
		s, err := os.Stdin.Stat()
		if err != nil {
			panic(err)
		}
		if s.Mode()&os.ModeNamedPipe == 0 {
			// no data
			usage()
		}
		d, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		process(d)
	case 1:
		d, err := ioutil.ReadFile(flag.Arg(0))
		if err != nil {
			panic(err)
		}
		process(d)
	default:
		usage()
	}
}

func process(data []byte) {
	if len(data) == 0 {
		usage()
	}
	chunks := bytes.Split(data, []byte{'.'})
	if len(chunks) != 3 {
		exit("data is not a jwt")
	}

	for i, c := range chunks {
		if i == 2 {
			fmt.Println(string(c))
			continue
		}
		d, err := base64.RawURLEncoding.DecodeString(string(c))
		if err != nil {
			exit("error decoding base64: %v", err)
		}

		m := make(map[string]interface{})
		if err := json.Unmarshal(d, &m); err != nil {
			exit("error parsing json: %v", err)
		}

		f, err := json.MarshalIndent(m, "", " ")
		if err != nil {
			exit("error formatting json: %v", err)
		}
		fmt.Println(string(f))
	}
}

func usage() {
	exit("cjwt [filepath] | [stdin]")
}

func exit(format string, a ...interface{}) {
	if len(a) > 0 {
		fmt.Println(format, a)
	} else {
		fmt.Println(format)
	}
	os.Exit(1)
}

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var service string

func init() {
	flag.StringVar(&service, "service", "", "filter which service to see")
}

func main() {

	flag.Parse()
	var b strings.Builder

	service := strings.ToLower(service)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		m := make(map[string]any)
		err := json.Unmarshal([]byte(s), &m)
		if err != nil {
			if service == "" {
				fmt.Println(s)
			}
			continue
		}
		if service != "" && strings.ToLower(m["service"].(string)) != service {
			continue
		}

		traceID := "00000000-0000-0000-0000-000000000000"
		if v, ok := m["trace_id"]; ok {
			traceID = fmt.Sprintf("%v", v)
		}

		// {"time":"2023-06-01T17:21:11.13704718Z","level":"INFO","msg":"startup","service":"SALES-API","GOMAXPROCS":1}

		// Build out the know portions of the log in the order
		// I want them in.

		b.Reset()
		b.WriteString(fmt.Sprintf("%s: %s: %s: %s: %s: %s: ", m["service"], m["time"], m["file"], m["level"], traceID, m["msg"]))

		for k, v := range m {
			switch k {
			case "service", "time", "file", "level", "trace_id", "msg":
				continue
			}
			b.WriteString(fmt.Sprintf("%s[%v]:", k, v))
		}
		out := b.String()
		fmt.Println(out[:len(out)-2])

	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

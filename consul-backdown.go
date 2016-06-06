package main

import (
	"github.com/docopt/docopt-go"
	"github.com/hashicorp/consul/api"

	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"sort"
	"strings"
)

const version = "consul-backdown 1.0.0"
const usage = `consul-backdown.

Usage:
  consul-backdown backup [--url <url> | -u <url>] [-t <token> | --token <token>]
				  [ -d | --debug ] [ -s <sep> | --separator <sep> ]
  consul-backdown restore [-u <url> | --url <url>] [-t <token> | --token <token>]
				  [ -d | --debug ] [ -i <file> | --input <file ]
				  [ -s <sep> | --separator <separator> ]
  consul-backdown -h | --help
  consul-backdown --version

Options:
  -u <url> --url=<url>         Consul's HTTP address [default: 127.0.0.1:8500].
  -t <token> --token=<token>   Consul's token [default: ].
  -i <file> --input=<file>     Input file for restore, defaults to stdin.
  -s <sep> --separator=<sep>   Key/Value separator [default: |||].
  -d --debug                   Show more values in the output (don't use it on production) [default: false].
  -h --help                    Show this screen.
  --version                    Show version.
`

type KVPairs api.KVPairs

func (a KVPairs) Len() int {
	return len(a)
}

func (a KVPairs) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a KVPairs) Less(i, j int) bool {
	return a[i].CreateIndex < a[j].CreateIndex
}

func backup(ip, token, separator string) {
	config := api.DefaultConfig()
	config.Address = ip
	config.Token = token

	client, _ := api.NewClient(config)
	kv := client.KV()

	pairs, _, err := kv.List("/", nil)
	if err != nil {
		panic(err)
	}

	sort.Sort(KVPairs(pairs))

	for _, element := range pairs {
		encoded_value := base64.StdEncoding.EncodeToString(element.Value)
		os.Stdout.WriteString(element.Key + separator + encoded_value + "\n")
	}
}

func restore(ip string, token string, separator string, input *os.File) {
	config := api.DefaultConfig()
	config.Address = ip
	config.Token = token

	client, _ := api.NewClient(config)
	kv := client.KV()

	scanner := bufio.NewScanner(input)

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	for scanner.Scan() {
		element := scanner.Text()
		kvp := strings.Split(element, separator)

		if len(kvp) > 1 {
			decoded_value, decode_err := base64.StdEncoding.DecodeString(kvp[1])
			if decode_err != nil {
				panic(decode_err)
			}

			p := &api.KVPair{Key: kvp[0], Value: decoded_value}
			_, err := kv.Put(p, nil)
			if err != nil {
				panic(err)
			}
		}
	}
}

func main() {
	args, _ := docopt.Parse(usage, nil, true, version, false)
	debug := false
	url, _ := args["--url"].(string)
	token, _ := args["--token"].(string)
	separator, _ := args["--separator"].(string)

	if args["--debug"] == true {
		debug = true
		fmt.Println("Arguments being passed:", args)
	}

	if args["backup"] == true {
		if debug {
			fmt.Println("Backing up from server:", url)
		}
		backup(url, token, separator)
	}

	if args["restore"] == true {
		if debug {
			fmt.Println("Restoring into server:", url)
		}

		if args["--input"] != nil {
			inputArg, _ := args["--input"].(string)

			if debug {
				fmt.Println("Will read backup from file given as argument", inputArg)
			}

			input, err := os.Open(inputArg)
			if err != nil {
				panic(err)
			}
			defer input.Close()
			restore(url, token, separator, input)
		} else {
			if debug {
				fmt.Println("Will read backup from STDIN")
			}

			input := os.Stdin
			restore(url, token, separator, input)
		}
	}
}

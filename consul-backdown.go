package main

import (
	"github.com/docopt/docopt-go"
	"github.com/hashicorp/consul/api"

	"encoding/base64"
	"fmt"
	"os"
	"sort"
)

const version = "consul-backdown 1.0.0"
const usage = `consul-backdown.

Usage:
  consul-backdown backup [--url <url> | -u <url>] [-t <token> | --token <token>]
				  [ -d | --debug ]
  consul-backdown restore [-u <url> | --url <url>] [-t <token> | --token <token>]
				  [ -d | --debug ]
  consul-backdown -h | --help
  consul-backdown --version

Options:
  -u <url> --url=<url>           Consul's HTTP address [default: 127.0.0.1:8500].
  -t <token> --token=<token>     Consul's token [default: ].
  -d --debug                     Show more values in the output (don't use it on production) [default: false].
  -h --help                      Show this screen.
  --version                      Show version.
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

func backup(ip, token string) {
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
		os.Stdout.WriteString(element.Key + ":" + encoded_value + "\n")
	}
}

func main() {
	args, _ := docopt.Parse(usage, nil, true, version, false)
	debug := false
	url, _ := args["--url"].(string)
	token, _ := args["--token"].(string)

	if args["--debug"] == true {
		debug = true
		fmt.Println("Arguments being passed:", args)
	}
	if args["backup"] == true {
		if debug {
			fmt.Println("Backing up from server:", url)
		}
		backup(url, token)
	}
	if args["restore"] == true {
	}
}

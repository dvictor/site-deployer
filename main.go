package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)
import "net/http"

type Config struct {
	ServerAddress string                `json:"server_address"`
	Sites         map[string]SiteConfig `json:"sites"`
}

type SiteConfig struct {
	Path string     `json:"path"`
	Cmd  [][]string `json:"cmd"`
}

var mq chan string
var config Config

func main() {
	bytes, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		panic(err)
	}

	mq = make(chan string)
	go func() {
		for true {
			name := <-mq
			cfg := config.Sites[name]
			for _, a := range cfg.Cmd {
				cmd := exec.Command(a[0], a[1:]...)
				cmd.Dir = cfg.Path
				out, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Println("Error deploying", name)
				}
				fmt.Println(strings.TrimSpace(string(out)))
			}
		}
	}()

	http.HandleFunc("/", Handle)
	fmt.Printf("Server started at http://%s\n", config.ServerAddress)
	err = http.ListenAndServe(config.ServerAddress, nil)
	if err != nil {
		panic(err)
	}
}

func Handle(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[1:]
	if _, ok := config.Sites[name]; !ok {
		w.WriteHeader(400)
		_, _ = fmt.Fprintf(w, "Site not found: %s", name)
		return
	}
	mq <- name
	_, _ = fmt.Fprintf(w, "Deployment queued: %s", name)
}

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/go-ini/ini"
)

func main() {
	// GetMyGerrit is a tool to aggregate all your patches and reviews from
	// different repos and show them in a single webpage with links to them.
	fmt.Println("Welcome to GetMyGerrit")

	cfg, err := ini.Load("config/repos.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	// Classic read of values, default section can be represented as empty string
	// fmt.Println("App Mode:", cfg.Section("").Key("app_mode").String())
	fmt.Println("Hola", cfg.SectionStrings())

	sectionNames := cfg.SectionStrings()
	for _, section := range sectionNames {
		fmt.Println("Ricky", section)
		if cfg.Section(section).Name() != "DEFAULT" {
			url := cfg.Section(section).Key("url").String()
			user := cfg.Section(section).Key("user").String()
			fmt.Println("Ricky URL", url)
			var finalURL []string
			finalURL = append(finalURL, url, "/changes/?q=reviewer:", user, "+status:open")
			fmt.Println("FinalURL:", strings.Join(finalURL, ""))

			req, err := http.NewRequest("GET", strings.Join(finalURL, ""), nil)
			if err != nil {
				//handle error
				fmt.Println("Wrong request")
			}
			fmt.Println("Great request:", req)
			res, _ := http.DefaultClient.Do(req)

			if res == nil {
				fmt.Println("Booooh")
			}
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)

			fmt.Println(string(body))
		}
	}
}

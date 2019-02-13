package main

import (
	"encoding/json"
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

	// This struct represents the information fetch from Gerrit for one patch.
	type PatchData struct {
		ChangeID string `json:"change_id"`
		Status   string `json:"status"`
		Number   int    `json:"_number"`
	}

	webContent := map[string][]string{}

	cfg, err := ini.Load("config/repos.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	// Classic read of values, default section can be represented as empty string
	fmt.Println("Reading through the config file", cfg.SectionStrings())

	for _, section := range cfg.SectionStrings() {
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
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				//handle error
				fmt.Println("Wrong response")
			} else {
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					panic(err)
				}
				var myPatches []PatchData
				err = json.Unmarshal(body[5:], &myPatches)
				if err != nil {
					fmt.Println(err)
				}
				for _, myPatch := range myPatches {
					fmt.Println("Response: ", section, url+"/"+fmt.Sprint(myPatch.Number))
					webContent[cfg.Section(section).Name()] = append(webContent[cfg.Section(section).Name()], (url + "/" + fmt.Sprint(myPatch.Number)))
				}
			}
			fmt.Println("FINISHED: ", webContent)
		}
	}
}

package postrepo

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func request(url, token string, response interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	return nil
}

func getTreeSHA(token string) (string, error) {
	var response struct {
		Commit struct {
			Commit struct {
				Tree struct {
					URL string `json:"url"`
					Sha string `json:"sha"`
				} `json:"tree"`
			} `json:"commit"`
		} `json:"commit"`
	}

	url := "https://api.github.com/repos/MihaiBlebea/mihaiblebea-content/branches/master"
	err := request(url, token, &response)
	if err != nil {
		return "", err
	}

	return response.Commit.Commit.Tree.URL, nil
}

func getFileURLs(url, token string) ([]string, error) {
	var response struct {
		Tree []struct {
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"tree"`
	}

	err := request(url, token, &response)
	if err != nil {
		return []string{}, err
	}

	var urls []string
	for _, url := range response.Tree {
		if url.Type != "blob" {
			continue
		}
		urls = append(urls, url.URL)
	}

	return urls, nil
}

func getFileContent(url, token string) (string, error) {
	var response struct {
		Content string `json:"content"`
	}

	err := request(url, token, &response)
	if err != nil {
		return "", err
	}

	b, err := base64.StdEncoding.DecodeString(response.Content)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

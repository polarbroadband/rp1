package gitlib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"

	"net/http"

	"encoding/json"

	"github.com/kr/pretty"
	"github.com/sirupsen/logrus"
)

type GitLab struct {
	URL         string
	AccessToken string
	Client      *http.Client
}

func (g *GitLab) Get(ep string, rb interface{}) (*http.Header, error) {
	url := fmt.Sprintf("%s/%s", g.URL, ep)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build request to %s: %v", url, err)
	}
	req.Header.Set("Authorization", "Bearer "+g.AccessToken)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := g.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unnable to GET data from %s: %v", url, err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(rb); err != nil {
		return nil, fmt.Errorf("unnable to decode response from %s: %v", url, err)
	}
	return &resp.Header, nil
}

func (g *GitLab) Post(ep string, rb interface{}, request map[string]interface{}) error {
	url := fmt.Sprintf("%s/%s", g.URL, ep)
	if rb == nil {
		var xm interface{}
		rb = &xm
	}
	if request == nil {
		request = map[string]interface{}{}
	}
	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("unable to encode request data %+v: %v", request, err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("unable to build request to %s: %v", url, err)
	}
	req.Header.Set("Authorization", "Bearer "+g.AccessToken)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := g.Client.Do(req)
	if err != nil {
		return fmt.Errorf("unnable to POST data from %s: %v", url, err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(rb); err != nil {
		return fmt.Errorf("fail to decode response from %s: %v", url, err)
	}
	return nil
}

func (g *GitLab) GetFile(ep string, rb *[]byte) error {
	url := fmt.Sprintf("%s/%s", g.URL, ep)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("unable to build request to %s: %v", url, err)
	}
	req.Header.Set("Authorization", "Bearer "+g.AccessToken)
	req.Header.Set("Accept", "application/octet-stream")

	resp, err := g.Client.Do(req)
	if err != nil {
		return fmt.Errorf("unnable to GET data from %s: %v", url, err)
	}
	defer resp.Body.Close()

	if *rb, err = ioutil.ReadAll(resp.Body); err != nil {
		return fmt.Errorf("fail to read content from %s: %v", url, err)
	}
	return nil
}

type GitLabCommit struct {
	*GitLab
	Repo   float64
	Commit string

	Log *logrus.Entry
}

func (g *GitLabCommit) GetRepoRawFiles(filter *regexp.Regexp) ([]GitLabBlob, error) {
	tree := []GitLabBlob{}
	page := 1
	for {
		treePaged := []GitLabBlob{}
		ep := fmt.Sprintf("projects/%v/repository/tree?ref=%s&recursive=true&pagination=keyset&page=%v&per_page=1&order_by=path&sort=asc", g.Repo, g.Commit, page)
		h, err := g.GitLab.Get(ep, &treePaged)
		if err != nil {
			g.Log.WithError(err).Errorf("unable to retrieve repository tree, page %v", page)
			return nil, err
		}
		tree = append(tree, treePaged...)
		if h.Get("X-Next-Page") == "" {
			break
		} else {
			page++
		}
	}
	pretty.Printf("\n--- TREE STRUCT ---\n%# v\n\n", tree)

	resp := []GitLabBlob{}
	for _, b := range tree {
		if filter == nil || filter.MatchString(b.Name) {
			g.Log.Info(b.Path)
			b.Content = []byte{}
			if err := g.GetFile(fmt.Sprintf("projects/%v/repository/blobs/%s/raw", g.Repo, b.ID), &b.Content); err != nil {
				g.Log.WithError(err).Errorf("unable to load blob %s content", b.Path)
				return nil, err
			}
			g.Log.Info(b.Content)
			resp = append(resp, b)
		}
	}
	return resp, nil
}

type GitLabBlob struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	Type    string `json:"type"` // "blob" "tree"
	Content []byte
}

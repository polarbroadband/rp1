package gitlib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"

	"net/http"
	"net/url"

	"encoding/json"

	"github.com/kr/pretty"
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
	*GitLabCheckoutCommit
	Repo float64
}

func (g *GitLabCommit) GetRepoRawFiles(filter *regexp.Regexp, checkout bool) ([]GitLabBlob, error) {
	url := fmt.Sprintf("projects/%v/repository/tree?", g.Repo)
	if !checkout {
		url += fmt.Sprintf("ref=%s&", g.SHA)
	}
	tree := []GitLabBlob{}
	page := 1
	for {
		treePaged := []GitLabBlob{}
		ep := fmt.Sprintf("%srecursive=true&pagination=keyset&page=%v&per_page=100&order_by=path&sort=asc", url, page)
		h, err := g.GitLab.Get(ep, &treePaged)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve repository tree, page %v, %v", page, err)
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
		if b.Type == "tree" {
			continue
		}
		if filter == nil || filter.MatchString(b.Name) {

			pretty.Printf("\n--- TREE PATH ---\n%# v\n\n", b.Path)

			b.Content = []byte{}
			if err := g.GetFile(fmt.Sprintf("projects/%v/repository/blobs/%s/raw", g.Repo, b.ID), &b.Content); err != nil {
				return nil, fmt.Errorf("unable to load blob %s content, %v", b.Path, err)
			}

			pretty.Printf("\n--- FILE CONTENT ---\n%# v\n\n", b.Content)

			resp = append(resp, b)
		}
	}
	return resp, nil
}

// GetAddedFiles download the content of every added files of this commit
func (g *GitLabCommit) GetAddedFiles() (map[string]*[]byte, error) {
	if len(g.Added) == 0 {
		return nil, nil
	}
	files := map[string]*[]byte{}
	for _, f := range g.Added {
		var content []byte
		ep := fmt.Sprintf("projects/%v/repository/files/%s/raw?ref=%s", g.Repo, url.QueryEscape(f), g.SHA)
		if err := g.GetFile(ep, &content); err != nil {
			return nil, err
		}
		files[f] = &content
	}
	return files, nil
}

func (g *GitLabCommit) GetModifiedFiles() (map[string]*[]byte, error) {
	if len(g.Modified) == 0 {
		return nil, nil
	}
	files := map[string]*[]byte{}
	for _, f := range g.Modified {
		var content []byte
		ep := fmt.Sprintf("projects/%v/repository/files/%s/raw?ref=%s", g.Repo, url.QueryEscape(f), g.SHA)
		if err := g.GetFile(ep, &content); err != nil {
			return nil, err
		}
		files[f] = &content
	}
	return files, nil
}

func (g *GitLabCommit) GetRemovedFiles() (map[string]*[]byte, error) {
	if len(g.Removed) == 0 {
		return nil, nil
	}
	files := map[string]*[]byte{}
	for _, f := range g.Removed {
		var content []byte
		ep := fmt.Sprintf("projects/%v/repository/files/%s/raw?ref=%s", g.Repo, url.QueryEscape(f), g.PreSHA)
		if err := g.GetFile(ep, &content); err != nil {
			return nil, err
		}
		files[f] = &content
	}
	return files, nil
}

type GitLabBlob struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	Type    string `json:"type"` // "blob" "tree"
	Content []byte
}

type GitLabCheckoutCommit struct {
	SHA      string   `json:"id"`
	Added    []string `json:"added"`
	Modified []string `json:"modified"`
	Removed  []string `json:"removed"`
	PreSHA   string
}

type GitLabWebhookEvent struct {
	PreviousCommit string                  `json:"before"`
	CheckoutCommit string                  `json:"checkout_sha"`
	Commits        []*GitLabCheckoutCommit `json:"commits"`
	Repo           struct {
		ID   float64 `json:"id"`
		Name string  `json:"name"`
	} `json:"project"`
}

func (w *GitLabWebhookEvent) GetCheckoutCommit(g *GitLab) *GitLabCommit {
	for _, c := range w.Commits {
		if c.SHA == w.CheckoutCommit {
			c.PreSHA = w.PreviousCommit
			return &GitLabCommit{
				GitLab:               g,
				GitLabCheckoutCommit: c,
				Repo:                 w.Repo.ID,
			}
		}
	}
	return nil
}

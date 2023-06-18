package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/gobwas/glob"
	json "github.com/json-iterator/go"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mangalorg/mangalcli/fs"
	"github.com/spf13/afero"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
)

var (
	gitHubAPIURL, _ = url.Parse("https://api.github.com")
)

type downloadCmd struct {
	Owner  string `help:"Owner of the GitHub repo" default:"mangalorg"`
	Branch string `help:"Repo branch" default:"main"`
	Repo   string `help:"GitHub repo that contains lua providers" default:"saturno"`
	Glob   string `help:"Glob pattern of the files to download" default:"*.lua"`
	Dir    string `help:"Output directory" default:"." type:"existingdir"`
	All    bool   `help:"Download all files that match the glob pattern"`
}

type gitHubTreeItem struct {
	Path string `json:"path"`
	URL  string `json:"url"`

	// Type of the tree item.
	//
	// `blob` is a file.
	// `tree` is a directory.
	Type string `json:"type"`
}

func (g gitHubTreeItem) downloadAndDecodeContents() ([]byte, error) {
	request, err := newGitHubApiRequest(g.URL)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}

	buffer, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var blob struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}

	err = json.Unmarshal(buffer, &blob)
	if err != nil {
		return nil, err
	}

	switch blob.Encoding {
	case "base64":
		encoding := base64.RawStdEncoding
		return encoding.DecodeString(blob.Content)
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", blob.Encoding)
	}
}

type gitHubTree struct {
	Tree []gitHubTreeItem `json:"tree"`
}

func newGitHubApiRequest(URL string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("Accept", "application/vnd.github+json")

	if token, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	return request, nil
}

func (d *downloadCmd) getTreeSHA() (string, error) {
	URL := gitHubAPIURL.JoinPath("repos", d.Owner, d.Repo, "branches", d.Branch)
	request, err := newGitHubApiRequest(URL.String())
	if err != nil {
		return "", err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", errors.New(response.Status)
	}

	buffer, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var r struct {
		Commit struct {
			SHA string `json:"sha"`
		} `json:"commit"`
	}

	err = json.Unmarshal(buffer, &r)
	if err != nil {
		return "", err
	}

	return r.Commit.SHA, nil
}

func (d *downloadCmd) getTree() (gitHubTree, error) {
	SHA, err := d.getTreeSHA()
	if err != nil {
		return gitHubTree{}, err
	}

	URL := gitHubAPIURL.JoinPath("repos", d.Owner, d.Repo, "git", "trees", SHA)

	params := url.Values{}
	params.Set("recursive", "1")

	URL.RawQuery = params.Encode()

	request, err := newGitHubApiRequest(URL.String())
	if err != nil {
		return gitHubTree{}, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return gitHubTree{}, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return gitHubTree{}, errors.New(response.Status)
	}

	buffer, err := io.ReadAll(response.Body)
	if err != nil {
		return gitHubTree{}, err
	}

	var tree gitHubTree

	err = json.Unmarshal(buffer, &tree)
	if err != nil {
		return gitHubTree{}, err
	}

	return tree, nil
}

func (d *downloadCmd) getFilteredFiles() ([]*gitHubTreeItem, error) {
	tree, err := d.getTree()
	if err != nil {
		return nil, err
	}

	var items []*gitHubTreeItem

	pattern, err := glob.Compile(d.Glob)
	if err != nil {
		return nil, err
	}

	for _, item := range tree.Tree {
		if item.Type == "blob" && pattern.Match(item.Path) {
			items = append(items, &item)
		}
	}

	return items, nil
}

func (d *downloadCmd) Run() error {
	files, err := d.getFilteredFiles()
	if err != nil {
		return err
	}

	var filesToDownload []*gitHubTreeItem

	if d.All {
		filesToDownload = files
	} else {
		indexes, err := fuzzyfinder.FindMulti(files, func(i int) string {
			return filepath.Base(files[i].Path)
		}, fuzzyfinder.WithHeader("Select providers to download. TAB to select multiple"))

		if err != nil {
			if errors.Is(err, fuzzyfinder.ErrAbort) {
				os.Exit(1)
			}

			return err
		}

		filesToDownload = make([]*gitHubTreeItem, len(indexes))
		for i, index := range indexes {
			filesToDownload[i] = files[index]
		}
	}

	var wg sync.WaitGroup

	wg.Add(len(filesToDownload))

	for _, file := range filesToDownload {
		go func(file *gitHubTreeItem) {
			defer wg.Done()

			log.Info("downloading", "file", file.Path)

			reader, err := file.downloadAndDecodeContents()
			if err != nil {
				log.Error(err, "file", file.Path)
				return
			}

			filename := filepath.Base(file.Path)

			err = afero.WriteFile(fs.FS, filepath.Join(d.Dir, filename), reader, 0755)
			if err != nil {
				log.Error(err, "file", file.Path)
				return
			}

			log.Info("done", "file", file.Path)
		}(file)
	}

	wg.Wait()

	return nil
}

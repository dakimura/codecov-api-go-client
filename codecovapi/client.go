package codecovapi

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type GitHosting int

const (
	GitHub GitHosting = iota
	GitLab
	BitBucket
)

const (
	BaseURL = "https://codecov.io/api"
)

func (gh GitHosting) String() string {
	return [...]string{"gh", "gl", "bb"}[gh]
}

type Client struct {
	Token      string
	HTTPClient *http.Client
}

func NewClient(token string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = new(http.Client)
	}

	return &Client{
		Token:      token,
		HTTPClient: httpClient,
	}
}

type Overview struct {
	Commits []Commit `json:"commits"`
	Repo    Repo     `json:"repo"`
	Meta    Meta     `json:"meta"`
	Owner   Owner    `json:"owner"`
	Commit  Commit   `json:"commit"`
}

type Commit struct {
	PullID       int    `json:"pullid"`
	Author       Author `json:"author"`
	Timestamp    string `json:"timestamp"`
	ParentTotals Totals `json:"parent_totals"`
	State        string `json:"state"`
	Totals       Totals `json:"totals"`
	CommitID     string `json:"commitid"`
	CIPassed     bool   `json:"ci_passed"`
	Branch       string `json:"branch"`
	Message      string `json:"message"`
	Merged       bool   `json:"merged"`
}

type Author struct {
	Username  string `json:"username"`
	ServiceID string `json:"service_id"`
	Name      string `json:"name"`
	Service   string `json:"service"`
	Email     string `json:"email"`
}

type Totals struct {
	C             int           `json:"C"` // ??? not explained in API doc... https://docs.codecov.io/reference#totals
	CoverageRatio string        `json:"c"`
	FilesCount    int           `json:"f"`
	LinesCount    int           `json:"n"`
	HitsCount     int           `json:"h"`
	MissedCount   int           `json:"m"`
	PartialsCount int           `json:"p"`
	BranchesCount int           `json:"b"`
	MethodsCount  int           `json:"d"`
	MessagesCount int           `json:"M"`
	SessionsCount int           `json:"s"`
	Diff          []interface{} `json:"diff"` // see https://docs.codecov.io/reference#totals
}

type Meta struct {
	Status int `json:"status"`
}

type Owner struct {
	Username         string `json:"username"`
	RemainingCredits int    `json:"remaining_credits"`
	Name             string `json:"name"`
	Service          string `json:"service"`
	UpdateStamp      string `json:"updatestamp"`
	AvatarURL        string `json:"avatar_url"`
	ServiceID        string `json:"service_id"`
}

type Repo struct {
	UsingIntegration bool   `json:"using_integration"`
	Name             string `json:"name"`
	Language         string `json:"language"`
	Deleted          bool   `json:"deleted"`
	BotUserName      string `json:"bot_username"`
	Activated        bool   `json:"activated"`
	Private          bool   `json:"private"`
	UpdateStamp      string `json:"updatestamp"`
	Branch           string `json:"branch"`
	UploadToken      string `json:"upload_token"`
	Active           bool   `json:"active"`
	ServiceID        string `json:"service_id"`
	ImageToken       string `json:"image_token"`
}

func (c *Client) Get(hosting GitHosting, repoOwner, repoName string) (Overview, error) {
	apiURL := fmt.Sprintf(BaseURL+"/%s/%s/%s", hosting.String(), repoOwner, repoName)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return Overview{}, errors.Wrap(err, "failed to create an http request.")
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.Token))

	resp, err := c.HTTPClient.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Overview{}, errors.Wrap(err, "failed to read response body.")
	}

	var o Overview
	err = json.Unmarshal(body, &o)
	if err != nil {
		return Overview{}, errors.Wrap(err, "failed to json-unmarshal the response body.")
	}

	return o, nil
}

func (c *Client) GetBranch(hosting GitHosting, repoOwner, repoName, branchName string) (Overview, error) {
	apiURL := fmt.Sprintf(BaseURL+"/%s/%s/%s/branch/%s", hosting.String(), repoOwner, repoName, branchName)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return Overview{}, errors.Wrap(err, "failed to create an http request.")
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.Token))

	resp, err := c.HTTPClient.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Overview{}, errors.Wrap(err, "failed to read response body.")
	}

	var o Overview
	err = json.Unmarshal(body, &o)
	if err != nil {
		return Overview{}, errors.Wrap(err, "failed to json-unmarshal the response body.")
	}

	return o, nil
}

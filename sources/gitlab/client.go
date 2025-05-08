package gitlab

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"vislab/config"

	"github.com/google/go-querystring/query"
)

const (
	defaultApiPrefix = "/api/v4/"
	defaultRateLimit = 1000 * time.Millisecond
	defaultPerPage   = 40
	defaultTimeout   = 20 * time.Second
)

type Client struct {
	client    *http.Client
	token     string
	baseURL   *url.URL
	rateLimit time.Duration
	perPage   int64

	Projects     *ProjectService
	Tags         *TagService
	Commits      *CommitService
	Branches     *BranchService
	Files        *FIleService
	Contributors *ContributorService
	Groups       *GroupService
	Search       *SearchService
	MergeRequest *MergeRequestService
}

func NewClient(token, baseUrl string, options ...ClientOption) (*Client, error) {
	client := &Client{
		client: &http.Client{Timeout: defaultTimeout, Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // TODO: make optional
		}},
		token:     token,
		rateLimit: defaultRateLimit,
		perPage:   defaultPerPage,
	}

	if err := client.setBaseURL(baseUrl); err != nil {
		return nil, err
	}

	for _, option := range options {
		if err := option(client); err != nil {
			return nil, err
		}
	}

	if !strings.Contains(client.baseURL.Path, "/api/") {
		if err := client.setApiPrefix(defaultApiPrefix); err != nil {
			return nil, err
		}
	}

	client.Projects = &ProjectService{client: client}
	client.Tags = &TagService{client: client}
	client.Commits = &CommitService{client: client}
	client.Branches = &BranchService{client: client}
	client.Files = &FIleService{client: client}
	client.Contributors = &ContributorService{client: client}
	client.Groups = &GroupService{client: client}
	client.Search = &SearchService{client: client}
	client.MergeRequest = &MergeRequestService{client: client}

	return client, nil
}

func (c *Client) setBaseURL(urlStr string) error {
	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	c.baseURL = baseURL

	return nil
}

func (c *Client) setApiPrefix(apiPrefix string) error {
	apiRegexStr := `^/api/v\d/$`

	apiRegex, err := regexp.Compile(apiRegexStr)
	if err != nil {
		return err
	}

	if !apiRegex.MatchString(apiPrefix) {
		return fmt.Errorf("invalid api prefix: %s", apiPrefix)
	}

	c.baseURL.Path += apiPrefix

	return nil
}

func GetOptions(config *config.GitLabClientConfig) []ClientOption {
	options := []ClientOption{}

	if config.GitLabRateLimit != 0 {
		options = append(options, WithRateLimit(config.GitLabRateLimit))
	}
	if config.GitlabPerPage != 0 {
		options = append(options, WithPerPage(config.GitlabPerPage))
	}
	if config.GitlabAPIPrefix != "" {
		options = append(options, WithAPIPrefix(config.GitlabAPIPrefix))
	}
	if config.GitlabTimeout != 0 {
		timeout := time.Duration(config.GitlabTimeout) * time.Millisecond
		options = append(options, WithTimeout(timeout))
	}

	return options
}

func (c *Client) NewRequest(ctx context.Context, method, path string, opt interface{}) (*http.Request, error) {
	u := *c.baseURL
	unescaped, err := url.PathUnescape(path)
	if err != nil {
		return nil, err
	}

	u.RawPath = c.baseURL.Path + path
	u.Path = c.baseURL.Path + unescaped

	reqHeaders := http.Header{}
	reqHeaders.Set("Accept", "application/json")

	if opt != nil {
		q, err := query.Values(opt)
		if err != nil {
			return nil, err
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	for k, v := range reqHeaders {
		req.Header[k] = v
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	if values := req.Header.Values("PRIVATE-TOKEN"); len(values) == 0 {
		req.Header.Set("PRIVATE-TOKEN", c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body)

	response := newResponse(resp)

	if err := CheckResponse(resp); err != nil {
		return response, err
	}

	err = json.NewDecoder(resp.Body).Decode(v)

	return response, err
}

package galleryruntime

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
)

// GitHubConfig holds the credentials required to authenticate as a GitHub App
// installation and file issues on a target repository.
type GitHubConfig struct {
	// AppID is the numeric GitHub App ID.
	AppID int64
	// InstallationID is the installation ID for the target repo/org.
	InstallationID int64
	// PrivateKeyPEM is the PEM-encoded RSA private key for the App.
	PrivateKeyPEM string
	// Repo is the target repository in "owner/repo" format.
	Repo string
}

// GitHubClient posts issues to a GitHub repository as a GitHub App installation.
type GitHubClient struct {
	cfg    GitHubConfig
	client *http.Client
}

// NewGitHubClient creates a GitHubClient using GitHub App installation credentials.
func NewGitHubClient(cfg GitHubConfig) (*GitHubClient, error) {
	if cfg.AppID == 0 || cfg.InstallationID == 0 || cfg.PrivateKeyPEM == "" || cfg.Repo == "" {
		return nil, fmt.Errorf("github: missing required config fields")
	}
	itr, err := ghinstallation.New(
		http.DefaultTransport,
		cfg.AppID,
		cfg.InstallationID,
		[]byte(cfg.PrivateKeyPEM),
	)
	if err != nil {
		return nil, fmt.Errorf("github: create installation transport: %w", err)
	}
	return &GitHubClient{
		cfg:    cfg,
		client: &http.Client{Transport: itr},
	}, nil
}

// issueRequest is the JSON body for POST /repos/{owner}/{repo}/issues.
type issueRequest struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels,omitempty"`
}

// issueResponse is the minimal JSON response from the GitHub Issues API.
type issueResponse struct {
	HTMLURL string `json:"html_url"`
	Number  int    `json:"number"`
}

// CreateIssue opens a new GitHub issue and returns the issue URL.
func (c *GitHubClient) CreateIssue(ctx context.Context, title, body string, labels []string) (string, error) {
	payload, err := json.Marshal(issueRequest{Title: title, Body: body, Labels: labels})
	if err != nil {
		return "", fmt.Errorf("github: marshal issue: %w", err)
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/issues", c.cfg.Repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("github: build request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("github: post issue: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var ghErr struct {
			Message string `json:"message"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&ghErr)
		return "", fmt.Errorf("github: unexpected status %d: %s", resp.StatusCode, ghErr.Message)
	}

	var result issueResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("github: decode response: %w", err)
	}
	return result.HTMLURL, nil
}

// BuildIssueContent constructs the GitHub issue title and markdown body
// aggregating all provided open feedback items for a gallery component.
func BuildIssueContent(comp GalleryComponent, items []Feedback, galleryBaseURL string) (title, body string) {
	n := len(items)
	title = fmt.Sprintf("[Gallery Feedback] %s (%d item", comp.Name, n)
	if n != 1 {
		title += "s"
	}
	title += ")"

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## Component: %s\n\n", comp.Name))
	sb.WriteString(fmt.Sprintf("**Category:** %s", comp.Category))
	if comp.Subcategory != "" {
		sb.WriteString(fmt.Sprintf(" › %s", comp.Subcategory))
	}
	sb.WriteString("\n")
	if comp.Description != "" {
		sb.WriteString(fmt.Sprintf("**Description:** %s\n", comp.Description))
	}
	if galleryBaseURL != "" {
		sb.WriteString(fmt.Sprintf("**Gallery URL:** %s/gallery/%s\n", galleryBaseURL, comp.Slug))
	}
	sb.WriteString("\n---\n\n### Feedback Items\n\n")

	for i, item := range items {
		sb.WriteString(fmt.Sprintf("#### %d. \"%s\"\n\n", i+1, item.Comment))

		var ctx map[string]interface{}
		if item.ContextJSON != "" && item.ContextJSON != "{}" {
			_ = json.Unmarshal([]byte(item.ContextJSON), &ctx)
		}
		if ctx != nil {
			if tag, _ := ctx["tagName"].(string); tag != "" {
				if sel, _ := ctx["selectorPath"].(string); sel != "" {
					sb.WriteString(fmt.Sprintf("- **Element:** `%s` @ `%s`\n", tag, sel))
				} else {
					sb.WriteString(fmt.Sprintf("- **Element:** `%s`\n", tag))
				}
			} else if sel, _ := ctx["selectorPath"].(string); sel != "" {
				sb.WriteString(fmt.Sprintf("- **Selector:** `%s`\n", sel))
			}
			if text, _ := ctx["innerText"].(string); text != "" {
				trimmed := strings.TrimSpace(text)
				if len(trimmed) > 80 {
					trimmed = trimmed[:80] + "…"
				}
				sb.WriteString(fmt.Sprintf("- **Inner text:** %s\n", trimmed))
			}
		}
		sb.WriteString(fmt.Sprintf("- **Submitted:** %s\n\n", item.CreatedAt.Format(time.RFC3339)))
	}

	sb.WriteString(fmt.Sprintf("---\n*Filed automatically by the go-daisy gallery on %s*\n", time.Now().UTC().Format("2006-01-02")))

	return title, sb.String()
}

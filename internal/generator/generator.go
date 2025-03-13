package generator

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/google/go-github/v69/github"
)

func GitHubGenerate(ctx context.Context, workflowRunId int, jobName, owner, repo string) error {
	archivePath := fmt.Sprintf("tw_github_workflow_logs_%d.zip", workflowRunId)
	if err := downloadArchive(ctx, workflowRunId, owner, repo, archivePath); err != nil {
		return fmt.Errorf("github generate download: %w", err)
	}

	zipReader, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("github generate open archive: %w", err)
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			// Only interested in log *files*
			continue
		}

		if !strings.HasPrefix(file.Name, fmt.Sprintf("%s/", jobName)) {
			// Unrelated logs
			continue
		}

		// TODO: Is it safe to rely on this name?
		if file.Name == fmt.Sprintf("%s/1_Set up job.txt", jobName) {
			logs, err := readZipFile(file)
			if err != nil {
				return fmt.Errorf("reading setup job log: %w", err)
			}
			logsReader := strings.NewReader(logs.String())
			actions, err := getActions(logsReader)
			if err != nil {
				return fmt.Errorf("getting actions from logs: %w", err)
			}
			for name, info := range actions {
				fmt.Println("action:", name, "info", info)
			}
		}
	}

	return nil
}

func downloadArchive(ctx context.Context, workflowRunId int, owner, repo, archivePath string) error {
	client := github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN"))

	// Get the log download URL
	url, _, err := client.Actions.GetWorkflowRunLogs(ctx, owner, repo, int64(workflowRunId), 20)
	if err != nil {
		return fmt.Errorf("github get workflow run logs: %w", err)
	}
	fmt.Println(url)

	// Download the log archive
	resp, err := http.Get(url.String())
	if err != nil {
		return fmt.Errorf("github download workflow logs: %w", err)
	}

	// Save the logs to a file
	outFile, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("github create logs zip archive: %w", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, resp.Body); err != nil {
		return fmt.Errorf("github writing logs zip archive: %w", err)
	}

	return nil
}

func readZipFile(zipFile *zip.File) (fmt.Stringer, error) {
	var buffer strings.Builder
	file, err := zipFile.Open()
	if err != nil {
		return &buffer, fmt.Errorf("opening zip file %s: %w", zipFile.Name, err)
	}
	defer file.Close()

	if _, err := io.Copy(&buffer, file); err != nil {
		return &buffer, fmt.Errorf("reading zip file %s: %w", zipFile.Name, err)
	}

	return &buffer, nil
}

type ActionReference struct {
	Name    string
	Version string
	Digest  string
}

var actionRegexp = regexp.MustCompile(`\bDownload action repository '(.+)' \(SHA:([a-f0-9]{40})\)`)

func getActions(log io.Reader) (map[string]ActionReference, error) {
	refs := map[string]ActionReference{}

	scanner := bufio.NewScanner(log)
	for scanner.Scan() {
		line := scanner.Bytes()
		if !actionRegexp.Match(line) {
			continue
		}
		matches := actionRegexp.FindSubmatch(line)
		if len(matches) != 3 {
			// Malformed action line. TODO: return error instead?
			continue
		}
		key := matches[1]

		parts := bytes.SplitN(key, []byte("@"), 2)
		name := parts[0]
		var version []byte
		if len(parts) > 0 {
			version = parts[1]
		}

		digest := matches[2]

		refs[string(key)] = ActionReference{
			Name:    string(name),
			Version: string(version),
			Digest:  string(digest),
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading actions from logs: %w", err)
	}

	return refs, nil
}

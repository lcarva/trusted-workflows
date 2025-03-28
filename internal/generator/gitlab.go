package generator

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"regexp"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func GitLabGenerate(ctx context.Context, pipelineId int, owner, repo string) error {
	client, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"))
	if err != nil {
		return fmt.Errorf("creating client: %w", err)
	}

	projectId := fmt.Sprintf("%s/%s", owner, repo)
	jobs, _, err := client.Jobs.ListPipelineJobs(projectId, pipelineId, nil)
	if err != nil {
		return fmt.Errorf("listing pipeline jobs: %w", err)
	}

	for _, job := range jobs {
		fmt.Println("JOB: ID:", job.ID, " NAME:", job.Name)

		// Get the logs for this job.
		logs, _, err := client.Jobs.GetTraceFile(projectId, job.ID, nil)
		if err != nil {
			return fmt.Errorf("getting trace file: %w", err)
		}

		var inPrepExec bool
		scanner := bufio.NewScanner(logs)
		for scanner.Scan() {
			line := scanner.Bytes()

			subLine := bytes.SplitN(line, []byte{13}, 2)[0]

			if gitlabPrepExecStartRegexp.Match(subLine) {
				inPrepExec = true
				continue
			}

			if gitlabSectionEndRegexp.Match(subLine) {
				if inPrepExec {
					break
				}
				continue
			}

			if inPrepExec {
				fmt.Println("PREPARE_EXECUTOR", string(line))
			}
		}

	}

	return nil
}

var gitlabPrepExecStartRegexp = regexp.MustCompile(`^section_start:\d+:prepare_executor$`)
var gitlabSectionEndRegexp = regexp.MustCompile(`^section_end:.*$`)

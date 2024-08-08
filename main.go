package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v2"
)

// PR represents a pull request with a count label
export GITHUB_TOKEN=ghp_h1n1WHdd0WMdjbodiRRxB3O8GvMIql1coH
type PR struct {
	ID     int
	Branch string
	Count  int
}

// Conf represents the structure of conf.yaml
type Conf struct {
	Schedule string `yaml:"schedule"`
}

// GitRepository interface for Git operations
type GitRepository interface {
	FetchPRs() ([]PR, error)
	FetchPRBranch(pr PR) (*git.Repository, error)
	GetFileContent(repo *git.Repository, filePath string) ([]byte, error)
	GetChangedFiles(pr PR) ([]string, error)
	UpdateCountLabel(pr PR, count int) error
	MergePR(pr PR) error
	UpdateBranch(pr PR) error
	ListLabels(pr PR) ([]string, error)
}

// RealGitRepository is a concrete implementation of GitRepository
type RealGitRepository struct {
	owner string
	repo  string
	token string
}

// NewRealGitRepository constructor
func NewRealGitRepository(owner, repo, token string) *RealGitRepository {
	return &RealGitRepository{
		owner: owner,
		repo:  repo,
		token: token,
	}
}

// FetchPRs fetches all open PRs from the repository
func (r *RealGitRepository) FetchPRs() ([]PR, error) {
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: fmt.Sprintf("https://github.com/%s/%s.git", r.owner, r.repo),
		Auth: &http.BasicAuth{
			Username: "your-username", // replace with your GitHub username
			Password: r.token,
		},
	})
	if err != nil {
		return nil, err
	}

	iter, err := repo.Branches()
	if err != nil {
		return nil, err
	}

	var prs []PR
	iter.ForEach(func(ref *plumbing.Reference) error {
		if strings.HasPrefix(ref.Name().String(), "refs/pull/") {
			prs = append(prs, PR{
				ID:     extractPRID(ref.Name().String()),
				Branch: ref.Name().String(),
				Count:  0,
			})
		}
		return nil
	})

	return prs, nil
}

func extractPRID(ref string) int {
	var id int
	fmt.Sscanf(ref, "refs/pull/%d/head", &id)
	return id
}

// FetchPRBranch fetches the branch of a PR
func (r *RealGitRepository) FetchPRBranch(pr PR) (*git.Repository, error) {
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:           fmt.Sprintf("https://github.com/%s/%s.git", r.owner, r.repo),
		ReferenceName: plumbing.ReferenceName(pr.Branch),
		SingleBranch:  true,
		Auth: &http.BasicAuth{
			Username: "your-username", // replace with your GitHub username
			Password: r.token,
		},
	})
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// GetFileContent retrieves the content of a file from a repository
func (r *RealGitRepository) GetFileContent(repo *git.Repository, filePath string) ([]byte, error) {
	if repo == nil {
		return nil, fmt.Errorf("invalid repository")
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	file, err := commit.File(filePath)
	if err != nil {
		return nil, err
	}

	reader, err := file.Reader()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	return buf.Bytes(), nil
}

// GetChangedFiles retrieves the list of files changed in the most recent commit of a PR
func (r *RealGitRepository) GetChangedFiles(pr PR) ([]string, error) {
	repo, err := r.FetchPRBranch(pr)
	if err != nil {
		return nil, err
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	stats, err := commit.Stats()
	if err != nil {
		return nil, err
	}

	var files []string
	for _, stat := range stats {
		files = append(files, stat.Name)
	}
	return files, nil
}

// ListLabels lists labels for a PR
func (r *RealGitRepository) ListLabels(pr PR) ([]string, error) {
	// For simplicity, we'll assume labels are stored in a map
	// In a real implementation, you would fetch labels from the repository
	labels := map[int][]string{
		pr.ID: {"count:1", "count:2"},
	}
	return labels[pr.ID], nil
}

// UpdateCountLabel updates the count label of a PR
func (r *RealGitRepository) UpdateCountLabel(pr PR, count int) error {
	// For simplicity, we'll print the updated label
	// In a real implementation, you would update the label in the repository
	fmt.Printf("Updated count label for PR %d: count:%d\n", pr.ID, count)
	return nil
}

// MergePR merges a PR
func (r *RealGitRepository) MergePR(pr PR) error {
	fmt.Printf("Merging PR %d\n", pr.ID)
	// For simplicity, we'll print the merge operation
	// In a real implementation, you would merge the PR in the repository
	fmt.Printf("Successfully merged PR %d\n", pr.ID)
	return nil
}

// UpdateBranch updates the branch of a PR
func (r *RealGitRepository) UpdateBranch(pr PR) error {
	// For simplicity, we'll print the update operation
	// In a real implementation, you would update the branch in the repository
	fmt.Printf("Updated branch for PR %d\n", pr.ID)
	return nil
}

// PRProcessor struct to hold dependencies
type PRProcessor struct {
	repo      GitRepository
	batchSize int
}

// NewPRProcessor constructor
func NewPRProcessor(repo GitRepository, batchSize int) *PRProcessor {
	return &PRProcessor{repo: repo, batchSize: batchSize}
}

// ProcessBatch processes a batch of PRs
func (p *PRProcessor) ProcessBatch() error {
	prs, err := p.repo.FetchPRs()
	if err != nil {
		return err
	}

	// Determine the largest existing count across all PRs
	maxCount := 0
	for _, pr := range prs {
		labels, err := p.repo.ListLabels(pr)
		if err != nil {
			return err
		}
		for _, label := range labels {
			if strings.HasPrefix(label, "count:") {
				var currentCount int
				fmt.Sscanf(label, "count:%d", &currentCount)
				if currentCount > maxCount {
					maxCount = currentCount
				}
			}
		}
	}

	// Process only the top batchSize PRs
	for i := 0; i < min(p.batchSize, len(prs)); i++ {
		pr := prs[i]

		changedFiles, err := p.repo.GetChangedFiles(pr)
		if err != nil {
			fmt.Printf("Error getting changed files for PR %d: %v\n", pr.ID, err)
			maxCount++
			p.repo.UpdateCountLabel(pr, maxCount)
			continue
		}

		fmt.Printf("Changed files for PR %d: %v\n", pr.ID, changedFiles) // Log the changed files
		if !allFilesAreConfYaml(changedFiles) {
			fmt.Printf("Not all files are conf.yaml for PR %d\n", pr.ID)
			maxCount++
			p.repo.UpdateCountLabel(pr, maxCount)
			continue
		}

		// Log the changed files and the label to be updated
		fmt.Printf("PR %d changed files: %v\n", pr.ID, changedFiles)
		fmt.Printf("Updating count label for PR %d\n", pr.ID)

		repo, err := p.repo.FetchPRBranch(pr)
		if err != nil {
			fmt.Printf("Error fetching branch for PR %d: %v\n", pr.ID, err)
			maxCount++
			p.repo.UpdateCountLabel(pr, maxCount)
			continue
		}

		// Iterate through the changed files to get the configuration file content
		for _, file := range changedFiles {
			if file == "conf.yaml" || strings.HasSuffix(file, "/conf.yaml") {
				content, err := p.repo.GetFileContent(repo, file)
				if err != nil {
					fmt.Printf("Error getting file content for PR %d: %v\n", pr.ID, err)
					maxCount++
					p.repo.UpdateCountLabel(pr, maxCount)
					continue
				}

				fmt.Printf("File content for PR %d: %s\n", pr.ID, string(content)) // Log the content of conf.yaml

				var conf Conf
				if err := yaml.Unmarshal(content, &conf); err != nil {
					maxCount++
					p.repo.UpdateCountLabel(pr, maxCount)
					continue
				}

				// Interpret and print the next schedule
				nextScheduleTime, err := getNextScheduleTime(conf.Schedule)
				if err != nil {
					maxCount++
					p.repo.UpdateCountLabel(pr, maxCount)
					continue
				}
				fmt.Printf("PR %d next schedule due: %v\n", pr.ID, nextScheduleTime)
				fmt.Printf("Current time: %v\n", time.Now())

				// Detailed logging for the time window check
				fmt.Printf("Checking if current time is within 60 seconds before or after the next schedule, or up to 30 minutes past due.\n")
				fmt.Printf("Current time: %v\n", time.Now())
				fmt.Printf("Next schedule time: %v\n", nextScheduleTime)
				fmt.Printf("30 seconds before next schedule time: %v\n", nextScheduleTime.Add(-60*time.Second))
				fmt.Printf("30 seconds after next schedule time: %v\n", nextScheduleTime.Add(30*time.Second))
				fmt.Printf("30 minutes after next schedule time: %v\n", nextScheduleTime.Add(30*time.Minute))

				// Check if the PR is due for merging (within 60 seconds window or up to 30 minutes past due)
				if (time.Now().After(nextScheduleTime.Add(-30*time.Second)) && time.Now().Before(nextScheduleTime.Add(30*time.Second))) ||
					(time.Now().After(nextScheduleTime) && time.Now().Before(nextScheduleTime.Add(30*time.Minute))) {
					fmt.Printf("Attempting to merge PR %d\n", pr.ID)
					if err := p.repo.MergePR(pr); err != nil {
						fmt.Printf("Failed to merge PR %d: %v\n", pr.ID, err)
						maxCount++
						p.repo.UpdateCountLabel(pr, maxCount)
					} else {
						fmt.Printf("Successfully merged PR %d\n", pr.ID)
					}
				} else {
					fmt.Printf("PR %d is not due for merging.\n", pr.ID)
					maxCount++
					p.repo.UpdateCountLabel(pr, maxCount)
				}
				break // If we find the configuration file, we don't need to check other files
			}
		}
	}

	return nil
}

// isValidCronExpression validates a cron expression
func isValidCronExpression(expr string) bool {
	_, err := cron.ParseStandard(expr)
	return err == nil
}

// allFilesAreConfYaml checks if all changed files are conf.yaml
func allFilesAreConfYaml(files []string) bool {
	for _, file := range files {
		if !(file == "conf.yaml" || strings.HasSuffix(file, "/conf.yaml")) {
			return false
		}
	}
	return true
}

// sortPRsByCount sorts PRs by the count label in ascending order
func sortPRsByCount(prs []PR) []PR {
	sort.SliceStable(prs, func(i, j int) bool {
		return prs[i].Count < prs[j].Count
	})
	return prs
}

// getNextScheduleTime parses the schedule and returns the next scheduled time
func getNextScheduleTime(schedule string) (time.Time, error) {
	sched, err := cron.ParseStandard(schedule)
	if err != nil {
		return time.Time{}, err
	}
	return sched.Next(time.Now()), nil
}

func main() {
	owner := "igboma"
	repo := "cd-pipeline"
	token := os.Getenv("GITHUB_TOKEN") // Make sure to set this environment variable

	repository := NewRealGitRepository(owner, repo, token)
	batchSize := 5

	processor := NewPRProcessor(repository, batchSize)
	if err := processor.ProcessBatch(); err != nil {
		log.Fatalf("Error processing PR batch: %v", err)
	}
}

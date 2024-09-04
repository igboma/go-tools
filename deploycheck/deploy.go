package deploycheck

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gopkg.in/yaml.v2"
)

// Config struct to match the conf.yaml structure
type Config struct {
	Version                        string `yaml:"version"`
	VersionOverride                string `yaml:"versionOverride"`
	Namespace                      string `yaml:"namespace"`
	GomplateDatasources            string `yaml:"gomplateDatasources"`
	HeoRoot                        string `yaml:"heoRoot"`
	HeoRevision                    string `yaml:"heoRevision"`
	HeoRevisionOverride            string `yaml:"heoRevisionOverride"`
	EnableArgoHookDeleteRedis      string `yaml:"enableArgoHookDeleteRedis"`
	EnableArgoHookDeleteRedisForce string `yaml:"enableArgoHookDeleteRedisForce"`
	SlackNotifyChannel             string `yaml:"slackNotifyChannel"`
	Onboarded                      string `yaml:"onboarded"`
	DeploymentSchedule             string `yaml:"deploymentSchedule"`
	DeploymentWindow               int    `yaml:"deploymentWindow"`
}

// GitRepository interface to abstract git repository operations
type GitRepository interface {
	GetFileContentFromBranch(branch, file string) (string, error)
	GetFileContentFromCommit(commitHash, file string) (string, error)
}

// GitRepositoryImpl is the actual implementation using go-git
type GitRepositoryImpl struct {
	repo *git.Repository
}

func NewGitRepository(repo *git.Repository) *GitRepositoryImpl {
	return &GitRepositoryImpl{repo: repo}
}

func (g *GitRepositoryImpl) GetFileContentFromBranch(branch, file string) (string, error) {
	ref, err := g.repo.Reference(plumbing.ReferenceName(branch), true)
	if err != nil {
		return "", err
	}

	commit, err := g.repo.CommitObject(ref.Hash())
	if err != nil {
		return "", err
	}

	return GetFileContentFromCommit(commit, file)
}

func (g *GitRepositoryImpl) GetFileContentFromCommit(commitHash, file string) (string, error) {
	commit, err := g.repo.CommitObject(plumbing.NewHash(commitHash))
	if err != nil {
		return "", err
	}

	return GetFileContentFromCommit(commit, file)
}

// VersionChecker struct that uses GitRepository to check versions and revisions
type VersionChecker struct {
	repo GitRepository
}

func NewVersionChecker(repo GitRepository) *VersionChecker {
	return &VersionChecker{repo: repo}
}

func (v *VersionChecker) GetVersionAndHeoRevision(branch, file string) (string, string, error) {
	content, err := v.repo.GetFileContentFromBranch(branch, file)
	if err != nil {
		return "", "", err
	}

	var config Config
	if err := yaml.Unmarshal([]byte(content), &config); err != nil {
		return "", "", err
	}

	return config.Version, config.HeoRevision, nil
}

func CheckVersionAndHeoRevisionDiff(r *git.Repository, file string) bool {
	// Read current file contents
	currentConfig, err := GetCurrentConfig(file)
	if err != nil {
		log.Fatalf("Failed to get current config: %v", err)
	}

	// Open previous config from origin/main
	gitRepo := NewGitRepository(r)
	previousConfigContent, err := gitRepo.GetFileContentFromBranch("refs/remotes/origin/main", file)
	if err != nil {
		log.Fatalf("Failed to get previous config: %v", err)
	}

	var previousConfig Config
	err = yaml.Unmarshal([]byte(previousConfigContent), &previousConfig)
	if err != nil {
		log.Fatalf("Failed to parse previous YAML: %v", err)
	}

	// Compare version and heoRevision
	if currentConfig.Version != previousConfig.Version || currentConfig.HeoRevision != previousConfig.HeoRevision {
		fmt.Println("Version or heoRevision has changed")
		return true
	}

	fmt.Println("No version or heoRevision change. Skipping creation of deployment.")
	return false
}

// Utility Functions
func GetFileContentFromCommit(commit *object.Commit, file string) (string, error) {
	tree, err := commit.Tree()
	if err != nil {
		return "", err
	}

	entry, err := tree.File(file)
	if err != nil {
		return "", err
	}

	content, err := entry.Contents()
	if err != nil {
		return "", err
	}

	return content, nil
}

func GetCurrentConfig(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func WriteOutput(outputFile, key, value string) {
	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalf("Failed to open GITHUB_OUTPUT: %v", err)
	}
	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%s=%s\n", key, value)); err != nil {
		log.Fatalf("Failed to write output: %v", err)
	}
}

func RemoveVersionAndHeoRevision(config *Config) string {
	config.Version = ""
	config.HeoRevision = ""
	jsonData, _ := json.Marshal(config)
	return string(jsonData)
}

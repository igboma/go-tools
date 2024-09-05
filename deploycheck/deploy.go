package deploycheck

import (
	"encoding/json"
	"fmt"
	"gitpkg/qgit"
	"io/ioutil"

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

// ConfigLoader interface
type ConfigLoader interface {
	LoadConfig(filePath string) (*Config, error)
}

// FileConfigLoader implementation
type FileConfigLoader struct{}

func (f *FileConfigLoader) LoadConfig(filePath string) (*Config, error) {
	data, err := ioutil.ReadFile(filePath)
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

// VersionChecker struct that uses GitRepository to check versions and revisions
type VersionChecker struct {
	repo qgit.GitRepository
}

func NewVersionChecker(repo qgit.GitRepository) *VersionChecker {
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

func CheckVersionAndHeoRevisionDiff(gitRepo qgit.GitRepository, loader ConfigLoader, file string) (bool, error) {
	// Read current file contents
	fmt.Printf("CheckVersionAndHeoRevisionDiff file Path ==> %v", file)
	currentConfig, err := loader.LoadConfig(file)
	if err != nil {
		fmt.Printf("CheckVersionAndHeoRevisionDiff file Path err ==> %v", err)
		return false, err
	}

	// Open previous config from origin/main
	previousConfigContent, err := gitRepo.GetFileContentFromBranch("refs/remotes/origin/main", file)
	if err != nil {
		return false, err
	}

	var previousConfig Config
	err = yaml.Unmarshal([]byte(previousConfigContent), &previousConfig)
	if err != nil {
		return false, err
	}

	// Compare version and heoRevision
	if currentConfig.Version != previousConfig.Version || currentConfig.HeoRevision != previousConfig.HeoRevision {
		fmt.Println("Version or heoRevision has changed")
		return true, nil
	}

	fmt.Println("No version or heoRevision change. Skipping creation of deployment.")
	return false, nil
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

func RemoveVersionAndHeoRevision(config *Config) string {
	config.Version = ""
	config.HeoRevision = ""
	jsonData, _ := json.Marshal(config)
	return string(jsonData)
}

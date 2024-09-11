package deploycheck

import (
	"encoding/json"
	"fmt"
	"gitpkg/qgit"
	"gitpkg/utilities"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config struct to match the conf.yaml structure
type ConfigFile struct {
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
type DeployCheckerOption struct {
	PrNumber          int
	Token             string
	Url               string
	Path              string
	OutputFile        string
	Action            string
	PrMerged          string
	SourceBranch      string
	DestinationBranch string
}

type DeployChecker struct {
	gitClient      *qgit.Client
	outputWriter   *utilities.FileOutputWriter
	option         DeployCheckerOption
	component      string
	environment    string
	needDeployment bool
	version        string
	heoRevision    string
	isRelease      string
}

func (gr *DeployChecker) GetConfFileChangedByPRNumber() ([]string, error) {
	return gr.gitClient.ChangedFiles(gr.option.DestinationBranch, gr.option.SourceBranch)
}

func (gr *DeployChecker) GetComponentConfFileChangedByPRNumber(pr int) (string, error) {
	files, _ := gr.gitClient.ChangedFiles(gr.option.DestinationBranch, gr.option.SourceBranch)
	if len(files) < 1 {
		return "", fmt.Errorf("no files found")
	}
	if len(files) > 1 {
		return "", fmt.Errorf("more than one file was changed")
	}

	file := files[0]
	// Ensure the file is a conf.yaml file
	if !strings.Contains(file, "components/") || !strings.HasSuffix(file, "conf.yaml") {
		return "", fmt.Errorf("the file is not a conf.yaml file")
	}

	return file, nil
}

func (gr *DeployChecker) getConfigData(file, ref string) (configData *ConfigFile, err error) {
	configContent, err := gr.gitClient.FileContentFromBranch(ref, file)
	if err != nil {
		return
	}

	err = yaml.Unmarshal([]byte(configContent), &configData)
	if err != nil {
		return
	}
	return
}

func (gr *DeployChecker) RemoveVersionAndHeoRevision(config *ConfigFile) string {
	config.Version = ""
	config.HeoRevision = ""
	jsonData, _ := json.Marshal(config)
	return string(jsonData)
}

func (gr *DeployChecker) GetSourceAndDestimationConf(file string, prNumber int, destimationBranch string) (currentConfig *ConfigFile, previousConfig *ConfigFile, err error) {
	prRef := fmt.Sprintf("refs/pull/%d/head", prNumber)
	currentConfig, err = gr.getConfigData(file, prRef)
	if err != nil {
		return
	}
	destRef := fmt.Sprintf("refs/remotes/origin/%v", destimationBranch)
	previousConfig, err = gr.getConfigData(file, destRef)
	if err != nil {
		return
	}
	return
}

func (gr *DeployChecker) WriteOutput(key, value string) error {
	return gr.outputWriter.WriteOutput(key, value)
}

func (gr *DeployChecker) Run() error {
	file, err := gr.GetComponentConfFileChangedByPRNumber(gr.option.PrNumber)
	if err != nil {
		return fmt.Errorf("error getting conf file %w", err)
	}
	if len(strings.Split(file, "/")) > 2 {
		gr.component = strings.Split(file, "/")[1]
		gr.environment = strings.Split(file, "/")[2]
	} else {
		return fmt.Errorf("invalid config file")
	}

	if gr.option.Action == "closed" && gr.option.PrMerged == "true" {
		//fmt.Println("PR is merged...")
		configData, err := gr.getConfigData(file, "refs/heads/main")
		fmt.Printf("configData: %v\n", configData)
		if err != nil {
			return fmt.Errorf("failed to get version and heoRevision: %w", err)
		}
		gr.version = configData.Version
		gr.heoRevision = configData.HeoRevision
	} else {
		source, destination, err := gr.GetSourceAndDestimationConf(file, gr.option.PrNumber, "main")
		if err != nil {
			return fmt.Errorf("error checking version and heoRevision: %w", err)
		}
		fmt.Printf("\nsource: %v\n", source.Version)
		fmt.Printf("destination: %v\n", destination.Version)

		gr.needDeployment = source.Version != destination.Version || source.HeoRevision != destination.HeoRevision

		gr.version = source.Version
		gr.heoRevision = source.HeoRevision

		// Compare non-version and non-heoRevision fields
		jsonCurrentOtherFields := gr.RemoveVersionAndHeoRevision(source)
		jsonPreviousOtherFields := gr.RemoveVersionAndHeoRevision(destination)

		fmt.Printf("+version: %s\n", gr.version)
		fmt.Printf("jsonCurrentOtherFields: %s\n", jsonCurrentOtherFields)
		fmt.Printf("jsonPreviousOtherFields: %s\n", jsonPreviousOtherFields)
	}

	// Determine if it is a release version
	gr.isRelease = "true"
	if strings.Contains(gr.version, "-") {
		gr.isRelease = "false"
	}

	gr.outputWriter.WriteOutput("COMPONENT", gr.component)
	gr.outputWriter.WriteOutput("ENVIRONMENT", gr.environment)
	gr.outputWriter.WriteOutput("VERSION", gr.version)
	gr.outputWriter.WriteOutput("IS_RELEASE", gr.isRelease)
	gr.outputWriter.WriteOutput("HEO_REVISION", gr.heoRevision)
	gr.outputWriter.WriteOutput("DEPLOYMENT_NEEDED", fmt.Sprintf("%t", gr.needDeployment))

	fmt.Printf("COMPONENT=%s\n", gr.component)
	fmt.Printf("ENVIRONMENT=%s\n", gr.environment)
	fmt.Printf("VERSION=%s\n", gr.version)
	fmt.Printf("IS_RELEASE=%s\n", gr.isRelease)
	fmt.Printf("HEO_REVISION=%s\n", gr.heoRevision)
	fmt.Printf("DEPLOYMENT_NEEDED=%s\n", fmt.Sprintf("%t", gr.needDeployment))

	return nil
}

func NewDeployChecker(opt DeployCheckerOption) *DeployChecker {
	client, err := qgit.NewClient(
		qgit.WithRepoPath(opt.Path),
		qgit.WithRepoUrl(opt.Url),
	)
	if err != nil {
		return nil
	}
	outputWriter := utilities.NewFileOutputWriter(opt.OutputFile)

	checker := &DeployChecker{gitClient: client,
		outputWriter: outputWriter,
		option:       opt}

	return checker
}

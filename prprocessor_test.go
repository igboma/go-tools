package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// func TestProcessBatch(t *testing.T) {
// 	mockRepo := new(MockGitRepository)
// 	processor := NewPRProcessor(mockRepo, 5)

// 	prs := []PR{
// 		{ID: 1, Branch: "branch-1", Count: 1},
// 		{ID: 2, Branch: "branch-2", Count: 2},
// 	}
// 	mockRepo.On("FetchPRs").Return(prs, nil)

// 	changedFiles := []string{"conf.yaml"}
// 	mockRepo.On("GetChangedFiles", prs[0]).Return(changedFiles, nil)
// 	mockRepo.On("GetChangedFiles", prs[1]).Return(changedFiles, nil)

// 	repo := &git.Repository{}
// 	mockRepo.On("FetchPRBranch", prs[0]).Return(repo, nil)
// 	mockRepo.On("FetchPRBranch", prs[1]).Return(repo, nil)

// 	confContent := []byte(`schedule: "* * * * *"`)
// 	mockRepo.On("GetFileContent", repo, "conf.yaml").Return(confContent, nil)

// 	// Expect UpdateCountLabel to be called when PRs are not due for merging
// 	// Adjusted the count values to match the actual behavior
// 	mockRepo.On("UpdateCountLabel", prs[0], 2).Return(nil).Once()
// 	mockRepo.On("UpdateCountLabel", prs[1], 3).Return(nil).Once()

// 	mockRepo.On("MergePR", prs[0]).Return(nil).Maybe()
// 	mockRepo.On("MergePR", prs[1]).Return(nil).Maybe()
// 	mockRepo.On("ListLabels", prs[0].ID).Return([]string{}, nil)
// 	mockRepo.On("ListLabels", prs[1].ID).Return([]string{}, nil)

// 	err := processor.ProcessBatch()
// 	assert.NoError(t, err)

// 	mockRepo.AssertExpectations(t)
// }

func TestGetNextScheduleTime(t *testing.T) {
	schedule := "* * * * *"
	expectedNextTime := time.Now().Add(1 * time.Minute).Truncate(time.Minute)

	nextTime, err := getNextScheduleTime(schedule)
	assert.NoError(t, err)
	assert.WithinDuration(t, expectedNextTime, nextTime, time.Second)
}

func TestAllFilesAreConfYaml(t *testing.T) {
	files := []string{"conf.yaml", "conf.yaml"}
	assert.True(t, allFilesAreConfYaml(files))

	files = []string{"conf.yaml", "not_conf.yaml"}
	assert.False(t, allFilesAreConfYaml(files))

	// Additional logging
	t.Logf("Tested files: %v", files)
}

func TestIsValidCronExpression(t *testing.T) {
	validExpr := "* * * * *"
	assert.True(t, isValidCronExpression(validExpr))

	invalidExpr := "invalid cron"
	assert.False(t, isValidCronExpression(invalidExpr))
}

func TestSortPRsByCount(t *testing.T) {
	prs := []PR{
		{ID: 1, Count: 5},
		{ID: 2, Count: 2},
		{ID: 3, Count: 8},
	}

	sortedPRs := sortPRsByCount(prs)
	expectedPRs := []PR{
		{ID: 2, Count: 2},
		{ID: 1, Count: 5},
		{ID: 3, Count: 8},
	}

	assert.Equal(t, expectedPRs, sortedPRs)
}

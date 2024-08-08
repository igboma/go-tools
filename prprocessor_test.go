package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestProcessBatch tests the ProcessBatch method
// func TestProcessBatch(t *testing.T) {
// 	mockRepo := new(MockGitRepository)
// 	processor := NewPRProcessor(mockRepo, 5)

// 	prs := []PR{
// 		{ID: 1, Branch: "branch-1", Count: 0},
// 		{ID: 2, Branch: "branch-2", Count: 0},
// 	}

// 	mockRepo.On("FetchPRs").Return(prs, nil)
// 	mockRepo.On("ListLabels", 1).Return([]string{"count:1"}, nil)
// 	mockRepo.On("ListLabels", 2).Return([]string{"count:2"}, nil)
// 	mockRepo.On("GetChangedFiles", prs[0]).Return([]string{"conf.yaml"}, nil)
// 	mockRepo.On("GetChangedFiles", prs[1]).Return([]string{"conf.yaml"}, nil)

// 	repo := &git.Repository{}
// 	mockRepo.On("FetchPRBranch", prs[0]).Return(repo, nil)
// 	mockRepo.On("FetchPRBranch", prs[1]).Return(repo, nil)

// 	mockRepo.On("GetFileContent", repo, "conf.yaml").Return([]byte(`schedule: "* * * * *"`), nil).Twice()
// 	mockRepo.On("UpdateCountLabel", prs[0], 1).Return(nil)
// 	mockRepo.On("UpdateCountLabel", prs[1], 2).Return(nil)

// 	mockRepo.On("MergePR", prs[0]).Return(nil)
// 	mockRepo.On("MergePR", prs[1]).Return(nil)

// 	err := processor.ProcessBatch()
// 	assert.NoError(t, err)

// 	mockRepo.AssertExpectations(t)
// }

// TestGetNextScheduleTime tests the getNextScheduleTime function
func TestGetNextScheduleTime(t *testing.T) {
	schedule := "* * * * *"
	expectedTime := time.Now().Add(time.Minute).Truncate(time.Minute)

	nextTime, err := getNextScheduleTime(schedule)
	assert.NoError(t, err)
	assert.True(t, nextTime.After(expectedTime) || nextTime.Equal(expectedTime))
}

// TestAllFilesAreConfYaml tests the allFilesAreConfYaml function
func TestAllFilesAreConfYaml(t *testing.T) {
	files := []string{"conf.yaml", "folder/conf.yaml"}
	assert.True(t, allFilesAreConfYaml(files))

	files = []string{"conf.yaml", "not_conf.yaml"}
	assert.False(t, allFilesAreConfYaml(files))

	files = []string{"folder/not_conf.yaml"}
	assert.False(t, allFilesAreConfYaml(files))
}

// TestIsValidCronExpression tests the isValidCronExpression function
func TestIsValidCronExpression(t *testing.T) {
	assert.True(t, isValidCronExpression("* * * * *"))
	assert.False(t, isValidCronExpression("invalid-cron"))
}

// TestSortPRsByCount tests the sortPRsByCount function
func TestSortPRsByCount(t *testing.T) {
	prs := []PR{
		{ID: 1, Count: 2},
		{ID: 2, Count: 1},
	}
	sortedPRs := sortPRsByCount(prs)
	assert.Equal(t, 2, sortedPRs[0].ID)
	assert.Equal(t, 1, sortedPRs[1].ID)
}

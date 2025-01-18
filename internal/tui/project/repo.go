package project

import (
	"fmt"
	tBox "yi/internal/tui/box/base"
	fuxo "yi/pkg/repo"
	t "yi/pkg/types"
)

type RepoUpdateJob struct {
	repo *fuxo.RepoConfig
}

func NewRepoUpdateJob(repo *fuxo.RepoConfig) *RepoUpdateJob {
	return &RepoUpdateJob{repo: repo}
}

func (job *RepoUpdateJob) Run() error {
	m := tBox.JobsModel{Jobs: []t.Job{*t.NewJob(fmt.Sprintf("Update %s Cache", job.repo.Name()), job.repo.GetIndex().UpdateCache)}}
	err := tBox.NewWaitingBox(m.Run).Run()

	return err
}

func NewRepoUpdateJobFromRepo(repo fuxo.RepoConfig) *RepoUpdateJob {
	return &RepoUpdateJob{repo: &repo}
}

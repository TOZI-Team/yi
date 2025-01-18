package project

import (
	"fmt"
	tBox "yi/internal/tui/box/base"
	fuxo "yi/pkg/repo"
	"yi/pkg/repo/cache"
	t "yi/pkg/types"
)

type PackageDownloadTUI struct {
	meta []fuxo.SimplePackageMeta
}

func NewPackageDownloadTUI(metas []fuxo.SimplePackageMeta) *PackageDownloadTUI {
	return &PackageDownloadTUI{metas}
}

//type DepDownloadJob struct {
//	meta fuxo.SimplePackageMeta
//	c    chan t.WaitingMessage
//}
//
//func (job DepDownloadJob) Run() string {
//	c <- t.WaitingMessage{t.Waiting}
//}
//
//func NewDepDownloadJob(meta fuxo.SimplePackageMeta, c chan t.WaitingMessage) *DepDownloadJob {
//	return &DepDownloadJob{meta, c}
//}

func (p *PackageDownloadTUI) Run() error {
	var jobs []t.Job
	for _, meta := range p.meta {
		jobs = append(jobs, *t.NewJob(fmt.Sprintf("Download %s", meta.Name), meta.Download))
	}
	m := tBox.JobsModel{Jobs: jobs}

	err := tBox.NewWaitingBox(m.Run).Run()
	if err != nil {
		return err
	}

	return nil
}

type DepMakeTUI struct {
	meta []fuxo.SimplePackageMeta
}

func NewDepMakeTUI(metas []fuxo.SimplePackageMeta) *DepMakeTUI {
	return &DepMakeTUI{metas}
}

func (p *DepMakeTUI) Run() error {
	c := cache.NewConfig(fuxo.GlobalConfig().GetDefault().Name())
	var jobs []t.Job
	for _, meta := range p.meta {
		jobs = append(jobs, *t.NewJob(fmt.Sprintf("Generate Config: %s", meta.Name), func() error {
			err := c.GenerateCache(meta.Name, meta.Ver)
			if err != nil {
				return err
			}
			return nil
		}))
	}
	m := tBox.JobsModel{Jobs: jobs}
	err := tBox.NewWaitingBox(m.Run).Run()
	if err != nil {
		return err
	}

	return nil
}

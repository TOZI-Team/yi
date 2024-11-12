package types

type JobStatue = int8

const (
	Waiting JobStatue = iota
	Success JobStatue = iota
	Err     JobStatue = iota
)

type WaitingMessage struct {
	Statue  JobStatue
	Message string
}

type Job struct {
	Name string
	f    func() error
}

func (j Job) Run() error {
	return j.f()
}

func NewJob(name string, f func() error) *Job {
	return &Job{
		Name: name,
		f:    f,
	}
}

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

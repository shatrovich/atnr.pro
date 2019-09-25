package core

type Repository interface {
	GetTaskDataByID(uint64) (TaskData, error)
	GetLastHashByID(uint64) (string, error)
	PushNewTask(data PushTaskData) (uint64, error)
	PushNewHashTask(uint64, string, int) error
}

type TaskData struct {
	ID           uint64
	Payload      string
	Hash         string
	RoundsCount  int
	CurrentRound int
}

type PushTaskData struct {
	Payload     string
	RoundsCount int
}

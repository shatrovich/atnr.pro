package storage

import (
	"database/sql"

	"github.com/shatrovich/atnr.pro/core"
	"go.uber.org/zap"
)

type Store struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewStore(db *sql.DB, logger *zap.Logger) core.Repository {
	var s core.Repository = &Store{db: db, logger: logger}

	return s
}

func (s *Store) GetTaskDataByID(id uint64) (td core.TaskData, err error) {
	if err = s.db.QueryRow("SELECT T.id, T.rounds, T.payload, TH.round as current_round, TH.hash FROM tasks T INNER JOIN (SELECT * FROM tasks_hashes) TH ON TH.task_id = T.id WHERE T.id = $1 ORDER BY TH.round DESC LIMIT 1", id).Scan(&td.ID, &td.RoundsCount, &td.Payload, &td.CurrentRound, &td.Hash); err != nil {
		s.logger.Error("error GetTaskDataByID", zap.Error(err), zap.Uint64("id", id))

		return td, err
	}

	return td, err
}

func (s *Store) GetLastHashByID(id uint64) (h string, err error) {
	return h, err
}

func (s *Store) PushNewTask(data core.PushTaskData) (id uint64, err error) {
	if err = s.db.QueryRow("INSERT INTO tasks (rounds, payload) VALUES ($1, $2) RETURNING id", data.RoundsCount, data.Payload).Scan(&id); err != nil {
		s.logger.Error("error PushNewTask", zap.Error(err), zap.Int("rounds", data.RoundsCount), zap.String("payload", data.Payload))

		return id, err
	}

	return id, err
}

func (s *Store) PushNewHashTask(taskID uint64, h string, round int) error {
	if _, err := s.db.Exec("INSERT INTO tasks_hashes (task_id, round, hash) VALUES ($1, $2, $3)", taskID, round, h); err != nil {
		s.logger.Error("error PushNewHashTask", zap.Error(err), zap.Uint64("task_id", taskID), zap.String("hash", h), zap.Int("round", round))

		return err
	}

	return nil
}

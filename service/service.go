package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/shatrovich/atnr.pro/core"
	"go.uber.org/zap"
)

type Service struct {
	storage core.Repository
	logger  *zap.Logger
}

type TaskData struct {
	ID          uint64 `json:"id"`
	Payload     string `json:"payload"`
	RoundsCount int    `json:"hash_rounds_cnt"`
	Status      string `json:"status"`
	Hash        string `json:"hash"`
}

type NewTaskData struct {
	ID          uint64 // optional param
	Payload     string
	RoundsCount int
}

func NewService(storage core.Repository, logger *zap.Logger) *Service {
	s := new(Service)
	s.storage = storage
	s.logger = logger

	return s
}

func (s *Service) GetTaskData(id uint64) (b []byte, err error) {
	td, err := s.storage.GetTaskDataByID(id)

	if err != nil {
		s.logger.Error("error GetTaskData", zap.Error(err), zap.Uint64("id", id))

		return b, err
	}

	d := TaskData{
		ID:          td.ID,
		Payload:     td.Payload,
		RoundsCount: td.RoundsCount,
		Hash:        td.Hash,
	}

	if td.RoundsCount == td.CurrentRound {
		d.Status = "finished"
	} else {
		d.Status = "in progress"
	}

	b, err = json.Marshal(&d)

	return b, err
}

func (s *Service) PushNewTask(data NewTaskData) (id uint64, err error) {
	if id, err = s.storage.PushNewTask(core.PushTaskData{Payload: data.Payload, RoundsCount: data.RoundsCount}); err != nil {
		return id, err
	}

	data.ID = id

	go func(s *Service, t *NewTaskData) {
		var hx string

		h := sha256.New()

		h.Write([]byte(t.Payload))

		hx = hex.EncodeToString(h.Sum(nil))

		_ = s.storage.PushNewHashTask(t.ID, hx, 1)

		for i := 1; i < t.RoundsCount; i++ {
			h := sha256.New()

			h.Write([]byte(hx))

			hx = hex.EncodeToString(h.Sum(nil))

			_ = s.storage.PushNewHashTask(t.ID, hx, (i + 1))
		}
	}(s, &data)

	return id, err
}

func (s *Service) GetTaskDataHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		s.logger.Info("error parse id to int", zap.Error(err))

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	b, err := s.GetTaskData(id)

	if err != nil {
		s.logger.Info("error GetTaskData result", zap.Error(err))

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (s *Service) PushTaskHTTP(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Payload     string `json:"payload"`
		RoundsCount int    `json:"hash_rounds_cnt"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		s.logger.Info("error decode body", zap.Error(err))

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	s.logger.Info("push track data", zap.String("payload", data.Payload), zap.Int("rounds_count", data.RoundsCount))

	if data.RoundsCount < 1 {
		http.Error(w, "rounds count must be greater 1", http.StatusBadRequest)

		return
	}

	id, err := s.PushNewTask(NewTaskData{Payload: data.Payload, RoundsCount: data.RoundsCount})

	if err != nil {
		s.logger.Error("error push new track", zap.Error(err))

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	var reply struct {
		ID string `json:"id"`
	}

	reply.ID = fmt.Sprint(id)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(&reply); err != nil {
		s.logger.Error("error encode response body", zap.Error(err))
	}
}

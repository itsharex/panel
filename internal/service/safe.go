package service

import (
	"net/http"

	"github.com/go-rat/chix"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
)

type SafeService struct {
	safeRepo biz.SafeRepo
}

func NewSafeService() *SafeService {
	return &SafeService{
		safeRepo: data.NewSafeRepo(),
	}
}

func (s *SafeService) GetSSH(w http.ResponseWriter, r *http.Request) {
	port, status, err := s.safeRepo.GetSSH()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, chix.M{
		"port":   port,
		"status": status,
	})
}

func (s *SafeService) UpdateSSH(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.SafeUpdateSSH](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.safeRepo.UpdateSSH(req.Port, req.Status); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *SafeService) GetPingStatus(w http.ResponseWriter, r *http.Request) {
	status, err := s.safeRepo.GetPingStatus()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, status)
}

func (s *SafeService) UpdatePingStatus(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.SafeUpdatePingStatus](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.safeRepo.UpdatePingStatus(req.Status); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}
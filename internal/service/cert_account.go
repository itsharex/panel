package service

import (
	"net/http"

	"github.com/go-rat/chix"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
)

type CertAccountService struct {
	certAccountRepo biz.CertAccountRepo
}

func NewCertAccountService() *CertAccountService {
	return &CertAccountService{
		certAccountRepo: data.NewCertAccountRepo(),
	}
}

func (s *CertAccountService) List(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.Paginate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	certDNS, total, err := s.certAccountRepo.List(req.Page, req.Limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, chix.M{
		"total": total,
		"items": certDNS,
	})
}

func (s *CertAccountService) Create(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.CertAccountCreate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	account, err := s.certAccountRepo.Create(req)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, account)
}

func (s *CertAccountService) Update(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.CertAccountUpdate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.certAccountRepo.Update(req); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *CertAccountService) Get(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	account, err := s.certAccountRepo.Get(req.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, account)
}

func (s *CertAccountService) Delete(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.certAccountRepo.Delete(req.ID); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}
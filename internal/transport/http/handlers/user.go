package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Jeskay/micsvc/internal/dto"
	"github.com/Jeskay/micsvc/internal/user"
)

type UserHandler struct {
	svc *user.Service
}

func NewUserHandler(svc *user.Service) *UserHandler {
	return &UserHandler{svc}
}

func (h *UserHandler) Add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var user dto.User
		if err := decoder.Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}
		if err := h.svc.Add(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *UserHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var user dto.User
		if err := decoder.Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}
		if err := h.svc.Update(int32(id), &user); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h *UserHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := h.svc.Delete(int32(id)); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h *UserHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := h.svc.GetAll()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(users); err != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

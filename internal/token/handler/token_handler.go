package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"seno-blackdragon/internal/token/service"
)

type TokenHandler struct {
	svc service.TokenService
}

func New(svc service.TokenService) *TokenHandler {
	return &TokenHandler{svc: svc}
}

func (h *TokenHandler) Create(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.FormValue("user_id")
	value := r.FormValue("value")
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)
	token, err := h.svc.Create(r.Context(), userID, value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(token)
}

func (h *TokenHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	token, err := h.svc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(token)
}

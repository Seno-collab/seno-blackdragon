package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"seno-blackdragon/internal/wallet/service"
)

type WalletHandler struct {
	svc service.WalletService
}

func New(svc service.WalletService) *WalletHandler {
	return &WalletHandler{svc: svc}
}

func (h *WalletHandler) Create(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.FormValue("user_id")
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)
	wallet, err := h.svc.Create(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(wallet)
}

func (h *WalletHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	wallet, err := h.svc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(wallet)
}

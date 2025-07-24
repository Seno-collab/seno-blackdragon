package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"seno-blackdragon/internal/dragon/service"
)

type DragonHandler struct {
	svc service.DragonService
}

func New(svc service.DragonService) *DragonHandler {
	return &DragonHandler{svc: svc}
}

func (h *DragonHandler) Create(w http.ResponseWriter, r *http.Request) {
	ownerIDStr := r.FormValue("owner_id")
	name := r.FormValue("name")
	ownerID, _ := strconv.ParseInt(ownerIDStr, 10, 64)
	dragon, err := h.svc.Create(r.Context(), ownerID, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(dragon)
}

func (h *DragonHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	dragon, err := h.svc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(dragon)
}

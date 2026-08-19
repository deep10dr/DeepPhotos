package handler

import (
	"encoding/json"
	"net/http"

	"deepphotos/backend/internal/model"
	"deepphotos/backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LockedHandler struct {
	lockedRepo *repository.LockedRepository
}

func NewLockedHandler(lockedRepo *repository.LockedRepository) *LockedHandler {
	return &LockedHandler{lockedRepo: lockedRepo}
}

func (h *LockedHandler) List(w http.ResponseWriter, r *http.Request) {
	folders, err := h.lockedRepo.ListAll()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folders)
}

func (h *LockedHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateLockedFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" || req.Passcode == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Folder name and 4-digit passcode are required"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Passcode), bcrypt.DefaultCost)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Failed to process passcode"})
		return
	}

	folder := &model.LockedFolder{
		ID:           "lck_" + uuid.New().String()[:8],
		UserID:       "usr_admin_1",
		Name:         req.Name,
		Description:  req.Description,
		PasscodeHash: string(hashed),
	}

	if err := h.lockedRepo.Create(folder); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(folder)
}

func (h *LockedHandler) VerifyPasscode(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	folder, err := h.lockedRepo.FindByID(id)
	if err != nil || folder == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Locked folder not found"})
		return
	}

	var req model.VerifyPasscodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Invalid JSON request"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(folder.PasscodeHash), []byte(req.Passcode)); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Incorrect 4-digit passcode"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.SuccessResponse{Message: "Passcode verified successfully"})
}

func (h *LockedHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.lockedRepo.Delete(id); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.SuccessResponse{Message: "Locked folder deleted successfully"})
}

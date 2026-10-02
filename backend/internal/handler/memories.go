package handler

import (
	"encoding/json"
	"net/http"

	"deepphotos/backend/internal/model"
	"deepphotos/backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type MemoriesHandler struct {
	memoryRepo *repository.MemoryRepository
}

func NewMemoriesHandler(memoryRepo *repository.MemoryRepository) *MemoriesHandler {
	return &MemoriesHandler{memoryRepo: memoryRepo}
}

func (h *MemoriesHandler) List(w http.ResponseWriter, r *http.Request) {
	memories, err := h.memoryRepo.ListAll()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(memories)
}

func (h *MemoriesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateMemoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Invalid JSON body"})
		return
	}

	memory := &model.Memory{
		ID:          "mem_" + uuid.New().String()[:8],
		Title:       req.Title,
		Description: req.Description,
	}

	if len(req.PhotoIDs) > 0 {
		memory.CoverPhotoID = req.PhotoIDs[0]
	}

	if err := h.memoryRepo.Create(memory); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	if len(req.PhotoIDs) > 0 {
		h.memoryRepo.AddPhotos(memory.ID, req.PhotoIDs)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(memory)
}

func (h *MemoriesHandler) GetDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	memory, err := h.memoryRepo.FindByID(id)
	if err != nil || memory == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Memory not found"})
		return
	}

	photos, err := h.memoryRepo.GetPhotosForMemory(id)
	if err == nil {
		memory.Photos = photos
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(memory)
}

func (h *MemoriesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.memoryRepo.Delete(id); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.SuccessResponse{Message: "Memory deleted successfully"})
}

func (h *MemoriesHandler) AddPhotos(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req model.AddPhotosToMemoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Invalid JSON body"})
		return
	}

	if err := h.memoryRepo.AddPhotos(id, req.PhotoIDs); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.SuccessResponse{Message: "Photos added to memory"})
}

func (h *MemoriesHandler) RemovePhoto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	photoID := chi.URLParam(r, "photoId")

	if err := h.memoryRepo.RemovePhotos(id, []string{photoID}); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.SuccessResponse{Message: "Photo removed from memory"})
}

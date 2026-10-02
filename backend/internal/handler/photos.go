package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"deepphotos/backend/internal/model"
	"deepphotos/backend/internal/repository"
	"deepphotos/backend/internal/service"
	"deepphotos/backend/internal/storage"
	"github.com/go-chi/chi/v5"
)

type PhotosHandler struct {
	photoRepo     *repository.PhotoRepository
	photoService  *service.PhotoService
	storageClient *storage.StorageClient
}

func NewPhotosHandler(photoRepo *repository.PhotoRepository, photoService *service.PhotoService, storageClient *storage.StorageClient) *PhotosHandler {
	return &PhotosHandler{
		photoRepo:     photoRepo,
		photoService:  photoService,
		storageClient: storageClient,
	}
}

// Upload handles multipart photo/video/document upload
func (h *PhotosHandler) Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(64 << 20) // 64MB max memory
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Failed to parse multipart form"})
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "No files uploaded under field 'files'"})
		return
	}

	var uploadedPhotos []*model.Photo
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}

		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		photo, err := h.photoService.UploadPhoto(r.Context(), fileHeader.Filename, contentType, file, fileHeader.Size)
		file.Close()

		if err == nil && photo != nil {
			uploadedPhotos = append(uploadedPhotos, photo)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(uploadedPhotos)
}

// UploadURL handles downloading & saving external website image URLs
func (h *PhotosHandler) UploadURL(w http.ResponseWriter, r *http.Request) {
	var req model.URLUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Invalid URL upload payload"})
		return
	}

	photo, err := h.photoService.UploadPhotoFromURL(r.Context(), req.URL)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	photo.URL = "/api/media/" + photo.ID + "/file"
	photo.ThumbnailURL = "/api/media/" + photo.ID + "/thumbnail"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(photo)
}

// List handles listing media with filters
func (h *PhotosHandler) List(w http.ResponseWriter, r *http.Request) {
	favOnly := r.URL.Query().Get("favorite") == "true"
	delOnly := r.URL.Query().Get("deleted") == "true"
	search := r.URL.Query().Get("search")
	fileType := r.URL.Query().Get("type")
	lockedFolderID := r.URL.Query().Get("locked_folder_id")

	photos, err := h.photoRepo.List(favOnly, delOnly, search, fileType, lockedFolderID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	for i := range photos {
		photos[i].URL = "/api/media/" + photos[i].ID + "/file"
		photos[i].ThumbnailURL = "/api/media/" + photos[i].ID + "/thumbnail"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(photos)
}

// ListGallery handles GET /api/gallery (e.g. /api/gallery?type=photos&deleted=false, /api/gallery?type=video, /api/gallery?type=all)
func (h *PhotosHandler) ListGallery(w http.ResponseWriter, r *http.Request) {
	fileType := r.URL.Query().Get("type")
	if fileType == "" {
		fileType = "gallery"
	}
	favOnly := r.URL.Query().Get("favorite") == "true"
	delOnly := r.URL.Query().Get("deleted") == "true"
	search := r.URL.Query().Get("search")
	lockedFolderID := r.URL.Query().Get("locked_folder_id")

	photos, err := h.photoRepo.List(favOnly, delOnly, search, fileType, lockedFolderID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	for i := range photos {
		photos[i].URL = "/api/media/" + photos[i].ID + "/file"
		photos[i].ThumbnailURL = "/api/media/" + photos[i].ID + "/thumbnail"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(photos)
}

// ListDocuments handles GET /api/documents (e.g. /api/documents?deleted=false)
func (h *PhotosHandler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	favOnly := r.URL.Query().Get("favorite") == "true"
	delOnly := r.URL.Query().Get("deleted") == "true"
	search := r.URL.Query().Get("search")

	photos, err := h.photoRepo.List(favOnly, delOnly, search, "document", "")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	for i := range photos {
		photos[i].URL = "/api/media/" + photos[i].ID + "/file"
		photos[i].ThumbnailURL = "/api/media/" + photos[i].ID + "/thumbnail"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(photos)
}

// GetDetail handles fetching single media metadata & technical EXIF details
func (h *PhotosHandler) GetDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	photo, err := h.photoRepo.FindByID(id)
	if err != nil || photo == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Media item not found"})
		return
	}

	photo.URL = "/api/media/" + photo.ID + "/file"
	photo.ThumbnailURL = "/api/media/" + photo.ID + "/thumbnail"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(photo)
}

// Update handles PUT updates (title, is_favorite, is_deleted, locked_folder_id)
func (h *PhotosHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	photo, err := h.photoRepo.FindByID(id)
	if err != nil || photo == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Media item not found"})
		return
	}

	var req model.UpdatePhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Invalid JSON request body"})
		return
	}

	if req.Title != "" {
		photo.Title = req.Title
	}
	if req.IsFavorite != nil {
		photo.IsFavorite = *req.IsFavorite
	}
	if req.IsDeleted != nil {
		photo.IsDeleted = *req.IsDeleted
	}
	if req.LockedFolderID != nil {
		photo.LockedFolderID = *req.LockedFolderID
	}

	if err := h.photoRepo.Update(photo); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(photo)
}

func (h *PhotosHandler) DeleteSingle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.photoService.DeleteSinglePhoto(r.Context(), id); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.SuccessResponse{Message: "Media deleted successfully"})
}

func (h *PhotosHandler) DeleteBatch(w http.ResponseWriter, r *http.Request) {
	var req model.BatchIDsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Invalid JSON payload"})
		return
	}

	count, err := h.photoService.DeleteBatchPhotos(r.Context(), req.IDs)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.SuccessResponse{Message: strconv.Itoa(count) + " items deleted successfully"})
}

func (h *PhotosHandler) RestoreBatch(w http.ResponseWriter, r *http.Request) {
	var req model.BatchIDsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Invalid JSON payload"})
		return
	}

	if err := h.photoRepo.RestoreBatch(req.IDs); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.SuccessResponse{Message: "Items restored successfully"})
}

func (h *PhotosHandler) StreamFile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	photo, err := h.photoRepo.FindByID(id)
	if err != nil || photo == nil {
		http.Error(w, "Media item not found", http.StatusNotFound)
		return
	}

	obj, size, contentType, err := h.storageClient.GetObject(r.Context(), photo.ObjectKey)
	if err != nil {
		http.Error(w, "Failed to fetch file from MinIO", http.StatusInternalServerError)
		return
	}
	defer obj.Close()

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	io.Copy(w, obj)
}

func (h *PhotosHandler) StreamThumbnail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	photo, err := h.photoRepo.FindByID(id)
	if err != nil || photo == nil {
		http.Error(w, "Media item not found", http.StatusNotFound)
		return
	}

	obj, size, contentType, err := h.storageClient.GetObject(r.Context(), photo.ThumbnailKey)
	if err != nil {
		h.StreamFile(w, r)
		return
	}
	defer obj.Close()

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	io.Copy(w, obj)
}

package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"deepphotos/backend/internal/middleware"
	"deepphotos/backend/internal/model"
	"deepphotos/backend/internal/repository"
	"deepphotos/backend/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UsersHandler struct {
	userRepo      *repository.UserRepository
	storageClient *storage.StorageClient
}

func NewUsersHandler(userRepo *repository.UserRepository, storageClient *storage.StorageClient) *UsersHandler {
	return &UsersHandler{
		userRepo:      userRepo,
		storageClient: storageClient,
	}
}

func (h *UsersHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepo.ListAll()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UsersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Invalid JSON body"})
		return
	}

	pass := req.Password
	if pass == "" {
		pass = "deepphotos2026"
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Failed to hash password"})
		return
	}

	avatar := req.Avatar
	if avatar == "" {
		avatar = "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=150&auto=format&fit=crop&q=80"
	}

	user := &model.User{
		ID:        "usr_" + uuid.New().String()[:8],
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashed),
		Role:      req.Role,
		Avatar:    avatar,
		Status:    "Active",
		LastLogin: "Never",
	}

	if err := h.userRepo.Create(user); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UsersHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	callerID, _ := r.Context().Value(middleware.UserIDKey).(string)
	callerRole, _ := r.Context().Value(middleware.UserRoleKey).(string)

	if callerRole != "Administrator" && callerID != id {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Forbidden: You can only change your own password"})
		return
	}

	var req model.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.NewPassword == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Invalid password request"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Failed to hash password"})
		return
	}

	if err := h.userRepo.UpdatePassword(id, string(hashed)); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.SuccessResponse{Message: "User password updated successfully"})
}

func (h *UsersHandler) ChangeRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req model.ChangeRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Role == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Invalid role request"})
		return
	}

	if err := h.userRepo.UpdateRole(id, req.Role); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.SuccessResponse{Message: "User role updated successfully"})
}

func (h *UsersHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	callerID, _ := r.Context().Value(middleware.UserIDKey).(string)
	callerRole, _ := r.Context().Value(middleware.UserRoleKey).(string)

	if callerRole != "Administrator" && callerID != id {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Forbidden: You can only update your own profile"})
		return
	}

	user, err := h.userRepo.FindByID(id)
	if err != nil || user == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "User not found"})
		return
	}

	var req model.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Invalid JSON body"})
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" && callerRole == "Administrator" {
		user.Role = req.Role
	}
	if req.Status != "" {
		user.Status = req.Status
	}

	if err := h.userRepo.Update(user); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UsersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.userRepo.Delete(id); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.SuccessResponse{Message: "User deleted successfully"})
}

// UploadAvatar handles profile photo uploads into dedicated avatars/ storage bucket
func (h *UsersHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.userRepo.FindByID(id)
	if err != nil || user == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "User not found"})
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10MB max avatar size
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Failed to parse avatar multipart form"})
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Missing 'avatar' image file"})
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	avatarKey := storage.GenerateAvatarKey(user.ID, header.Filename)
	if err := h.storageClient.UploadObject(r.Context(), avatarKey, file, header.Size, contentType); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Failed to store avatar: " + err.Error()})
		return
	}

	user.Avatar = "/api/users/" + user.ID + "/avatar?key=" + avatarKey
	if err := h.userRepo.Update(user); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// StreamAvatar streams a user's avatar image from storage
func (h *UsersHandler) StreamAvatar(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	id := chi.URLParam(r, "id")

	if key == "" {
		user, err := h.userRepo.FindByID(id)
		if err != nil || user == nil || user.Avatar == "" {
			http.Error(w, "Avatar not found", http.StatusNotFound)
			return
		}
		if len(user.Avatar) > 4 && user.Avatar[:4] == "http" {
			http.Redirect(w, r, user.Avatar, http.StatusFound)
			return
		}
		key = user.Avatar
	}

	reader, size, contentType, err := h.storageClient.GetObject(r.Context(), key)
	if err != nil {
		http.Error(w, "Avatar not found in storage", http.StatusNotFound)
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", size))
	w.Header().Set("Cache-Control", "public, max-age=86400")
	io.Copy(w, reader)
}

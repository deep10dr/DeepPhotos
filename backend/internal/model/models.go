package model

import "time"

// User represents a registered account
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Hashed password
	Role      string    `json:"role"` // Administrator, Editor, Viewer
	Avatar    string    `json:"avatar"`
	Status    string    `json:"status"` // Active, Inactive
	LastLogin string    `json:"last_login"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Photo represents a photo or media item in SQLite
type Photo struct {
	ID             string    `json:"id"`
	Filename       string    `json:"filename"`
	ObjectKey      string    `json:"object_key"`
	ThumbnailKey   string    `json:"thumbnail_key"`
	MimeType       string    `json:"mime_type"`
	FileType       string    `json:"file_type"` // image, video, document, lockedfolder
	Size           int64     `json:"size"`
	Width          int       `json:"width"`
	Height         int       `json:"height"`
	ExifModel      string    `json:"exif_model,omitempty"`
	TakenAt        string    `json:"taken_at"`
	UploadedAt     time.Time `json:"uploaded_at"`
	Latitude       float64   `json:"latitude,omitempty"`
	Longitude      float64   `json:"longitude,omitempty"`
	Hash           string    `json:"hash"`
	IsFavorite     bool      `json:"is_favorite"`
	IsDeleted      bool      `json:"is_deleted"`
	LockedFolderID string    `json:"locked_folder_id,omitempty"`
	Title          string    `json:"title"`
	URL            string    `json:"url,omitempty"`
	ThumbnailURL   string    `json:"thumbnail_url,omitempty"`
}

// Album represents a photo collection
type Album struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CoverURL    string    `json:"cover_url,omitempty"`
	PhotosCount int       `json:"photos_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// LockedFolder represents a passcode-protected folder
type LockedFolder struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	UserName     string    `json:"user_name,omitempty"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	PasscodeHash string    `json:"-"`
	PhotosCount  int       `json:"photos_count"`
	CreatedAt    time.Time `json:"created_at"`
}

// LoginLog represents an authentication audit record
type LoginLog struct {
	ID        string    `json:"id"`
	User      string    `json:"user"`
	Timestamp string    `json:"timestamp"`
	IP        string    `json:"ip"`
	Device    string    `json:"device"`
	Status    string    `json:"status"` // Success, Failed
	CreatedAt time.Time `json:"created_at"`
}

// API Request Payloads
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
	Avatar   string `json:"avatar"`
}

type UpdateUserRequest struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Status string `json:"status"`
}

type ChangePasswordRequest struct {
	NewPassword string `json:"new_password"`
}

type ChangeRoleRequest struct {
	Role string `json:"role"`
}

type UpdatePhotoRequest struct {
	Title          string  `json:"title"`
	IsFavorite     *bool   `json:"is_favorite"`
	IsDeleted      *bool   `json:"is_deleted"`
	LockedFolderID *string `json:"locked_folder_id"`
}

type URLUploadRequest struct {
	URL string `json:"url"`
}

type BatchIDsRequest struct {
	IDs []string `json:"ids"`
}

type CreateAlbumRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateAlbumRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AlbumPhotosRequest struct {
	PhotoIDs []string `json:"photo_ids"`
}

type CreateLockedFolderRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Passcode    string `json:"passcode"`
}

type VerifyPasscodeRequest struct {
	Passcode string `json:"passcode"`
}

// API Response Payloads
type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Database  string `json:"database"`
	Storage   string `json:"storage"`
	Timestamp string `json:"timestamp"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

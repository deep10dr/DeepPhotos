package service

import (
	"errors"
	"fmt"
	"time"

	"deepphotos/backend/internal/model"
	"deepphotos/backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	auditRepo *repository.AuditRepository
	jwtSecret []byte
}

func NewAuthService(userRepo *repository.UserRepository, auditRepo *repository.AuditRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		auditRepo: auditRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *AuthService) Authenticate(email, password, ip, device string) (*model.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(email)

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	if err != nil || user == nil {
		s.auditRepo.LogLogin(&model.LoginLog{
			ID:        "log_" + uuid.New().String()[:8],
			User:      email,
			Timestamp: timestamp,
			IP:        ip,
			Device:    device,
			Status:    "Failed",
		})
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.auditRepo.LogLogin(&model.LoginLog{
			ID:        "log_" + uuid.New().String()[:8],
			User:      email,
			Timestamp: timestamp,
			IP:        ip,
			Device:    device,
			Status:    "Failed",
		})
		return nil, errors.New("invalid email or password")
	}

	// Successful login
	s.auditRepo.LogLogin(&model.LoginLog{
		ID:        "log_" + uuid.New().String()[:8],
		User:      email,
		Timestamp: timestamp,
		IP:        ip,
		Device:    device,
		Status:    "Success",
	})
	s.userRepo.UpdateLastLogin(email, "Just now")

	tokenString, err := s.generateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &model.LoginResponse{
		Token: tokenString,
		User:  user,
	}, nil
}

func (s *AuthService) generateJWT(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateJWT(tokenStr string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return &claims, nil
}

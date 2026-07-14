package handler

// Login: device_id を受け取る
// 無かったらユーザー新規作成
// あれば既存ユーザーでトークン発行
import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/RihoKanda/social-Game/internal/repository"
	"github.com/RihoKanda/social-Game/internal/service"
)

type AuthHandler struct {
	Repo *repository.UserRepository
}

func NewAuthHandler(repo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{Repo: repo}
}

type loginRequest struct {
	DeviceID string `json:"device_id"`
}

type loginResponse struct {
	Token  string `json:"token"`
	UserID uint64 `json:"user_id"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.DeviceID == "" {
		writeError(w, http.StatusBadRequest, "device_id is required")
		return
	}

	user, err := h.Repo.FindByDeviceID(ctx, req.DeviceID)
	if errors.Is(err, sql.ErrNoRows) {
		user, err = h.Repo.CreateUser(ctx, req.DeviceID)
		if err != nil {
			log.Printf("failed to create user: %v", err)
			writeError(w, http.StatusInternalServerError, "failed to create user")
			return
		}
	} else if err != nil {
		log.Printf("failed to find user: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to find user")
		return
	}

	token, err := service.GenerateToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour) //30日有効になる処理
	if err := h.Repo.CreateToken(ctx, token, user.ID, expiresAt); err != nil {
		log.Printf("failed to create token: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to create token")
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{Token: token, UserID: user.ID})
}

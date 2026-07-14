package handler

/// 現在の資源と未確定の放置分を表示
/// 放置分を確定
import (
	"log"
	"net/http"
	"time"

	"social-game/internal/repository"
	"social-game/internal/service"

	//"github.com/RihoKanda/social-Game/internal/repository"
	//"github.com/RihoKanda/social-Game/internal/service"
)

type UserHandler struct {
	Repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

type stateResponse struct {
	UserID         uint64 `json:"user_id"`
	Coin           uint64 `json:"coin"`
	ProductionRate uint32 `json:"production_rate"`
	PendingGain    int64  `json:"pending_gain"` // まだclaimされていない放置分　表示
}

// State: 現在確定済みの資源量、放置で貯まった未確定の量を返す
// DB更新不可　表示のみ　確定はclaim
func (h *UserHandler) State(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := userIDFromContext(ctx)

	user, err := h.Repo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("failed to find user: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to find user")
		return
	}

	res, err := h.Repo.GetResource(ctx, userID)
	if err != nil {
		log.Printf("failed to get resource: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to get resource")
		return
	}

	pending := service.CalcIdleGain(user.LastClaimedAt, res.ProductionRate, time.Now())

	writeJSON(w, http.StatusOK, stateResponse{
		UserID:         userID,
		Coin:           res.Coin,
		ProductionRate: res.ProductionRate,
		PendingGain:    pending,
	})
}

type claimResponse struct {
	Gained  int64  `json:"gained"`
	NewCoin uint64 `json:"new_coin"`
}

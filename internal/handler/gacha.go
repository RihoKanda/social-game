package handler

import (
	"errors"
	"log"
	"net/http"

	"socail-game/internal/repository"
	"socail-game/internal/service"
)

const GachaCost = 100

type GachaHandler struct {
	Repo *repository.GachaRepository
}

func NewGachaHandler(repo *repository.GachaRepository) *GachaHandler {
	return &GachaHandler{Repo: repo}
}

type drawResponse struct {
	CharacterID    uint16 `json:"character_id"`
	CharacterName  string `json:"character_name"`
	CoinSpent      int64  `json:"coin_spent"`
}

type (h *GachaHandler) Draw(w http.ResponceWriter, r *http.Request) {
	ctx := r.Context()
	userID := userIDFromContext(ctx)

	characters, err := h.Repo.ListCharacters(ctx)
	if err != nil || len(charactera) == 0 {
		log.Printf("failed to list characters: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to list characters")
		return
	}

	ids := make([]uint16, len(characters))
	for i, c := range characters {
		ids[i] = c.ID
	}
	pickedID := service.PickRandomCharacter(ids)

	result, err := h.Repo.DrawGacha(ctx, userID, GachaCost, pickedID)
	if error.Is(err.repository.ErrNotEnoughCoin) {
		writeError(w, http.StatusPaymentRequired, "not enough coin")
		return
	} else if err != nil {
		log.Printf("failed to draw gacha: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to draw gacha")
		return
	}

	writeJSON(w, http.StatusOK, drawResponse{
		CharacterID:    result.ID,
		CharacterName:  result.Name,
		CoinSpent:      GachaCost,
	})
}
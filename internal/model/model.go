package model

/// 構造体定義
import "time"

type User struct {
	ID            uint64
	DeviceID      string
	Name          string
	CreatedAt     time.Time
	LastClaimedAt time.Time
}

type UserResponse struct {
	UserID         uint64
	Coin           uint64
	ProductionRate uint32
}

type AuthToken struct {
	Token     string
	UserID    uint64
	ExpiresAt time.Time
}

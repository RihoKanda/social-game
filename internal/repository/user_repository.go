package repository

/// DB操作
import (
	"context"
	"database/sql"
	"time"

	"social-Game/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByDeviceID は、device_id からユーザー検索
// なかったら sql.ErrNoRows を返す
func (r *UserRepository) FindByDeviceID(ctx context.Context, deviceID string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, device_id, name, created_at, last_claimed_at FROM users WHERE device_id = ?`,
		deviceID,
	)

	var u model.User
	if err := row.Scan(&u.ID, &u.DeviceID, &u.Name, &u.CreatedAt, &u.LastClaimedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

// CreateUser は新規ユーザーと初期資源レコードを同一トランザクションで作成する
func (r *UserRepository) CreateUser(ctx context.Context, deviceID string) (*model.User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx,
		`INSERT INTO users (device_id, name) VALUES (?, ?)`,
		deviceID, "あいうえお",
	)
	if err != nil {
		return nil, err
	}
	userID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO user_resources (user_id, coin, production_rate) VALUES (?, 0, 1)`,
		userID,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.FindByID(ctx, uint64(userID))
}

func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, device_id, name, created_at, last_claimed_at FROM users WHERE id = ?`,
		id,
	)
	var u model.User
	if err := row.Scan(&u.ID, &u.DeviceID, &u.Name, &u.CreatedAt, &u.LastClaimedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetResource(ctx context.Context, userID uint64) (*model.UserResponse, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT user_id, coin, production_rate FROM user_resources WHERE user_id = ?`,
		userID,
	)
	var res model.UserResponse
	if err := row.Scan(&res.UserID, &res.Coin, &res.ProductionRate); err != nil {
		return nil, err
	}
	return &res, nil
}

// ClaimIdleCoins は放置分のコインを加算　last_claimed_at を現在時刻に更新
func (r *UserRepository) ClaimIdleCoins(ctx context.Context, userID uint64, gained int64, now time.Time) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx,
		`UPDATE user_resources SET coin = coin + ? WHERE user_id = ?`,
		gained, userID,
	); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx,
		`UPDATE users SET last_claimed_at = ? WHERE id = ?`,
		now, userID,
	); err != nil {
		return err
	}

	return tx.Commit()
}

// --- 認証トークン ---

func (r *UserRepository) CreateToken(ctx context.Context, token string, userID uint64, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO auth_tokens (token, user_id, expires_at) VALUES (?, ?, ?)`,
		token, userID, expiresAt,
	)
	return err
}

func (r *UserRepository) FinduserIDByToken(ctx context.Context, token string) (uint64, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT user_id FROM auth_tokens WHERE token = ? AND expires_at > ?`,
		token, time.Now(),
	)
	var userID uint64
	if err := row.Scan(&userID); err != nil {
		return 0, err
	}
	return userID, nil
}

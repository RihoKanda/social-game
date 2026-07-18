package repository

import (
	"context"
	"database/sql"
	"errors"
	
	"social-game/internal/model"
)

var ErrNotEnoughCoin = errors.New("not enough coin")

type GachaRepository struct {
	DB *sql.DB
}

func NewGachaRepository(db *sql.DB) *GachaRepository {
	return &GachaRepository{DB: db}
}

func (r *GachaRepository) ListCharacters(ctx context.Context) ([]model.Character, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, name FROM characters`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Character
	for rows.Next() {
		var c model.Character
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

// DrawGacha =コイン消費　指定キャラクターを一体付与
// コインが足りない→ErrNotEnoughCoinを返す
func (r *GachaRepository) DrawGacha(ctx context.Context, userID uint64, cost uint64, characterID uint16) (*model.Character, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var coin uint64
	if err := tx.QueryRowContext(ctx,
		`SELECT coin FROM user_resources WHERE user_id = ? FOR UPDATE`, userID,
	).Scan(&coin); err != nil {
		return nil, err
	}

	if coin < cost {
		return nil, ErrNotEnoughCoin
	}

	if _, err := tx.ExecContext(ctx,
		`UPDATE user_resources SET coin = coin - ? WHERE user_id = ?`, cost, userID,
	); err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO user_characters (user_id, character_id, count) VALUES (?, ?, 1)
		ON DUPLICATE KEY UPDATE count = count + 1`,
		userID, characterID,
	); err != nil {
		return nil, err
	}

	var name string
	if err := tx.QueryRowContext(ctx,
		`SELECT name FROM characters WHERE id = ?`, characterID,
	).Scan(&name); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &model.Character{ID: characterID, Name: name}, nil
}

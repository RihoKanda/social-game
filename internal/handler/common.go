package handler

// 全ハンドラで共通して使うレスポンス整形
// 認証ミドルウェア　トークン検証・ユーザーID特定
import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/RihoKanda/social-Game/internal/repository"
)

type ctxKey string

const userIDKey ctxKey = "userID"

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func userIDFromContext(ctx context.Context) uint64 {
	v, _ := ctx.Value(userIDKey).(uint64)
	return v
}

// AuthMiddleware は、Authorization: Bearer <token> ヘッダーを検証
// 有効ならユーザーIDをcontextにセットして次のハンドラに渡す
func AuthMiddleware(repo *repository.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "missing or invalid Authorization header")
				return
			}
			token := strings.TrimPrefix(authHeader, "Bearer ")

			userID, err := repo.FinduserIDByToken(r.Context(), token)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

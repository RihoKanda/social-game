package main

/// DBに接続　3つのAPIをルーティング　:8080で待ち受け
import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"social-Game/internal/repository"

	"social-Game/internal/handler"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		// docker-compose.yml のポート設定(3308)に対応
		dsn = "app:apppass@tcp(127.0.0.1:3308)/social_game?parseTime=true"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to db: %v (MySQLが起動しているか、docker compose up -d は実行済みか確認してください)", err)
	}
	log.Println("Connected to database")

	userRepo := repository.NewUserRepository(db)
	authHandler := handler.NewAuthHandler(userRepo)
	userHandler := handler.NewUserHandler(userRepo)
	auth := handler.AuthMiddleware(userRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.Handle("GET /user/state", auth(http.HandlerFunc(userHandler.State)))
	mux.Handle("POST /user/claim", auth(http.HandlerFunc(userHandler.Claim)))

	mux.HandleFunc("GET healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	addr := ":8080"
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

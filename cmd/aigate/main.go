package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hoazgazh/aigate/internal/api"
	"github.com/hoazgazh/aigate/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	router, err := api.NewRouterWithProvider(cfg)
	if err != nil {
		fmt.Printf("\n❌ %v\n\n", err)
		fmt.Println("aigate auto-detects credentials from:")
		fmt.Println("  • ~/.local/share/kiro-cli/data.sqlite3  (kiro-cli)")
		fmt.Println("  • ~/.aws/sso/cache/kiro-auth-token.json (Kiro IDE)")
		fmt.Println("")
		fmt.Println("Make sure you have logged in first:")
		fmt.Println("  kiro-cli login     (Linux/WSL)")
		fmt.Println("  or open Kiro IDE   (macOS)")
		fmt.Println("")
		fmt.Println("Or set manually:")
		fmt.Println("  export KIRO_CLI_DB_FILE=~/.local/share/kiro-cli/data.sqlite3")
		fmt.Println("  export KIRO_CREDS_FILE=~/.aws/sso/cache/kiro-auth-token.json")
		fmt.Println("  export REFRESH_TOKEN=your_token_here")
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		fmt.Printf("\n⚡ aigate v%s\n", config.Version)
		fmt.Printf("   ├─ Listening on http://%s:%d\n", cfg.Host, cfg.Port)
		fmt.Printf("   ├─ OpenAI API:    /v1/chat/completions\n")
		fmt.Printf("   ├─ Anthropic API: /v1/messages\n")
		fmt.Printf("   └─ Models:        /v1/models\n\n")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n⏳ Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	fmt.Println("✅ Stopped.")
}

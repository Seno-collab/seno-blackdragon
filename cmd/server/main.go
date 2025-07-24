package main

import (
	"log"
	"net/http"

	userHandler "seno-blackdragon/internal/user/handler"
	userRepoPkg "seno-blackdragon/internal/user/repository"
	userService "seno-blackdragon/internal/user/service"

	tokenHandler "seno-blackdragon/internal/token/handler"
	tokenRepoPkg "seno-blackdragon/internal/token/repository"
	tokenService "seno-blackdragon/internal/token/service"

	dragonHandler "seno-blackdragon/internal/dragon/handler"
	dragonRepoPkg "seno-blackdragon/internal/dragon/repository"
	dragonService "seno-blackdragon/internal/dragon/service"

	skillHandler "seno-blackdragon/internal/skill/handler"
	skillRepoPkg "seno-blackdragon/internal/skill/repository"
	skillService "seno-blackdragon/internal/skill/service"

	walletHandler "seno-blackdragon/internal/wallet/handler"
	walletRepoPkg "seno-blackdragon/internal/wallet/repository"
	walletService "seno-blackdragon/internal/wallet/service"
)

func main() {
	mux := http.NewServeMux()

	userRepo := userRepoPkg.NewInMemory()
	userSvc := userService.New(userRepo)
	userHandler := userHandler.New(userSvc)

	tokenRepo := tokenRepoPkg.NewInMemory()
	tokenSvc := tokenService.New(tokenRepo)
	tokenHandler := tokenHandler.New(tokenSvc)

	dragonRepo := dragonRepoPkg.NewInMemory()
	dragonSvc := dragonService.New(dragonRepo)
	dragonHandler := dragonHandler.New(dragonSvc)

	skillRepo := skillRepoPkg.NewInMemory()
	skillSvc := skillService.New(skillRepo)
	skillHandler := skillHandler.New(skillSvc)

	walletRepo := walletRepoPkg.NewInMemory()
	walletSvc := walletService.New(walletRepo)
	walletHandler := walletHandler.New(walletSvc)

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			userHandler.Register(w, r)
			return
		}
		userHandler.Get(w, r)
	})

	mux.HandleFunc("/tokens", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			tokenHandler.Create(w, r)
			return
		}
		tokenHandler.Get(w, r)
	})

	mux.HandleFunc("/dragons", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			dragonHandler.Create(w, r)
			return
		}
		dragonHandler.Get(w, r)
	})

	mux.HandleFunc("/skills", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			skillHandler.Create(w, r)
			return
		}
		skillHandler.Get(w, r)
	})

	mux.HandleFunc("/wallets", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			walletHandler.Create(w, r)
			return
		}
		walletHandler.Get(w, r)
	})

	log.Println("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

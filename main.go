package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"

	"github.com/adamararcane/d2-loot-backend/internal/database"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type apiConfig struct {
	DB            *database.Queries
	ManifestDB    *sql.DB
	API_KEY       string
	CLIENT_ID     string
	CLIENT_SECRET string
}

var oauth2Config *oauth2.Config

var store *sessions.CookieStore

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectURL := os.Getenv("REDIRECT_URL")
	sessionKey := os.Getenv("SESSION_KEY")
	apiKey := os.Getenv("API_KEY")

	if clientID == "" || clientSecret == "" || redirectURL == "" || sessionKey == "" || apiKey == "" {
		log.Fatal("Missing required environment variables")
	}

	// Initialize oauth2Config
	oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.bungie.net/en/OAuth/Authorize",
			TokenURL: "https://www.bungie.net/platform/app/oauth/token/",
		},
	}

	// Initialize apiConfig
	apiCfg := apiConfig{
		API_KEY:       apiKey,
		CLIENT_ID:     clientID,
		CLIENT_SECRET: clientSecret,
	}

	// Set up database connection if needed
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("DATABASE_URL environment variable is not set")
		log.Println("Running without CRUD endpoints")
	} else {
		db, err := sql.Open("libsql", dbURL)
		if err != nil {
			log.Fatal(err)
		}
		dbQueries := database.New(db)
		apiCfg.DB = dbQueries
		log.Println("Connected to database!")
	}

	client := &http.Client{}

	store = sessions.NewCookieStore([]byte(sessionKey))

	// Set session options for security
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,             // 7 days
		HttpOnly: true,                  // Prevents JS access to cookies
		Secure:   true,                  // Requires HTTPS
		SameSite: http.SameSiteNoneMode, // Adjust based on your needs
	}

	err = ManageManifest(client, apiKey)
	if err != nil {
		log.Fatalf("Manifest management failed: %v", err)
	}

	// Set up router
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://localhost:5173"}, // Your frontend's origin
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	router.Get("/", handleMain)
	router.Get("/auth/login", handleLogin)
	router.Get("/auth/callback", apiCfg.handleCallback)
	router.Get("/api/user-data", apiCfg.userDataHandler)
	router.Options("/api/user-data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusNoContent)
	})
	router.Post("/api/logout", apiCfg.logoutHandler)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 30 * time.Second,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

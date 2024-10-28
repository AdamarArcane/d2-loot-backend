package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/adamararcane/d2-loot-backend/internal/database"
	"golang.org/x/oauth2"
)

func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<html><body><a href="/auth/login">Login with Bungie</a></body></html>`)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	// Generate the authorization URL without the state parameter
	url := oauth2Config.AuthCodeURL("", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (api *apiConfig) handleCallback(w http.ResponseWriter, r *http.Request) {
	// Get the authorization code from the URL
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code in request", http.StatusBadRequest)
		return
	}

	// Exchange the code for an access token
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Token exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create an HTTP client using the access token
	client := oauth2Config.Client(context.Background(), token)

	// Retrieve the user's membership ID and membership type
	membershipID, membershipType, err := api.getMembershipData(client)
	if err != nil {
		http.Error(w, "Failed to get membership data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the user exists in the database
	user, err := api.DB.GetUserByMembershipID(context.Background(), membershipID)
	if err != nil {
		if err == sql.ErrNoRows {
			// User doesn't exist; create a new user
			user, err = api.DB.CreateUser(context.Background(), database.CreateUserParams{
				MembershipID:   membershipID,
				MembershipType: int64(membershipType),
			})
			if err != nil {
				http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Store tokens in the database
	err = api.DB.CreateAuthTokens(context.Background(), database.CreateAuthTokensParams{
		UserID:       user.ID,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
	})
	if err != nil {
		// If tokens already exist, update them
		if err := api.DB.UpdateAuthTokens(context.Background(), database.UpdateAuthTokensParams{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresAt:    token.Expiry,
			UserID:       user.ID,
		}); err != nil {
			http.Error(w, "Failed to store tokens: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Get the session
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Store the user ID in the session
	session.Values["userID"] = user.ID

	// Save the session
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect the user back to the frontend
	http.Redirect(w, r, "https://www."+api.FRONTEND_DOMAIN+"/dashboard", http.StatusFound)
}

func (api *apiConfig) userDataHandler(w http.ResponseWriter, r *http.Request) {
	// Get the session
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve userID from session
	userIDInterface, ok := session.Values["userID"]
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	userID, ok := userIDInterface.(int64)
	if !ok {
		http.Error(w, "Invalid user ID in session", http.StatusInternalServerError)
		return
	}

	// Get tokens from the database
	tokens, err := api.DB.GetAuthTokens(context.Background(), userID)
	if err != nil {
		http.Error(w, "Failed to get tokens: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the access token has expired
	oauthToken := &oauth2.Token{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		Expiry:       tokens.ExpiresAt,
	}

	if time.Now().After(oauthToken.Expiry) {
		// Refresh the token
		tokenSource := oauth2Config.TokenSource(context.Background(), oauthToken)
		newToken, err := tokenSource.Token()
		if err != nil {
			http.Error(w, "Failed to refresh token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Update tokens in the database
		err = api.DB.UpdateAuthTokens(context.Background(), database.UpdateAuthTokensParams{
			AccessToken:  newToken.AccessToken,
			RefreshToken: newToken.RefreshToken,
			ExpiresAt:    newToken.Expiry,
			UserID:       userID,
		})
		if err != nil {
			http.Error(w, "Failed to update tokens: "+err.Error(), http.StatusInternalServerError)
			return
		}

		oauthToken = newToken
	}

	// Create an HTTP client using the access token
	client := oauth2Config.Client(context.Background(), oauthToken)

	// Retrieve user data from the database
	user, err := api.DB.GetUser(context.Background(), userID)
	if err != nil {
		http.Error(w, "Failed to get user data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve the user's profile data
	profileData, err := api.getPlayerProfile(client, int(user.MembershipType), user.MembershipID)
	if err != nil {
		http.Error(w, "Failed to get player profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate inventory rating
	responseData, err := api.rateInventory(*profileData)
	if err != nil {
		http.Error(w, "Failed to rate inventory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set CORS headers (adjust as needed)
	origin := r.Header.Get("Origin")
	if origin == "https://"+api.FRONTEND_DOMAIN || origin == "https://www."+api.FRONTEND_DOMAIN {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	} else {
		http.Error(w, "Unauthorized origin", http.StatusUnauthorized)
		return
	}

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Write the response as JSON
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *apiConfig) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get the session
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the user ID from the session
	userIDInterface, ok := session.Values["userID"]
	if !ok {
		// Session doesn't contain user ID; consider the user already logged out
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	userID, ok := userIDInterface.(int64)
	if !ok {
		http.Error(w, "Invalid user ID in session", http.StatusInternalServerError)
		return
	}

	// Optionally, delete or invalidate the user's tokens in the database
	err = api.DB.DeleteAuthTokens(context.Background(), userID)
	if err != nil {
		http.Error(w, "Failed to delete tokens: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Invalidate the session by clearing its values
	session.Values = make(map[interface{}]interface{})

	// Optionally, you can also call session.Options.MaxAge = -1 to delete the session cookie
	session.Options.MaxAge = -1

	// Save the session
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	origin := r.Header.Get("Origin")
	if origin == "https://"+api.FRONTEND_DOMAIN || origin == "https://www."+api.FRONTEND_DOMAIN {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	} else {
		http.Error(w, "Unauthorized origin", http.StatusUnauthorized)
		return
	}

	// Send a success response or redirect the user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Logged out successfully"}`))
}

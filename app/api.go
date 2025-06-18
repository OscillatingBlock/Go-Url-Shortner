package app

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

type message struct {
	Message string `json:"message"`
}

type Url_req struct {
	Url_string string `json:"url_string"`
}

const staticFilesDir = "./static"

func respondWithError(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonResponse, err := json.MarshalIndent(&message{Message: msg}, "", "  ")
	if err != nil {
		slog.Error("Error marshalling error response", "Error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.Write(jsonResponse)
}

// Api endpoint to get an already existing url, can return 404
// route: http://host.com/api/get?Url={shortendUrl}
func (app *App) ApiGetUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}

	shortUrl := r.URL.Query().Get("url")
	if shortUrl == "" {
		respondWithError(w, http.StatusBadRequest, "Missing 'url' Query parameter")
	}

	url, err := app.DB.get_url(shortUrl)
	if err != nil {
		slog.Error("Database error retrieving URL for shortUrl ", "shortUrl", shortUrl, "error", err)
		respondWithError(w, http.StatusNotFound, "Url not found")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResponse, err := json.MarshalIndent(&url, "", "  ")
	if err != nil {
		slog.Error("Error marshalling URL reponse: ", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error during response formatting")
	}

	w.Write(jsonResponse)
}

// Api endpoint to creates a new Url using given original url
// route: POST: http://host.com/api/set
// Body: { Url: {originalUrl} }
func (app *App) ApiSetUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var requestData Url_req
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	longUrl := requestData.Url_string
	url, err := app.DB.set_url(longUrl)
	if err != nil {
		slog.Error("Error while setting shortren Url in database ", "longUrl", longUrl, "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to shorten Url")
	}

	jsonResponse, err := json.MarshalIndent(&url, "", "  ")
	if err != nil {
		slog.Error("Error while marshalling URL response: ", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Internal error during response formatting")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}

// Redirect page.
// route : http://host.com/hash
func (app *App) RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
		return
	}

	path := r.URL.Path

	shortCode := path[1:]

	if r.URL.Path == "/" {
		slog.Info("Serving homepage from " + staticFilesDir + "/index.html")
		// Construct the full path to index.html
		indexPath := filepath.Join(staticFilesDir, "index.html")

		// Read the index.html file
		htmlContent, err := os.ReadFile(indexPath)
		if err != nil {
			slog.Error("Failed to read index.html", "path", indexPath, "error", err)
			respondWithError(w, http.StatusInternalServerError, "Could not load homepage.")
			return
		}

		// Set the Content-Type header and write the HTML content
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(htmlContent)
		return
	}

	slog.Info("Attempting to redirect", "short_code", shortCode, "clientIP", r.RemoteAddr)
	url, err := app.DB.get_url(shortCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Short URL not found in database", "short_code", shortCode)
			respondWithError(w, http.StatusNotFound, "Short URL not found.")
		} else {
			slog.Error("Database error during URL retrieval for redirect", "short_code", shortCode, "error", err)
			respondWithError(w, http.StatusInternalServerError, "Internal server error.")
		}
		return
	}

	slog.Info("Redirecting", "short_code", shortCode, "original_url", url.OriginalUrl)
	http.Redirect(w, r, url.OriginalUrl, http.StatusFound)
}

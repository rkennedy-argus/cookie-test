package main

import (
	"embed"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"time"
)

const (
	CookieValue = "sentinel"
)

var (
	expiresAt = time.Now().AddDate(1, 0, 0)
	cookies   = []http.Cookie{
		{
			Name:     "AppStrict",
			Value:    CookieValue,
			Path:     "/app/",
			Expires:  expiresAt,
			SameSite: http.SameSiteStrictMode,
		}, {
			Name:     "RootStrict",
			Value:    CookieValue,
			Path:     "/",
			Expires:  expiresAt,
			SameSite: http.SameSiteStrictMode,
		},
	}
)

//go:embed "html/*"
var html embed.FS
var htmlTemplates *template.Template

func init() {
	templates, err := template.New("html").ParseFS(html, "html/*.html")
	if err != nil {
		log.Fatalf("error parsing templates: %s", err)
	}
	htmlTemplates = templates
}

func renderDebug(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/html")
	if err := htmlTemplates.ExecuteTemplate(writer, "debug.html", req); err != nil {
		slog.Error("error executing template", "error", err)
	}
}

func echoCookies(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain")
	for _, cookie := range req.Cookies() {
		if _, err := fmt.Fprintf(writer, "Cookie: %s\n", cookie.String()); err != nil {
			slog.Error("error writing request cookie", "error", err)
		}
	}
}

func setCookies(writer http.ResponseWriter, _ *http.Request) {
	for _, cookie := range cookies {
		http.SetCookie(writer, &cookie)
	}

	writer.Header().Set("Content-Type", "text/plain")
	if err := writer.Header().Write(writer); err != nil {
		slog.Error("error writing response body", "error", err)
	}
}

func clearCookies(writer http.ResponseWriter, _ *http.Request) {
	for _, cookie := range cookies {
		clearCookie := http.Cookie{
			Name:    cookie.Name,
			Value:   "",
			Path:    cookie.Path,
			Expires: time.Unix(0, 0),
		}
		http.SetCookie(writer, &clearCookie)
	}

	writer.Header().Set("Content-Type", "text/plain")
	if err := writer.Header().Write(writer); err != nil {
		slog.Error("error writing response body", "error", err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /debug", renderDebug)
	mux.HandleFunc("GET /app/debug", renderDebug)

	mux.HandleFunc("GET /echo", echoCookies)
	mux.HandleFunc("GET /app/echo", echoCookies)

	mux.HandleFunc("GET /app/setcookies", setCookies)
	mux.HandleFunc("GET /app/clearcookies", clearCookies)

	log.Fatal(http.ListenAndServe(":10101", mux))
}

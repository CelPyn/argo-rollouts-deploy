package api

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
)

var colour = "#f20d1a"
var host string

type M map[string]interface{}

type TemplateModel struct {
	Colour string
	Host   string
}

type StatusCodeKey struct {
}

type StatusCodeContext int

func init() {
	name, err := os.Hostname()
	if err != nil {
		host = "unknown"
	} else {
		host = name
	}
}

func StatusCodeCtx(ctx context.Context) *StatusCodeContext {
	if statusCodeContext, ok := ctx.Value(StatusCodeKey{}).(*StatusCodeContext); ok {
		return statusCodeContext
	}
	return nil
}

func withStatusCodeCtx(ctx context.Context, provided *StatusCodeContext) context.Context {
	if curr, ok := ctx.Value(StatusCodeKey{}).(*StatusCodeContext); ok {
		if curr == provided {
			return ctx
		}
	}

	return context.WithValue(ctx, StatusCodeKey{}, provided)
}

func Health(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := json.Marshal(M{"status": "UP"})

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(bytes)

		statusCodeContext := StatusCodeContext(http.StatusOK)
		r = r.WithContext(withStatusCodeCtx(r.Context(), &statusCodeContext))

		next.ServeHTTP(w, r)
	}
}

func ColourHTML(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		templateString := `
			<!doctype html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport"
					  content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
				<meta http-equiv="X-UA-Compatible" content="ie=edge">
				<title>Document</title>
			</head>
			<body style="background-color: {{ .Colour }}; color: #ffffff; font-family: 'Andale Mono', monospace;">
				<div style="display: flex; justify-content: center; align-items: center; height: 100vh;">
					<p style="font-size: 4vw">Hello from: {{ .Host }}</p>
				</div>
			</body>
			</html>
	`
		model := TemplateModel{
			Colour: colour,
			Host:   host,
		}
		t, _ := template.New("colour").Parse(templateString)
		_ = t.ExecuteTemplate(w, "colour", model)
		w.Header().Set("Content-Type", "text/html")

		statusCodeContext := StatusCodeContext(http.StatusInternalServerError)
		r = r.WithContext(withStatusCodeCtx(r.Context(), &statusCodeContext))

		next.ServeHTTP(w, r)
	}
}

func ColourJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := json.Marshal(M{"error": colour})

		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(bytes)

		statusCodeContext := StatusCodeContext(http.StatusInternalServerError)
		r = r.WithContext(withStatusCodeCtx(r.Context(), &statusCodeContext))

		next.ServeHTTP(w, r)
	}
}

package main

import (
	"net/http"
)

func (app *app) internalServer(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err)

	app.writeJSONErr(w, http.StatusInternalServerError, "internal server error")
}

func (app *app) badRequest(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Warnw("bad request", "method", r.Method, "path", r.URL.Path, "error", err)

	app.writeJSONErr(w, http.StatusBadRequest, err.Error())
}

func (app *app) notFound(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Warnw("resource not found", "method", r.Method, "path", r.URL.Path, "error", err)

	app.writeJSONErr(w, http.StatusNotFound, "not found")
}

func (app *app) unauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err)

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)

	app.writeJSONErr(w, http.StatusConflict, "unauthorized error")
}

func (app *app) unauthorizedErrorOthers(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	app.writeJSONErr(w, http.StatusConflict, "unauthorized error")
}

func (app *app) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnf("forbidden error", "method", r.Method, "path", r.URL.Path, "error")

	app.writeJSONErr(w, http.StatusForbidden, "forbidden error")
}

func (app *app) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	app.writeJSONErr(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}

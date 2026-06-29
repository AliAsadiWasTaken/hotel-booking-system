package router

import "net/http"

func New() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthCheck)

	return mux
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte("OK"))
}

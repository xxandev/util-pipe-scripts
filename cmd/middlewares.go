package main

import (
	"net/http"
	"util-pipe/internal/dbg"
	"util-pipe/internal/xj"
)

func mwAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !authorization.IsActive() {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		if err := authorization.CheckCredentials(r.BasicAuth()); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(xj.Errf("auth - %v", err).Bts())
			return
		}
		next.ServeHTTP(w, r)
	})
}

func mwWhereAmI(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbg.Log.Println("call url path:", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

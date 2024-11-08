package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_routes_exits(t *testing.T) {
	testApp := Config{}

	testRoutes := testApp.routes()
	chiRoutes, ok := testRoutes.(chi.Router)
	if !ok {
		t.Fatalf("expected chi.Mux, got %T", testRoutes)
	}
	routes := []string{"/authenticate"}
	
	for _, route := range routes {
		routeExists(t, chiRoutes, route)
	}
}

func routeExists(t *testing.T, routes chi.Router, route string) {
	found := false

	_ = chi.Walk(routes, func(method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if route == foundRoute {
			found = true
		}
		return nil
	})

	if !found {
		t.Errorf("did not find %s in registered routes", route)
	}
}



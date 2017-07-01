package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestHTTP(t *testing.T) {
	server := httptest.NewServer(router())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("Root", func(t *testing.T) {
		e.GET("/").
			Expect().
			Status(http.StatusOK).Text().Equal("welcome")
	})

	t.Run("Plus", func(t *testing.T) {
		e.POST("/plus").
			WithJSON(map[string]float64{"a": 2, "b": 2}).
			Expect().
			Status(http.StatusOK).JSON().Equal(map[string]float64{"answer": 4})
	})

	t.Run("Minus", func(t *testing.T) {
		e.POST("/minus").
			WithJSON(map[string]float64{"a": 2, "b": 2}).
			Expect().
			Status(http.StatusOK).JSON().Equal(map[string]float64{"answer": 0})
	})

	t.Run("Mul", func(t *testing.T) {
		os.Setenv("PLUS_SVC_URL", server.URL+"/plus")
		e.POST("/mul").
			WithJSON(map[string]float64{"a": 3, "b": 4}).
			Expect().
			Status(http.StatusOK).JSON().Equal(map[string]float64{"answer": 12})
	})
}

package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"encoding/json"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Printf("Listening on :%s", port)
	http.ListenAndServe(":"+port, router())
}

func router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Post("/plus", func(w http.ResponseWriter, r *http.Request) {
		payload, err := parsePayload(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := map[string]float64{
			"answer": payload["a"] + payload["b"],
		}
		fmt.Printf("plus at=info %f + %f = %f\n", payload["a"], payload["b"], result["answer"])
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	r.Post("/minus", func(w http.ResponseWriter, r *http.Request) {
		payload, err := parsePayload(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := map[string]float64{
			"answer": payload["a"] - payload["b"],
		}
		fmt.Printf("minus at=info %f - %f = %f\n", payload["a"], payload["b"], result["answer"])
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	r.Post("/mul", func(w http.ResponseWriter, r *http.Request) {
		payload, err := parsePayload(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		answer := payload["b"]
		for i := 1; i < int(payload["a"]); i++ {
			answer, err = doPlus(answer, payload["b"], getHdrs(&r.Header))
			if err != nil {
				fmt.Printf("mul at=error %s\n", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		result := map[string]float64{
			"answer": answer,
		}
		fmt.Printf("mul at=info %f * %f = %f\n", payload["a"], payload["b"], result["answer"])
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	return r
}

func parsePayload(body io.ReadCloser) (map[string]float64, error) {
	var payload map[string]float64
	err := json.NewDecoder(body).Decode(&payload)
	body.Close()
	if err != nil {
		return nil, fmt.Errorf("Bad payload, expecting {\"a\": 1.0, \"b\": 2.0}")
	}
	if _, ok := payload["a"]; !ok {
		return nil, fmt.Errorf("Bad payload, expecting {\"a\": 1.0, \"b\": 2.0}")
	}
	if _, ok := payload["b"]; !ok {
		return nil, fmt.Errorf("Bad payload, expecting {\"a\": 1.0, \"b\": 2.0}")
	}

	return payload, nil
}

func doPlus(a, b float64, traceHdrs map[string]string) (float64, error) {
	payload, err := json.Marshal(map[string]float64{"a": a, "b": b})
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", os.Getenv("PLUS_SVC_URL"), bytes.NewReader(payload))
	if err != nil {
		return 0, err
	}
	req.Header.Add("Content-Type", "application/json")
	for k, v := range traceHdrs {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	var result map[string]float64
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()

	return result["answer"], nil
}

func getHdrs(hdrs *http.Header) map[string]string {
	r := map[string]string{}
	s := []string{"x-request-id", "x-ot-span-context", "x-b3-traceid", "x-b3-spanid", "x-b3-parentspanid", "x-b3-sampled", "x-b3-flags"}
	for _, ss := range s {
		r[ss] = hdrs.Get(ss)
	}
	return r
}

// Mock Disassembly.AI API for local demos of the Terraform provider.
// Implements POST /v1/scans and GET /v1/scans/{id} with canned findings.
//
//	go run ./mockserver        # listens on :8080
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type finding struct {
	RuleID   string `json:"rule_id"`
	Title    string `json:"title"`
	Severity string `json:"severity"`
}

type scan struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	Findings  []finding `json:"findings"`
	SARIF     string    `json:"sarif"`
	ReportURL string    `json:"report_url"`
	Tokens    int64     `json:"tokens"`
	Runs      int64     `json:"runs"`
}

func sarif(target string) string {
	doc := map[string]any{
		"$schema": "https://json.schemastore.org/sarif-2.1.0.json",
		"version": "2.1.0",
		"runs": []any{map[string]any{
			"tool": map[string]any{"driver": map[string]any{"name": "disassembly"}},
			"results": []any{
				map[string]any{"ruleId": "authz.idor", "level": "error",
					"message": map[string]any{"text": "IDOR on GET /api/v1/orders/{id} at " + target}},
				map[string]any{"ruleId": "secrets.jwt", "level": "warning",
					"message": map[string]any{"text": "Weak JWT secret in .env.example"}},
			},
		}},
	}
	b, _ := json.Marshal(doc)
	return string(b)
}

func result(id, target string) scan {
	return scan{
		ID:     id,
		Status: "completed",
		Findings: []finding{
			{RuleID: "authz.idor", Title: "IDOR on GET /api/v1/orders/{id}", Severity: "error"},
			{RuleID: "secrets.jwt", Title: "Weak JWT secret in .env.example", Severity: "warning"},
		},
		SARIF:     sarif(target),
		ReportURL: "https://app.disassembly.ai/scans/" + id,
		Tokens:    182000,
		Runs:      1,
	}
}

func main() {
	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}
	// naive in-memory store: id -> target
	store := map[string]string{}

	auth := func(w http.ResponseWriter, r *http.Request) bool {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, `{"error":"missing bearer token"}`, http.StatusUnauthorized)
			return false
		}
		return true
	}

	http.HandleFunc("/v1/scans", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || !auth(w, r) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}
		var body struct{ Target, Effort string }
		_ = json.NewDecoder(r.Body).Decode(&body)
		id := fmt.Sprintf("scan_%06d", len(store)+1)
		store[id] = body.Target
		log.Printf("start scan %s target=%s effort=%s", id, body.Target, body.Effort)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(scan{ID: id, Status: "running"})
	})

	http.HandleFunc("/v1/scans/", func(w http.ResponseWriter, r *http.Request) {
		if !auth(w, r) {
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/v1/scans/")
		target := store[id]
		if target == "" {
			target = "unknown"
		}
		log.Printf("get scan %s -> completed", id)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(result(id, target))
	})

	log.Printf("mock Disassembly API on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

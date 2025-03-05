package testing

import (
	"calc_parallel/internal/application"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_Calc(t *testing.T) {
	app := application.NewOrchestrator()
	handler_calc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Handler_Calc(w, r)
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(`{"expression": "(1+2)*3"}`))
	w := httptest.NewRecorder()
	app.Handler_Calc(w, req)
	res := w.Result()
	defer res.Body.Close()
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler_calc.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Errorf("Ожидался %d, полученный статус %d", http.StatusCreated, rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("Ошибка при декодировании: %v", err)
	}
	if id, ok := resp["id"]; !ok || id == "" {
		t.Errorf("Неправильный id: %v", resp)
	}
}

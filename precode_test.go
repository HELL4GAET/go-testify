package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandleValidRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cafe?count=2&city=moscow", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Ожидался статус 200")
	assert.NotEmpty(t, rr.Body.String(), "Тело ответа не должно быть пустым")
}

func TestMainHandleWrongCity(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cafe?count=2&city=novosibirsk", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Ожидался статус 400")
	assert.Equal(t, "wrong city value", rr.Body.String())
}

func TestMainHandleCountMoreThanAvailable(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cafe?count=100&city=moscow", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, "Ожидался статус 200")

	expected := strings.Join(cafeList["moscow"], ",")
	assert.Equal(t, expected, rr.Body.String(), "Ответ должен содержать все доступные кафе")

	cafes := strings.Split(rr.Body.String(), ",")
	assert.Len(t, cafes, len(cafeList["moscow"]), "Количество возвращённых кафе должно совпадать с доступным")
}

package httpresourcemux

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golangsimple/httpurlparthandler"
)

func TestNewMux(t *testing.T) {
	tt := []struct {
		name         string
		method       string
		target       string
		body         io.Reader
		responseCode int
		responseBody string
	}{
		{
			name:         "/ 404",
			method:       http.MethodGet,
			target:       "/",
			responseCode: http.StatusNotFound,
			responseBody: "404 page not found\n",
		},
		{
			name:         "/api/ 404",
			method:       http.MethodGet,
			target:       "/api/",
			responseCode: http.StatusNotFound,
			responseBody: "404 page not found\n",
		},
		{
			name:         "List /api/tasks/ OK",
			method:       http.MethodGet,
			target:       "/api/tasks/",
			responseCode: http.StatusOK,
			responseBody: `[]`,
		},
		{
			name:         "Get /api/tasks/1000 OK",
			method:       http.MethodGet,
			target:       "/api/tasks/1000",
			responseCode: http.StatusOK,
			responseBody: `[{"ID":1000}]`,
		},
		{
			name:         "Create /api/tasks Created",
			method:       http.MethodPost,
			target:       "/api/tasks/",
			responseCode: http.StatusCreated,
			responseBody: `{"ID":1000}`,
		},
		{
			name:         "Update /api/tasks/1000 Accepted",
			method:       http.MethodPut,
			target:       "/api/tasks/1000",
			responseCode: http.StatusAccepted,
			responseBody: `{"ID":1000}`,
		},
		{
			name:         "Delete /api/tasks/1000 OK",
			method:       http.MethodDelete,
			target:       "/api/tasks/1000",
			responseCode: http.StatusOK,
			responseBody: ``,
		},
	}

	handler := SetupHandler()

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest(test.method, test.target, nil))
			if response.Code != test.responseCode {
				t.Errorf("Incorrect response code %v", response.Code)
			}
			if string(response.Body.Bytes()) != test.responseBody {
				t.Errorf("Incorrect response body %v", string(response.Body.Bytes()))
			}
		})
	}
}

func SetupHandler() http.Handler {
	mux := http.NewServeMux()

	tasksHandler := httpurlparthandler.NewHandler("/api/", "tasks/", func(taskID string) http.Handler {
		return NewMux(taskID, nil, ResourceHandlers{
			ListHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("[]"))
			}),
			CreateHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{"ID":1000}`))
			}),
			ReadHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[{"ID":1000}]`))
			}),
			UpdateHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte(`{"ID":1000}`))
			}),
			DeleteHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(``))
			}),
		})
	})
	mux.Handle(tasksHandler.Route, tasksHandler)

	return mux
}

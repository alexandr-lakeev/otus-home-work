package internalhttp

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app/calendar/usecase"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/infrastructure/logger"
	memorystorage "github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/infrastructure/storage/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createTestServer() (*Server, error) {
	logger, err := logger.New(config.LoggerConf{
		Env:   "dev",
		Level: "INFO",
	})
	if err != nil {
		return nil, err
	}

	usecase := usecase.New(memorystorage.New(), logger)

	return NewServer(config.ServerConf{
		BindAddress: ":8080",
	}, usecase, logger), nil
}

// nolint
func TestServer_CreateEvent(t *testing.T) {
	t.Run("create and get event", func(t *testing.T) {
		createRequest := `
			{
				"title":"Dummy event",
				"description":"Event description",
				"date":"2021-12-31T23:59:59Z",
				"duration":30
			}
		`

		userID := uuid.New()

		rec := httptest.NewRecorder()
		req, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			"/api/calendar/v1/events",
			strings.NewReader(createRequest),
		)
		req.Header.Add("x-user-id", userID.String())

		server, err := createTestServer()
		require.NoError(t, err)

		server.server.Handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Result().StatusCode)

		response := struct {
			Data   map[string]interface{} `json:"data"`
			Errors map[string]string      `json:"errors"`
		}{}

		err = json.NewDecoder(rec.Body).Decode(&response)

		require.NoError(t, err)
		require.Empty(t, response.Errors)
		require.Contains(t, response.Data, "id")

		id, err := uuid.Parse(response.Data["id"].(string))

		require.NoError(t, err)

		rec = httptest.NewRecorder()
		req, _ = http.NewRequestWithContext(
			context.Background(),
			http.MethodGet,
			"/api/calendar/v1/events/"+id.String(),
			nil,
		)
		req.Header.Add("x-user-id", userID.String())

		server.server.Handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Result().StatusCode)

		response.Data = map[string]interface{}{}
		response.Errors = map[string]string{}
		err = json.NewDecoder(rec.Body).Decode(&response)

		require.NoError(t, err)
		require.Equal(t, "Dummy event", response.Data["title"])
		require.Equal(t, "Event description", response.Data["description"])
		require.Equal(t, "2021-12-31T23:59:59Z", response.Data["date"])
		require.Equal(t, float64(30), response.Data["duration"])
	})

	t.Run("bad request", func(t *testing.T) {
		postRequest := `{}`
		userID := uuid.New()

		rec := httptest.NewRecorder()
		req, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			"/api/calendar/v1/events",
			strings.NewReader(postRequest),
		)
		req.Header.Add("x-user-id", userID.String())

		server, err := createTestServer()
		require.NoError(t, err)

		server.server.Handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)

		response := struct {
			Data   map[string]string `json:"data"`
			Errors map[string]string `json:"errors"`
		}{}

		err = json.NewDecoder(rec.Body).Decode(&response)

		require.NoError(t, err)
		require.Contains(t, response.Errors, "Title")
		require.Contains(t, response.Errors, "Date")
		require.Contains(t, response.Errors, "Duration")
	})
}

// nolint
func TestServer_UpdateEvent(t *testing.T) {
	t.Run("update event", func(t *testing.T) {
		createRequest := `
			{
				"title":"Dummy event",
				"description":"Event description",
				"date":"2021-12-31T23:59:59Z",
				"duration":30
			}
		`

		userID := uuid.New()

		rec := httptest.NewRecorder()
		req, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			"/api/calendar/v1/events",
			strings.NewReader(createRequest),
		)
		req.Header.Add("x-user-id", userID.String())

		server, err := createTestServer()
		require.NoError(t, err)

		server.server.Handler.ServeHTTP(rec, req)

		result := rec.Result()
		require.Equal(t, http.StatusOK, result.StatusCode)

		response := struct {
			Data   map[string]interface{} `json:"data"`
			Errors map[string]string      `json:"errors"`
		}{}

		err = json.NewDecoder(rec.Body).Decode(&response)

		require.NoError(t, err)
		require.Empty(t, response.Errors)
		require.Contains(t, response.Data, "id")

		id, err := uuid.Parse(response.Data["id"].(string))

		log.Println("ID:", id, "/api/calendar/v1/events/"+id.String())

		require.NoError(t, err)

		updateRequest := `
			{
				"title":"Cool event",
				"description":"",
				"date":"2021-01-01T00:00:01Z",
				"duration":180
			}
		`

		rec = httptest.NewRecorder()
		req, _ = http.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			"/api/calendar/v1/events/"+id.String(),
			strings.NewReader(updateRequest),
		)
		req.Header.Add("x-user-id", userID.String())

		server.server.Handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNoContent, rec.Result().StatusCode)

		rec = httptest.NewRecorder()
		req, _ = http.NewRequestWithContext(
			context.Background(),
			http.MethodGet,
			"/api/calendar/v1/events/"+id.String(),
			nil,
		)
		req.Header.Add("x-user-id", userID.String())

		server.server.Handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Result().StatusCode)

		response.Data = map[string]interface{}{}
		response.Errors = map[string]string{}
		err = json.NewDecoder(rec.Body).Decode(&response)

		require.NoError(t, err)
		require.Equal(t, "Cool event", response.Data["title"])
		require.Equal(t, "", response.Data["description"])
		require.Equal(t, "2021-01-01T00:00:01Z", response.Data["date"])
		require.Equal(t, float64(180), response.Data["duration"])
	})

	t.Run("bad request", func(t *testing.T) {
		postRequest := `{}`
		id := uuid.New()
		userID := uuid.New()

		rec := httptest.NewRecorder()
		req, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			"/api/calendar/v1/events/"+id.String(),
			strings.NewReader(postRequest),
		)
		req.Header.Add("x-user-id", userID.String())

		server, err := createTestServer()
		require.NoError(t, err)

		server.server.Handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)

		response := struct {
			Data   map[string]string `json:"data"`
			Errors map[string]string `json:"errors"`
		}{}

		err = json.NewDecoder(rec.Body).Decode(&response)

		require.NoError(t, err)
		require.Contains(t, response.Errors, "Title")
		require.Contains(t, response.Errors, "Date")
		require.Contains(t, response.Errors, "Duration")
	})
}

// nolint
func TestServer_Errors(t *testing.T) {
	t.Run("event not found", func(t *testing.T) {
		id := uuid.New()
		userID := uuid.New()

		rec := httptest.NewRecorder()
		req, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodGet,
			"/api/calendar/v1/events/"+id.String(),
			nil,
		)
		req.Header.Add("x-user-id", userID.String())

		server, err := createTestServer()
		require.NoError(t, err)

		server.server.Handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNotFound, rec.Result().StatusCode)
	})

	t.Run("unauthorized", func(t *testing.T) {
		id := uuid.New()

		rec := httptest.NewRecorder()
		req, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodGet,
			"/api/calendar/v1/events/"+id.String(),
			nil,
		)

		server, err := createTestServer()
		require.NoError(t, err)

		server.server.Handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Result().StatusCode)
	})
}

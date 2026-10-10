package apiserver

import (
	"http-rest-api/internal/app/store/sqlstore"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_HandleUsersCreate(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	databaseURL := "host=localhost dbname=restapi_test sslmode=disable"
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := newServer(sqlstore.New(db))
	s.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusOK)
}

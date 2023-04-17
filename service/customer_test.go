package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	log *logrus.Logger
	cs  *CustomerService
)

func TestMain(m *testing.M) {
	code := 0
	defer func() {
		os.Exit(code)
	}()

	log = logrus.New()
	log.Formatter = &logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
		ForceColors:     true,
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.WithError(err).Fatal("Could not connect to docker")
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	src := map[string]string{
		"user":     "postgres",
		"password": "password",
		"db":       "integra_test",
	}

	runOpts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12-alpine",
		Env: []string{
			"POSTGRES_USER=" + src["user"],
			"POSTGRES_PASSWORD=" + src["password"],
			"POSTGRES_DB=" + src["db"],
		},
	}
	resource, err := pool.RunWithOptions(&runOpts,
		func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		})

	if err != nil {
		log.WithError(err).Fatal("could not start postgres container")
	}

	defer func() {
		err = pool.Purge(resource)
		if err != nil {
			log.WithError(err).Error("Could not purge resource")
		}
	}()

	// Tell docker to hard kill the container in 120 seconds
	if err := resource.Expire(120); err != nil {
		log.WithError(err)
	}

	logWaiter, err := pool.Client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container:    resource.Container.ID,
		OutputStream: log.Writer(),
		ErrorStream:  log.Writer(),
		Stderr:       true,
		Stdout:       true,
		Stream:       true,
	})
	if err != nil {
		log.WithError(err).Fatal("could not connect to postgres container log output")
	}
	defer func() {
		err = logWaiter.Close()
		if err != nil {
			log.WithError(err).Error("Could not wait for container log to close")
		}
	}()

	pool.MaxWait = 120 * time.Second
	link := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", src["user"], src["password"], resource.GetHostPort("5432/tcp"), src["db"])
	if err = pool.Retry(func() error {
		db, err := sql.Open("postgres", link)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.WithError(err).Fatal("Could not connect to postgres server")
	}

	var dbErr error
	cs, dbErr = NewCustomerServices(
		WithPGXConfiguration(logrus.New(), link),
	)
	if dbErr != nil {
		log.WithError(dbErr).Fatal("could not connect postgres container")
	}

	code = m.Run()
	// run this after all test has run
}

func TestCreateUser(t *testing.T) {
	testCase := []struct {
		name     string
		testData string
		message  string
		code     int
		response map[string]any
	}{
		{
			name: "successful response",
			testData: `{
					"userName": "willi",
					"firstName": "john",
					"lastName": "peter",
					"email": "john@gmaily.com",
					"department": "computer",
					"userStatus": 2
				}`,
			message:  "successful",
			response: make(map[string]any),
			code:     http.StatusCreated,
		},
		{
			name: "email required error response",
			testData: `{
					"userName": "williss",
					"firstName": "john",
					"lastName": "peter",
					"email": "",
					"department": "computer",
					"userStatus": 2
				}`,
			message:  "failed to validation user",
			code:     http.StatusBadRequest,
			response: make(map[string]any),
		},
		{
			name: "username required error response",
			testData: `{
					"userName": "",
					"firstName": "john",
					"lastName": "peter",
					"email": "wlli@gmail.com",
					"department": "computer",
					"userStatus": 1
				}`,
			message:  "failed to validation user",
			code:     http.StatusBadRequest,
			response: make(map[string]any),
		},
		{
			name: "invalid user status error response",
			testData: `{
					"userName": "hhdd",
					"firstName": "john",
					"lastName": "peter",
					"email": "wlli@gmail.com",
					"department": "computer",
					"userStatus": 4
				}`,
			message:  "failed to validation user",
			code:     http.StatusBadRequest,
			response: make(map[string]any),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(tc.testData))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			// Assertions
			if assert.NoError(t, cs.Create(ctx)) {
				assert.Equal(t, tc.code, rec.Code)
				log.Println(rec.Body.String())
				if assert.NoError(t, json.NewDecoder(rec.Body).Decode(&tc.response)) {
					log.Println(tc.response)
					assert.Equal(t, tc.response["message"], tc.message)
				}
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	testCase := []struct {
		name    string
		message string
		code    int
	}{
		{
			name:    "successful fetch response",
			message: "successful",
			code:    http.StatusOK,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/user", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			// Assertions
			if assert.NoError(t, cs.FetchAll(ctx)) {
				assert.Equal(t, tc.code, rec.Code)
				log.Println(rec.Body.String())
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	testCase := []struct {
		name     string
		testData string
		message  string
		code     int
		response map[string]any
	}{
		{
			name: "successful response",
			testData: `{
					"id": "1",
					"userName": "new user name",
					"firstName": "john",
					"lastName": "peter",
					"email": "john@gmailyhh.com",
					"department": "computer",
					"userStatus": 1
				}`,
			message:  "successful",
			response: make(map[string]any),
			code:     http.StatusOK,
		},
		{
			name: "email required error response",
			testData: `{
					"userName": "williss",
					"id": "1",
					"firstName": "john",
					"lastName": "peter",
					"email": "",
					"department": "computer",
					"userStatus": 2
				}`,
			message:  "failed to validation user",
			code:     http.StatusBadRequest,
			response: make(map[string]any),
		},
		{
			name: "username required error response",
			testData: `{
					"userName": "",
					"id": "1",
					"firstName": "john",
					"lastName": "peter",
					"email": "wlli@gmail.com",
					"department": "computer",
					"userStatus": 1
				}`,
			message:  "failed to validation user",
			code:     http.StatusBadRequest,
			response: make(map[string]any),
		},
		{
			name: "invalid user status error response",
			testData: `{
					"userName": "hhdd",
					"id": "1",
					"firstName": "john",
					"lastName": "peter",
					"email": "wlli@gmail.com",
					"department": "computer",
					"userStatus": 4
				}`,
			message:  "failed to validation user",
			code:     http.StatusBadRequest,
			response: make(map[string]any),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/user", strings.NewReader(tc.testData))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			// Assertions
			if assert.NoError(t, cs.Update(ctx)) {
				assert.Equal(t, tc.code, rec.Code)
				log.Println(rec.Body.String())
				if assert.NoError(t, json.NewDecoder(rec.Body).Decode(&tc.response)) {
					log.Println(tc.response)
					assert.Equal(t, tc.response["message"], tc.message)
				}
			}
		})
	}
}

package api

import (
	"encoding/json"
	"fmt"
	"github.com/S-Ryouta/sample-blog/models"
	"github.com/S-Ryouta/sample-blog/serializers/entities"
	"github.com/S-Ryouta/sample-blog/test/factories"
	"github.com/S-Ryouta/sample-blog/test/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	DBConn *gorm.DB
)

func dbSetUp() {
	testDb := helpers.InitDb()
	testDb.Setup()
	DBConn = testDb.Connect(true)
	helpers.TableMigrate(DBConn)
}

func dbCleanUp() {
	testDb := helpers.InitDb()
	defer testDb.CleanUp()
}

func TestGetEntries(t *testing.T) {
	dbSetUp()
	helpers.TestDbMock(t)

	app := fiber.New()
	app.Get("/entities", func(c *fiber.Ctx) error {
		return GetEntries(c)
	})

	var entries []models.Entry
	entries = append(entries, factories.CreateEntry(DBConn))

	entryJson, err := json.Marshal(entities.IndexSerializer(entries))
	if err != nil {
		fmt.Println(err)
		return
	}
	expectedBody := string(entryJson)

	tests := []struct {
		description   string
		route         string
		expectedError bool
		expectedCode  int
		expectedBody  string
	}{
		{
			description:   "index route",
			route:         "/entities",
			expectedError: false,
			expectedCode:  200,
			expectedBody:  expectedBody,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(
			http.MethodGet,
			test.route,
			nil,
		)

		res, err := app.Test(req, -1)
		assert.Equalf(t, test.expectedError, err != nil, test.description)

		if test.expectedError {
			continue
		}

		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		body, err := ioutil.ReadAll(res.Body)
		assert.Equalf(t, test.expectedBody, string(body), test.description)
		assert.Nilf(t, err, test.description)
	}
	t.Cleanup(func() {
		dbCleanUp()
	})
}

func TestGetEntry(t *testing.T) {
	dbSetUp()
	helpers.TestDbMock(t)

	app := fiber.New()
	app.Get("/entities/:id", func(c *fiber.Ctx) error {
		return GetEntry(c)
	})

	var entry models.Entry
	entry = factories.CreateEntry(DBConn)

	entryJson, err := json.Marshal(entry)
	if err != nil {
		fmt.Println(err)
		return
	}
	expectedBody := string(entryJson)

	tests := []struct {
		description   string
		route         string
		expectedError bool
		expectedCode  int
		expectedBody  string
	}{
		{
			description:   "index route",
			route:         fmt.Sprintf("/entities/%s", entry.ID),
			expectedError: false,
			expectedCode:  200,
			expectedBody:  expectedBody,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(
			http.MethodGet,
			test.route,
			nil,
		)

		res, err := app.Test(req, -1)
		assert.Equalf(t, test.expectedError, err != nil, test.description)

		if test.expectedError {
			continue
		}

		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		body, err := ioutil.ReadAll(res.Body)
		assert.Equalf(t, test.expectedBody, string(body), test.description)
		assert.Nilf(t, err, test.description)
	}
	t.Cleanup(func() {
		dbCleanUp()
	})
}

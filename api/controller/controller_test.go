package controller

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emamex98/academy-go-q32021/model"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type mockUseCase struct {
	MockUseCase mockusecase
}

type mockusecase interface {
	FetchContestans() ([]model.Contestant, int)
}

type mockConUseCase struct {
	MockConUseCase mockccu
}

type mockccu interface {
	FetchContestansConcurrently(class string, max int, ixw int) ([]model.Contestant, int)
}

func (m mockConUseCase) FetchContestansConcurrently(class string, max int, ixw int) ([]model.Contestant, int) {
	return tt.list, tt.errCode
}

func (m mockUseCase) FetchContestans() ([]model.Contestant, int) {
	return tt.list, tt.errCode
}

func genetareSampleListResponse() []model.Contestant {
	var mockList []model.Contestant
	mockList = append(mockList, model.Contestant{
		ID:           100,
		Contestant:   "Guadalupe Fierce",
		RealName:     "Emanuel Estrada",
		Age:          23,
		CurrentCity:  "Guadalajara",
		CurrentScore: 1000,
		Bio:          "Guadalupe Fierce is an imaginary drag queen.",
	})
	return mockList
}

func genetareSampleSingleResponse() model.Contestant {
	resp := model.Contestant{
		ID:           100,
		Contestant:   "Guadalupe Fierce",
		RealName:     "Emanuel Estrada",
		Age:          23,
		CurrentCity:  "Guadalajara",
		CurrentScore: 1000,
		Bio:          "Guadalupe Fierce is an imaginary drag queen.",
	}
	return resp
}

type mockListResponses struct {
	list    []model.Contestant
	errCode int
}

var tt mockListResponses

func TestHelloWorld(t *testing.T) {

	tCases := []struct {
		name     string
		expected string
		hasError bool
	}{
		{
			name:     "get hello world",
			expected: `{"message":"hello, world!"}`,
			hasError: false,
		},
	}

	for _, tc := range tCases {

		t.Run(tc.name, func(t *testing.T) {
			s := CreateControllers(mockUseCase{}, mockConUseCase{})

			req := httptest.NewRequest(http.MethodGet, "/api", nil)
			w := httptest.NewRecorder()
			s.HelloWorld(w, req)

			res := w.Result()
			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			if tc.hasError {
				assert.NotEqual(t, 0, err)
			} else {
				assert.EqualValues(t, tc.expected, string(data))
			}

		})
	}
}

func TestGetContestants(t *testing.T) {

	tCases := []struct {
		name     string
		expect   string
		resp     []model.Contestant
		err      int
		hasError bool
	}{
		{
			name:     "get contestants list",
			expect:   `[{"ID":100,"Contestant":"Guadalupe Fierce","Real Name":"Emanuel Estrada","Age":23,"Current City":"Guadalajara","Score":1000,"Bio":"Guadalupe Fierce is an imaginary drag queen."}]`,
			resp:     genetareSampleListResponse(),
			err:      0,
			hasError: false,
		},
		{
			name:     "failed call to usecase",
			expect:   `{"error":"500"}`,
			resp:     nil,
			err:      500,
			hasError: true,
		},
		{
			name:     "bad request to usecase",
			expect:   `{"error":"400"}`,
			resp:     nil,
			err:      400,
			hasError: true,
		},
	}

	for _, tc := range tCases {

		tt.list = tc.resp
		tt.errCode = tc.err

		t.Run(tc.name, func(t *testing.T) {
			s := CreateControllers(mockUseCase{}, mockConUseCase{})

			req := httptest.NewRequest(http.MethodGet, "/api/contestants", nil)
			w := httptest.NewRecorder()
			s.GetContestans(w, req)

			res := w.Result()
			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			if tc.hasError {
				assert.NotEqual(t, 0, err)
			} else {
				assert.EqualValues(t, tc.expect, string(data))
			}

		})
	}
}

func TestGetSingleContestant(t *testing.T) {

	tCases := []struct {
		name     string
		expect   string
		id       string
		resp     []model.Contestant
		err      int
		hasError bool
	}{
		{
			name:     "get contestant",
			expect:   `{"ID":100,"Contestant":"Guadalupe Fierce","Real Name":"Emanuel Estrada","Age":23,"Current City":"Guadalajara","Score":1000,"Bio":"Guadalupe Fierce is an imaginary drag queen."}`,
			id:       "100",
			resp:     genetareSampleListResponse(),
			err:      0,
			hasError: false,
		},
		{
			name:     "failed call to usecase",
			expect:   `{"error":"500"}`,
			id:       "100",
			resp:     nil,
			err:      500,
			hasError: true,
		},
		{
			name:     "bad request to usecase",
			expect:   `{"error":"400"}`,
			id:       "100",
			resp:     nil,
			err:      400,
			hasError: true,
		},
		{
			name:     "bad request to endpoint",
			expect:   `{"error":"400"}`,
			id:       "notanid",
			resp:     nil,
			err:      400,
			hasError: true,
		},
		{
			name:     "id not found",
			expect:   `{"error":"404"}`,
			id:       "300",
			resp:     nil,
			err:      404,
			hasError: true,
		},
	}

	for _, tc := range tCases {

		tt.list = tc.resp
		tt.errCode = tc.err

		t.Run(tc.name, func(t *testing.T) {
			s := CreateControllers(mockUseCase{}, mockConUseCase{})

			req := httptest.NewRequest(http.MethodGet, "/api/contestants/100", nil)
			w := httptest.NewRecorder()

			vars := map[string]string{
				"id": tc.id,
			}
			req = mux.SetURLVars(req, vars)

			s.GetSingleContestant(w, req)

			res := w.Result()
			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			if tc.hasError {
				assert.NotEqual(t, 0, err)
			} else {
				assert.EqualValues(t, tc.expect, string(data))
			}

		})
	}
}

func TestGetContestansConcurrently(t *testing.T) {

	tCases := []struct {
		name     string
		expect   string
		resp     []model.Contestant
		err      int
		class    string
		items    string
		ixw      string
		hasError bool
	}{
		{
			name:     "get contestants list",
			expect:   `[{"ID":100,"Contestant":"Guadalupe Fierce","Real Name":"Emanuel Estrada","Age":23,"Current City":"Guadalajara","Score":1000,"Bio":"Guadalupe Fierce is an imaginary drag queen."}]`,
			resp:     genetareSampleListResponse(),
			err:      0,
			class:    "even",
			items:    "1",
			ixw:      "1",
			hasError: false,
		},
		{
			name:     "failed call to usecase",
			expect:   `{"error":"500"}`,
			resp:     nil,
			err:      500,
			class:    "even",
			items:    "1",
			ixw:      "1",
			hasError: true,
		},
		{
			name:     "bad request to usecase",
			expect:   `{"error":"400"}`,
			resp:     nil,
			err:      400,
			class:    "notcorrect",
			items:    "notanint",
			ixw:      "notanint",
			hasError: true,
		},
	}

	for _, tc := range tCases {

		tt.list = tc.resp
		tt.errCode = tc.err

		t.Run(tc.name, func(t *testing.T) {
			s := CreateControllers(mockUseCase{}, mockConUseCase{})

			req := httptest.NewRequest(http.MethodGet, "/api/contestants_concurrent", nil)
			w := httptest.NewRecorder()

			q := req.URL.Query()
			q.Add("type", tc.class)
			q.Add("items", tc.items)
			q.Add("items_per_workers", tc.ixw)
			req.URL.RawQuery = q.Encode()

			s.GetContestansConcurrently(w, req)

			res := w.Result()
			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			if tc.hasError {
				assert.NotEqual(t, 0, err)
			} else {
				assert.EqualValues(t, tc.expect, string(data))
			}

		})
	}
}

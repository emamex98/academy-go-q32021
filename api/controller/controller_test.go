package controller

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emamex98/academy-go-q32021/model"
	"github.com/stretchr/testify/assert"
)

type mockUseCase struct {
	MockUseCase mockusecase
}

type mockusecase interface {
	FetchContestans() ([]model.Contestant, int)
}

func (m mockUseCase) FetchContestans() ([]model.Contestant, int) {
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
	return mockList, 0
}

func TestHelloWorld(t *testing.T) {

	tCases := []struct {
		name     string
		response string
		err      error
		hasError bool
	}{
		{
			name:     "test 1",
			response: "hello world",
			err:      nil,
			hasError: false,
		},
	}

	for _, tc := range tCases {

		exp := `{"message":"hello, world!"}`

		t.Run(tc.name, func(t *testing.T) {
			s := CreateControllers(mockUseCase{})

			req := httptest.NewRequest(http.MethodGet, "/api", nil)
			w := httptest.NewRecorder()
			s.HelloWorld(w, req)
			res := w.Result()
			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			if tc.hasError {
				assert.NotEqual(t, 0, err)
			} else {
				assert.EqualValues(t, exp, string(data))
			}

		})
	}
}

func TestGetContestans(t *testing.T) {

	tCases := []struct {
		name     string
		response string
		err      error
		hasError bool
	}{
		{
			name:     "get contestants list",
			response: "",
			err:      nil,
			hasError: false,
		},
	}

	for _, tc := range tCases {

		exp := `[{"ID":100,"Contestant":"Guadalupe Fierce","Real Name":"Emanuel Estrada","Age":23,"Current City":"Guadalajara","Score":1000,"Bio":"Guadalupe Fierce is an imaginary drag queen."}]`

		t.Run(tc.name, func(t *testing.T) {
			s := CreateControllers(mockUseCase{})

			req := httptest.NewRequest(http.MethodGet, "/api/contestants", nil)
			w := httptest.NewRecorder()
			s.GetContestans(w, req)
			res := w.Result()
			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			if tc.hasError {
				assert.NotEqual(t, 0, err)
			} else {
				assert.EqualValues(t, exp, string(data))
			}

		})
	}
}

func TestGetSingleContestant(t *testing.T) {

	tCases := []struct {
		name     string
		response string
		err      error
		hasError bool
	}{
		{
			name:     "get contestant by id",
			response: "",
			err:      nil,
			hasError: false,
		},
	}

	for _, tc := range tCases {

		exp := `[{"ID":100,"Contestant":"Guadalupe Fierce","Real Name":"Emanuel Estrada","Age":23,"Current City":"Guadalajara","Score":1000,"Bio":"Guadalupe Fierce is an imaginary drag queen."}]`

		t.Run(tc.name, func(t *testing.T) {
			s := CreateControllers(mockUseCase{})

			req := httptest.NewRequest(http.MethodGet, "/api/contestants/68190", nil)
			w := httptest.NewRecorder()
			s.GetSingleContestant(w, req)
			res := w.Result()
			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			if tc.hasError {
				assert.NotEqual(t, 0, err)
			} else {
				assert.EqualValues(t, exp, string(data))
			}

		})
	}
}

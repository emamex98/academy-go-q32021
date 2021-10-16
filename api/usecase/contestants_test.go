package usecase

import (
	"strings"
	"testing"

	"github.com/emamex98/academy-go-q32021/model"
	"github.com/stretchr/testify/assert"
)

type testResp struct {
	Model   []model.Contestant
	ErrCode int
}

type mockExtApi struct {
	ExtApi mockextapi
}

type mockextapi interface {
	FetchBiosAndScores() (map[int]model.ContestantInfo, error)
}

func (m mockExtApi) FetchBiosAndScores() (map[int]model.ContestantInfo, error) {
	var testMap = make(map[int]model.ContestantInfo)
	testMap[100] = model.ContestantInfo{
		Bio:   "Guadalupe Fierce is an imaginary drag queen.",
		Score: 1000,
	}
	return testMap, nil
}

type mockCsvUtil struct {
	CsvUtil mockcsvutil
}

type mockcsvutil interface {
	ReadCSV() ([][]string, error)
	WriteCSV(records []model.Contestant) error
}

func (m mockCsvUtil) ReadCSV() ([][]string, error) {
	var arrs [][]string
	arrs = append(arrs, strings.Split("ID,Contestant,Real Name,Age,Current City,Score,Bio", ","))
	arrs = append(arrs, strings.Split("100,Guadalupe Fierce,Emanuel Estrada,23,Guadalajara,1000,Guadalupe Fierce is an imaginary drag queen.", ","))
	return arrs, nil
}

func (m mockCsvUtil) WriteCSV(records []model.Contestant) error {
	return nil
}

func TestFetchContestans(t *testing.T) {

	tCases := []struct {
		name     string
		response testResp
		err      error
		hasError bool
	}{
		{
			name:     "test 1",
			response: testResp{},
			err:      nil,
			hasError: false,
		},
	}

	for _, tc := range tCases {

		var mcon []model.Contestant
		mcon = append(mcon, model.Contestant{
			ID:           100,
			Contestant:   "Guadalupe Fierce",
			RealName:     "Emanuel Estrada",
			Age:          23,
			CurrentCity:  "Guadalajara",
			CurrentScore: 1000,
			Bio:          "Guadalupe Fierce is an imaginary drag queen.",
		})

		t.Run(tc.name, func(t *testing.T) {
			s := CreateUseCase(mockExtApi{}, mockCsvUtil{})
			contestants, err := s.FetchContestans()

			if tc.hasError {
				assert.NotEqual(t, 0, err)
			} else {
				assert.EqualValues(t, mcon, contestants)
			}

		})
	}

}

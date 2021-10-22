package usecase

import (
	"errors"
	"strings"
	"testing"

	"github.com/emamex98/academy-go-q32021/model"
	"github.com/stretchr/testify/assert"
)

type mockExtApiReturns struct {
	tmap map[int]model.ContestantInfo
	err  error
}

func makeSampleExtApiResp() map[int]model.ContestantInfo {
	var testMap = make(map[int]model.ContestantInfo)
	testMap[100] = model.ContestantInfo{
		Bio:   "Guadalupe Fierce is an imaginary drag queen.",
		Score: 1000,
	}
	return testMap
}

type testExpect struct {
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
	return tm.tmap, tm.err
}

type mockCsvUtil struct {
	CsvUtil mockcsvutil
}

type mockReadCsvReturns struct {
	lines [][]string
	err   error
}

func makeSampleCsvResp() [][]string {
	var arrs [][]string
	arrs = append(arrs, strings.Split("ID,Contestant,Real Name,Age,Current City,Score,Bio", ","))
	arrs = append(arrs, strings.Split("100,Guadalupe Fierce,Emanuel Estrada,23,Guadalajara,1000,Guadalupe Fierce is an imaginary drag queen.", ","))
	return arrs
}

type mockcsvutil interface {
	ReadCSV() ([][]string, error)
	WriteCSV(records []model.Contestant) error
}

func (m mockCsvUtil) ReadCSV() ([][]string, error) {
	return tt.lines, tt.err
}

func (m mockCsvUtil) WriteCSV(records []model.Contestant) error {
	return nil
}

var c = []model.Contestant{{
	ID:           100,
	Contestant:   "Guadalupe Fierce",
	RealName:     "Emanuel Estrada",
	Age:          23,
	CurrentCity:  "Guadalajara",
	CurrentScore: 1000,
	Bio:          "Guadalupe Fierce is an imaginary drag queen.",
}}

var tm mockExtApiReturns
var tt mockReadCsvReturns

func TestFetchContestans(t *testing.T) {

	tCases := []struct {
		name     string
		expect   testExpect
		passApi  mockExtApiReturns
		passCsv  mockReadCsvReturns
		hasError bool
	}{
		{
			name: "fetch contestants list",
			expect: testExpect{
				Model:   c,
				ErrCode: 0,
			},
			passApi: mockExtApiReturns{
				tmap: makeSampleExtApiResp(),
				err:  nil,
			},
			passCsv: mockReadCsvReturns{
				lines: makeSampleCsvResp(),
				err:   nil,
			},
			hasError: false,
		},
		{
			name: "csv read error",
			expect: testExpect{
				Model:   nil,
				ErrCode: 500,
			},
			passCsv: mockReadCsvReturns{
				lines: nil,
				err:   errors.New("error"),
			},
			hasError: true,
		},
		{
			name: "api fetch error",
			expect: testExpect{
				Model:   nil,
				ErrCode: 500,
			},
			passApi: mockExtApiReturns{
				tmap: nil,
				err:  errors.New("error"),
			},
			hasError: true,
		},
	}

	for _, tc := range tCases {

		var mcon = c

		t.Run(tc.name, func(t *testing.T) {

			tm.tmap = tc.passApi.tmap
			tm.err = tc.passApi.err

			tt.lines = tc.passCsv.lines
			tt.err = tc.passCsv.err

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

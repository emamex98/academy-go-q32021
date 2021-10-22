package usecase

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/emamex98/academy-go-q32021/model"
	"github.com/stretchr/testify/assert"
)

type mockCCUC struct {
	CsvUtil mockCsvCncUtil
}

type mockCsvCncUtil interface {
	CreateCsvReader() (*csv.Reader, error)
}

func (m mockCCUC) CreateCsvReader() (*csv.Reader, error) {
	return ttc.reader, ttc.err
}

type mockCsvReaderReturn struct {
	reader *csv.Reader
	err    error
}

func createSampleCsvReader() *csv.Reader {
	file, err := os.Create("./test.csv")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	header := strings.Split("ID,Contestant,Real Name,Age,Current City,Score,Bio", ",")
	if err := w.Write(header); err != nil {
		fmt.Println(err)
	}

	bodyEven := strings.Split("100,Guadalupe Fierce,Emanuel Estrada,23,Guadalajara,1000,Guadalupe Fierce is an imaginary drag queen.", ",")
	if err := w.Write(bodyEven); err != nil {
		fmt.Println(err)
	}

	bodyOdd := strings.Split("101,Guadalupe Fierce,Emanuel Estrada,23,Guadalajara,1000,Guadalupe Fierce is an imaginary drag queen.", ",")
	if err := w.Write(bodyOdd); err != nil {
		fmt.Println(err)
	}

	csvf, err := os.Open("./test.csv")
	if err != nil {
		fmt.Println(err)
	}

	reader := csv.NewReader(csvf)
	return reader
}

var d = []model.Contestant{{
	ID:           101,
	Contestant:   "Guadalupe Fierce",
	RealName:     "Emanuel Estrada",
	Age:          23,
	CurrentCity:  "Guadalajara",
	CurrentScore: 1000,
	Bio:          "Guadalupe Fierce is an imaginary drag queen.",
}}

var ttc mockCsvReaderReturn

func TestFetchContestansConcurrently(t *testing.T) {

	tCases := []struct {
		name     string
		expect   testExpect
		passCsv  mockCsvReaderReturn
		class    string
		items    int
		ixw      int
		hasError bool
	}{
		{
			name: "fetch contestants list concurrently even",
			expect: testExpect{
				Model:   c,
				ErrCode: 0,
			},
			passCsv: mockCsvReaderReturn{
				reader: createSampleCsvReader(),
				err:    nil,
			},
			class:    "even",
			items:    1,
			ixw:      1,
			hasError: false,
		},
		{
			name: "fetch contestants list concurrently odd",
			expect: testExpect{
				Model:   d,
				ErrCode: 0,
			},
			passCsv: mockCsvReaderReturn{
				reader: createSampleCsvReader(),
				err:    nil,
			},
			class:    "odd",
			items:    1,
			ixw:      1,
			hasError: false,
		},
		{
			name: "error creating csv reader",
			expect: testExpect{
				Model:   nil,
				ErrCode: 500,
			},
			passCsv: mockCsvReaderReturn{
				reader: nil,
				err:    errors.New("error"),
			},
			class:    "even",
			items:    1,
			ixw:      1,
			hasError: true,
		},
		{
			name: "wrong paramter for type",
			expect: testExpect{
				Model:   nil,
				ErrCode: 400,
			},
			passCsv: mockCsvReaderReturn{
				reader: nil,
				err:    errors.New("error"),
			},
			class:    "wrong",
			items:    1,
			ixw:      1,
			hasError: true,
		},
	}

	for _, tc := range tCases {

		t.Run(tc.name, func(t *testing.T) {

			ttc.reader = tc.passCsv.reader
			ttc.err = tc.passCsv.err

			s := CreateConcurrentUseCase(mockCCUC{})
			contestants, err := s.FetchContestansConcurrently(tc.class, tc.items, tc.ixw)

			if tc.hasError {
				assert.NotEqual(t, 0, err)
			} else {
				assert.EqualValues(t, tc.expect.Model, contestants)
			}

		})
	}
	os.Remove("./test.csv")
}

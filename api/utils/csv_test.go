package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/emamex98/academy-go-q32021/model"
	"github.com/stretchr/testify/assert"
)

var c = []model.Contestant{{
	ID:           100,
	Contestant:   "Guadalupe Fierce",
	RealName:     "Emanuel Estrada",
	Age:          23,
	CurrentCity:  "Guadalajara",
	CurrentScore: 1000,
	Bio:          "Guadalupe Fierce is an imaginary drag queen.",
}}

func mockReadCsvLines() [][]string {
	var arrs [][]string
	arrs = append(arrs, strings.Split("ID,Contestant,Real Name,Age,Current City,Score,Bio", ","))
	arrs = append(arrs, strings.Split("100,Guadalupe Fierce,Emanuel Estrada,23,Guadalajara,1000,Guadalupe Fierce is an imaginary drag queen.", ","))
	return arrs
}

func createSampleCsv() {
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

	body := strings.Split("100,Guadalupe Fierce,Emanuel Estrada,23,Guadalajara,1000,Guadalupe Fierce is an imaginary drag queen.", ",")
	if err := w.Write(body); err != nil {
		fmt.Println(err)
	}
}

var mockedLength = 15

func TestReadCSV(t *testing.T) {
	tCases := []struct {
		name     string
		path     string
		expected [][]string
		hasError bool
	}{
		{
			name:     "read csv",
			path:     "./test.csv",
			hasError: false,
		},
		{
			name:     "invalid csv path",
			path:     "./inexsistent.csv",
			hasError: true,
		},
	}

	for _, tc := range tCases {

		createSampleCsv()

		t.Run(tc.name, func(t *testing.T) {
			s := CreateCsvUtil(tc.path, tc.path)
			csvLines, err := s.ReadCSV()

			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.IsType(t, tc.expected, csvLines)
			}

		})
	}

	os.Remove("./test.csv")
}

func TestWriteCSV(t *testing.T) {

	tCases := []struct {
		name     string
		expected [][]string
		path     string
		hasError bool
	}{
		{
			name:     "write to csv",
			expected: mockReadCsvLines(),
			path:     "./test.csv",
			hasError: false,
		},
		{
			name:     "invalid path",
			expected: mockReadCsvLines(),
			path:     "D:/",
			hasError: true,
		},
	}

	for _, tc := range tCases {

		createSampleCsv()

		t.Run(tc.name, func(t *testing.T) {
			s := CreateCsvUtil("./test.csv", tc.path)
			err := s.WriteCSV(c)

			if tc.hasError {
				assert.Error(t, err)
			} else {

				csvf, err := os.Open(tc.path)
				if err != nil {
					fmt.Println(err)
				}

				csvLines, err := csv.NewReader(csvf).ReadAll()
				if err != nil {
					fmt.Println(err)
				}

				defer csvf.Close()
				assert.EqualValues(t, tc.expected, csvLines)
			}

		})
		os.Remove(tc.path)
	}

	os.Remove("./test.csv")
}

func TestCreateCsvReader(t *testing.T) {

	tCases := []struct {
		name     string
		expected *csv.Reader
		path     string
		hasError bool
	}{
		{
			name:     "create reader",
			expected: &csv.Reader{},
			path:     "./test.csv",
			hasError: false,
		},
		{
			name:     "inexistent file",
			expected: nil,
			path:     "nonexistent.csv",
			hasError: true,
		},
	}

	for _, tc := range tCases {

		createSampleCsv()

		s := CreateCsvUtil(tc.path, tc.path)
		res, err := s.CreateCsvReader()

		if tc.hasError {
			assert.Error(t, err)
		} else {
			assert.IsType(t, tc.expected, res)
		}

	}

	os.Remove("./test.csv")
}

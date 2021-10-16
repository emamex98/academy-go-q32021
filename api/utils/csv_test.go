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

var mockedLength = 15

func TestReadCSV(t *testing.T) {
	tCases := []struct {
		name     string
		response []model.Contestant
		err      error
		hasError bool
	}{
		{
			name:     "read csv",
			response: c,
			err:      nil,
			hasError: false,
		},
	}

	for _, tc := range tCases {

		t.Run(tc.name, func(t *testing.T) {
			s := CreateCsvUtil("../lmd.csv", "test.csv")
			csvLines, err := s.ReadCSV()

			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.EqualValues(t, mockedLength, len(csvLines))
			}

		})
	}

	os.Remove("./test.csv")
}

func TestWriteCSV(t *testing.T) {

	tCases := []struct {
		name     string
		response []model.Contestant
		err      error
		hasError bool
	}{
		{
			name:     "save contestant",
			response: c,
			err:      nil,
			hasError: false,
		},
	}

	for _, tc := range tCases {

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

		t.Run(tc.name, func(t *testing.T) {
			s := CreateCsvUtil("test.csv", "test.csv")
			err := s.WriteCSV(c)

			if tc.hasError {
				assert.Error(t, err)
			} else {
				mockedLines := mockReadCsvLines()

				csvf, err := os.Open("test.csv")
				if err != nil {
					fmt.Println(err)
				}

				csvLines, err := csv.NewReader(csvf).ReadAll()
				if err != nil {
					fmt.Println(err)
				}

				defer csvf.Close()
				assert.EqualValues(t, mockedLines, csvLines)
			}

		})
	}

	os.Remove("./test.csv")
}

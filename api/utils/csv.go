package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/emamex98/academy-go-q32021/model"
)

type contestantsUseCase struct {
	InputPath  string
	OutputPath string
}

func CreateCsvUtil(inpath string, outpath string) contestantsUseCase {
	return contestantsUseCase{
		InputPath:  inpath,
		OutputPath: outpath,
	}
}

func (c contestantsUseCase) ReadCSV() ([][]string, error) {

	csvf, err := os.Open(c.InputPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	csvLines, err := csv.NewReader(csvf).ReadAll()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer csvf.Close()
	return csvLines, nil
}

func (c contestantsUseCase) WriteCSV(records []model.Contestant) error {

	file, err := os.Create(c.OutputPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	header := strings.Split("ID,Contestant,Real Name,Age,Current City,Score,Bio", ",")
	if err := w.Write(header); err != nil {
		fmt.Println(err)
		return err
	}

	for _, record := range records {

		row := []string{
			strconv.Itoa(record.ID),
			record.Contestant,
			record.RealName,
			strconv.Itoa(record.Age),
			record.CurrentCity,
			strconv.Itoa(record.CurrentScore),
			record.Bio}

		if err := w.Write(row); err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

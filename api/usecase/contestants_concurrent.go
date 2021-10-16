package usecase

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/emamex98/academy-go-q32021/model"
)

type contestantsConcurrentUseCase struct {
	CsvUtil csvCncUtil
}

type csvCncUtil interface {
	ReadCSV() ([][]string, error)
	WriteCSV(records []model.Contestant) error
}

func CreateConcurrentUseCase(csvu csvUtil) contestantsConcurrentUseCase {
	return contestantsConcurrentUseCase{
		CsvUtil: csvu,
	}
}

func worker(workerId int, jobs <-chan []string, res chan<- model.Contestant) {
	fmt.Println(workerId, jobs, res)
	for line := range jobs {

		id, err := strconv.Atoi(line[0])
		if err != nil {
			fmt.Println("35", err)
		}

		age, err := strconv.Atoi(line[3])
		if err != nil {
			fmt.Println("40", err)
		}

		contestant := model.Contestant{
			ID:           id,
			Contestant:   line[1],
			RealName:     line[2],
			Age:          age,
			CurrentCity:  line[4],
			CurrentScore: 0,
			Bio:          "pending...",
		}
		res <- contestant
	}
}

func ReadCsvLineByLine() (*csv.Reader, error) {

	csvf, err := os.Open("./lmd.csv")
	if err != nil {
		fmt.Println("60", err)
		return nil, err
	}

	reader := csv.NewReader(csvf)

	if err != nil {
		fmt.Println("67", err)
		return nil, err
	}

	defer csvf.Close()
	return reader, nil
}

func (uc contestantsConcurrentUseCase) FetchContestansConcurrently(class string, max int, ixw int) ([]model.Contestant, int) {

	jobs := make(chan []string)
	res := make(chan model.Contestant)

	noWorkers := 1
	if max >= ixw {
		noWorkers = int(max / ixw)
	}

	for i := 1; i <= noWorkers; i++ {
		go worker(i, jobs, res)
	}

	csvReader, err := ReadCsvLineByLine()
	if err != nil {
		fmt.Println("91", err)
		return nil, -1
	}

	var Contestants []model.Contestant
	i := 0

	for {
		line, err := csvReader.Read()
		fmt.Println(len(line))

		if err == io.EOF || i == max || len(line) <= 1 {
			close(jobs)
			break
		}

		if line[0] == "id" {
			continue
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			fmt.Println("112", err)
		}

		switch class {
		case "odd":
			if id%2 == 0 {
				jobs <- line
				con := <-res
				Contestants = append(Contestants, con)
				i++
			}
		case "even":
			if id%2 != 0 {
				jobs <- line
				con := <-res
				Contestants = append(Contestants, con)
				i++
			}
		default:
			jobs <- line
			con := <-res
			Contestants = append(Contestants, con)
			i++
		}

		if err != nil {
			log.Fatal("135", err)
		}
	}

	return Contestants, 0
}

package usecase

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/emamex98/academy-go-q32021/model"
)

type contestantsConcurrentUseCase struct {
	CsvUtil csvCncUtil
}

type csvCncUtil interface {
	CreateCsvReader() (*csv.Reader, error)
}

func CreateConcurrentUseCase(csvu csvCncUtil) contestantsConcurrentUseCase {
	return contestantsConcurrentUseCase{
		CsvUtil: csvu,
	}
}

func worker(workerId int, jobs <-chan []string, res chan<- model.Contestant) {
	fmt.Println(workerId, jobs, res)
	for line := range jobs {

		id, err := strconv.Atoi(line[0])
		if err != nil {
			fmt.Println(err)
			id = 0
		}

		age, err := strconv.Atoi(line[3])
		if err != nil {
			fmt.Println(err)
			age = 0
		}

		score, err := strconv.Atoi(line[5])
		if err != nil {
			fmt.Println(err)
			score = 0
		}

		contestant := model.Contestant{
			ID:           id,
			Contestant:   line[1],
			RealName:     line[2],
			Age:          age,
			CurrentCity:  line[4],
			CurrentScore: score,
			Bio:          line[6],
		}
		res <- contestant
	}
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

	csvReader, err := uc.CsvUtil.CreateCsvReader()
	if err != nil {
		fmt.Println(err)
		return nil, 500
	}

	var Contestants []model.Contestant
	i := 0

	for {

		line, err := csvReader.Read()
		if err == io.EOF || i == max {
			close(jobs)
			break
		}

		if line[0] == "ID" {
			continue
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			fmt.Println(err)
			return nil, 400
		}

		switch class {
		case "even":
			if id%2 == 0 {
				jobs <- line
				con := <-res
				Contestants = append(Contestants, con)
				i++
			}
		case "odd":
			if id%2 != 0 {
				jobs <- line
				con := <-res
				Contestants = append(Contestants, con)
				i++
			}
		default:
			return nil, 400
		}

	}

	return Contestants, 0
}

package usecase

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"strconv"
	"sync"

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

func worker(workerId int, maxJobs int, jobs <-chan []string, res chan<- model.Contestant) {

	fmt.Println("Worker", workerId, "started.")
	jobsDone := 0

	for {

		if jobsDone == maxJobs {
			break
		}

		line, ok := <-jobs

		if !ok {
			break
		}

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
		jobsDone++
	}

	fmt.Println("Worker", workerId, "finished, processed", jobsDone, "jobs in total.")
}

func (uc contestantsConcurrentUseCase) FetchContestansConcurrently(class string, max int, ixw int) ([]model.Contestant, int) {

	var Contestants []model.Contestant

	linesChan := make(chan []string, max)
	respChan := make(chan model.Contestant, max)

	totalWorkers := 1
	if max >= ixw {
		totalWorkers = int(math.Ceil(float64(max) / float64(ixw)))
	}

	csvReader, err := uc.CsvUtil.CreateCsvReader()
	if err != nil {
		fmt.Println(err)
		return nil, 500
	}

	i := 0
	for {

		line, err := csvReader.Read()
		if err == io.EOF || i == max {
			close(linesChan)
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
				linesChan <- line
				i++
			}
		case "odd":
			if id%2 != 0 {
				linesChan <- line
				i++
			}
		default:
			return nil, 400
		}

	}

	var wg sync.WaitGroup

	for i := 1; i <= totalWorkers; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			worker(i, ixw, linesChan, respChan)
		}()
	}

	wg.Wait()
	close(respChan)

	for item := range respChan {
		Contestants = append(Contestants, item)
	}

	return Contestants, 0
}

package usecase

import (
	"fmt"
	"strconv"

	"github.com/emamex98/academy-go-q32021/model"
)

type contestantsUseCase struct {
	ExtApi  extapi
	CsvUtil csvUtil
}

type extapi interface {
	FetchBiosAndScores() (map[int]model.ContestantInfo, error)
}

type csvUtil interface {
	ReadCSV() ([][]string, error)
	WriteCSV(records []model.Contestant) error
}

func CreateUseCase(extApi extapi, csvu csvUtil) contestantsUseCase {
	return contestantsUseCase{
		ExtApi:  extApi,
		CsvUtil: csvu,
	}
}

func (uc contestantsUseCase) FetchContestans() ([]model.Contestant, int) {

	var Contestants []model.Contestant

	csvLines, err := uc.CsvUtil.ReadCSV()
	if err != nil {
		fmt.Println(err)
		return nil, 500
	}

	info, err := uc.ExtApi.FetchBiosAndScores()
	if err != nil {
		return nil, 500
	}

	for i, line := range csvLines {

		if i == 0 {
			continue
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			fmt.Println(err)
			return nil, 400
		}

		age, err := strconv.Atoi(line[3])
		if err != nil {
			fmt.Println(err)
			return nil, 400
		}

		contestant := model.Contestant{
			ID:           id,
			Contestant:   line[1],
			RealName:     line[2],
			Age:          age,
			CurrentCity:  line[4],
			CurrentScore: info[id].Score,
			Bio:          info[id].Bio,
		}

		Contestants = append(Contestants, contestant)
	}

	err = uc.CsvUtil.WriteCSV(Contestants)
	if err != nil {
		fmt.Println(err)
		return nil, 500
	}

	return Contestants, 0
}

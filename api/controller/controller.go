package controller

import (
	"fmt"
	"net/http"

	"github.com/emamex98/academy-go-q32021/model"
	"github.com/unrolled/render"
)

type controllers struct {
	UseCase    usecase
	ConUseCase conusecase
}

type usecase interface {
	FetchContestans() ([]model.Contestant, int)
}

type conusecase interface {
	FetchContestansConcurrently(class string, max int, ixw int) ([]model.Contestant, int)
}

func CreateControllers(uc usecase, cuc conusecase) controllers {
	return controllers{
		UseCase:    uc,
		ConUseCase: cuc,
	}
}

func returnError(resp *render.Render, w http.ResponseWriter, statusCode int, err error) {
	fmt.Println(err)
	resp.JSON(w, statusCode, map[string]string{"error": err.Error()})
}

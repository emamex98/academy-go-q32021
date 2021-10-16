package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/unrolled/render"
)

func (c controllers) GetContestans(w http.ResponseWriter, r *http.Request) {

	resp := render.New()

	contestants, errCode := c.UseCase.FetchContestans()
	if errCode != 0 {
		switch errCode {
		case 400:
			returnError(resp, w, errCode, errors.New("bad request"))
			return
		default:
			returnError(resp, w, errCode, errors.New("something happened while processing your request, try again"))
			return
		}
	}

	fmt.Println("Endpoint reached: /contestants")
	resp.JSON(w, http.StatusOK, contestants)
}

func (c controllers) GetSingleContestant(w http.ResponseWriter, r *http.Request) {

	resp := render.New()
	args := mux.Vars(r)

	id, err := strconv.Atoi(args["id"])
	if err != nil {
		fmt.Println(err)
		returnError(resp, w, http.StatusBadRequest, errors.New("bad request"))
		return
	}

	contestants, errCode := c.UseCase.FetchContestans()
	if errCode != 0 {
		switch errCode {
		case 400:
			returnError(resp, w, errCode, errors.New("bad request"))
			return
		case 500:
			returnError(resp, w, errCode, errors.New("something happened while processing your request, try again"))
			return
		}
	}

	fmt.Println("Endpoint reached: /contestants/" + strconv.Itoa(id))

	for i := range contestants {
		if contestants[i].ID == id {
			con := contestants[i]
			resp.JSON(w, http.StatusOK, con)
			return
		}
	}

	resp.JSON(w, http.StatusNotFound, map[string]string{"error": "id not found"})
}

func (c controllers) GetContestansConcurrently(w http.ResponseWriter, r *http.Request) {

	resp := render.New()
	query := r.URL.Query()

	class := query["type"]
	maxStr := query["items"]
	ixwStr := query["items_per_workers"]

	if len(class) == 0 {
		returnError(resp, w, http.StatusBadRequest, errors.New("bad request"))
		return
	}

	max, err := strconv.Atoi(maxStr[0])
	if err != nil {
		returnError(resp, w, http.StatusBadRequest, errors.New("bad request"))
		return
	}

	ixw, err := strconv.Atoi(ixwStr[0])
	if err != nil {
		returnError(resp, w, http.StatusBadRequest, errors.New("bad request"))
		return
	}

	contestants, errCode := c.ConUseCase.FetchContestansConcurrently(class[0], max, ixw)
	if errCode != 0 {
		switch errCode {
		case 400:
			returnError(resp, w, errCode, errors.New("bad request"))
			return
		default:
			returnError(resp, w, errCode, errors.New("something happened while processing your request, try again"))
			return
		}
	}

	fmt.Println("Endpoint reached: /contestants-concurrent")
	resp.JSON(w, http.StatusOK, contestants)
}

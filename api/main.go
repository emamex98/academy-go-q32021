package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/emamex98/academy-go-q32021/config"
	"github.com/emamex98/academy-go-q32021/controller"
	"github.com/emamex98/academy-go-q32021/extapi"
	"github.com/emamex98/academy-go-q32021/router"
	"github.com/emamex98/academy-go-q32021/usecase"
	"github.com/emamex98/academy-go-q32021/utils"
)

func main() {

	conf, err := config.ReadConfig("config.json")

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Server:", conf.Server.Address+"/api")

	csvu := utils.CreateCsvUtil(conf.CSV.Input, conf.CSV.Output)

	httpc := http.Client{}
	ac := extapi.CreateApiClient(conf.API.Host, &httpc)

	uc := usecase.CreateUseCase(ac, csvu)
	cuc := usecase.CreateConcurrentUseCase(csvu)

	c := controller.CreateControllers(uc, cuc)
	r := router.NewRouter(c)

	log.Fatal(http.ListenAndServe(conf.Server.Address, r))

}

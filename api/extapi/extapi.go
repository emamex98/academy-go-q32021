package extapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

	"github.com/emamex98/academy-go-q32021/model"
)

type extapi struct {
	Host       string
	HttpClient httpClient
}

type httpClient interface {
	Get(url string) (resp *http.Response, err error)
}

func CreateApiClient(host string, http httpClient) extapi {
	return extapi{
		Host:       host,
		HttpClient: http,
	}
}

func (e extapi) FetchBiosAndScores() (map[int]model.ContestantInfo, error) {

	resp, err := e.HttpClient.Get(e.Host)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	m := map[int]interface{}{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var ContestantsInfo = make(map[int]model.ContestantInfo)

	for _, page := range m {

		var conMap = make(map[string]interface{})
		v := reflect.ValueOf(page)

		if v.Kind() == reflect.Map {
			for _, key := range v.MapKeys() {
				str := fmt.Sprintf("%v", key.Interface())
				strct := v.MapIndex(key)
				conMap[str] = strct.Interface()
			}
		}

		id, err := strconv.Atoi(fmt.Sprintf("%v", conMap["pageid"]))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		score, err := strconv.Atoi(fmt.Sprintf("%v", conMap["score"]))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		ContestantsInfo[id] = model.ContestantInfo{
			Bio:   fmt.Sprintf("%v", conMap["extract"]),
			Score: score,
		}

	}

	return ContestantsInfo, nil
}

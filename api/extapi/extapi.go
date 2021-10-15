package extapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/emamex98/academy-go-q32021/model"
)

type extapi struct{}

func CreateApiClient() extapi {
	return extapi{}
}

func FetchInfo(idParam int) string {

	id := strconv.Itoa(idParam)

	resp, err := http.Get("https://rupaulsdragrace.fandom.com/api.php/?action=query&prop=extracts&exlimit=1&explaintext=true&pageids=" + id + "&format=json")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	plainBody := string(body)
	object := strings.Split(plainBody, "{")[4]
	extract := strings.Split(object, "\"extract\":")[1]
	bio := strings.Split(extract, "==")[0]
	trimmedBio := strings.Replace(
		strings.Replace(
			strings.Replace(bio, "\"", "", -1), "\\n", "", -1), ",", "", -1)

	return trimmedBio

}

func (e extapi) FetchBiosAndScores(host string) (map[int]model.ContestantInfo, error) {

	resp, err := http.Get(host)
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

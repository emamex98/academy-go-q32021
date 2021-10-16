package extapi

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/emamex98/academy-go-q32021/model"
	"github.com/stretchr/testify/assert"
)

type testResp struct {
	resp map[int]model.ContestantInfo
	err  error
}

type mockextapi struct {
	HttpClient mhttpClient
}

type mockGetResp struct {
	resp http.Response
	err  error
}

type mhttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

func (m mockextapi) Get(url string) (resp *http.Response, err error) {
	return &tt.resp, tt.err
}

func generateSampleResponse() (resp *http.Response) {
	json := `{"100": {"extract": "Guadalupe Fierce is an imaginary drag queen.", "ns": 0, "pageid": 100, "score": 1000, "title": "Guadalupe Fierce"}}`
	body := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	return &http.Response{
		StatusCode: 200,
		Body:       body,
	}
}

func generateFaultyJsonResp() (resp *http.Response) {
	json := `{"100": "not what you expect"}`
	body := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	return &http.Response{
		StatusCode: 200,
		Body:       body,
	}
}

func generateExpectValue() map[int]model.ContestantInfo {
	var con = make(map[int]model.ContestantInfo)
	con[1000] = model.ContestantInfo{
		Bio:   "Guadalupe Fierce is an imaginary drag queen.",
		Score: 1000,
	}
	return con
}

var tt mockGetResp

func TestFetchBiosAndScores(t *testing.T) {

	tCases := []struct {
		name     string
		expect   testResp
		httpc    *http.Response
		hasError bool
	}{
		{
			name: "call api and proccess response",
			expect: testResp{
				resp: generateExpectValue(),
				err:  nil,
			},
			httpc:    generateSampleResponse(),
			hasError: false,
		},
		{
			name: "failed api response",
			expect: testResp{
				resp: nil,
				err:  errors.New("error"),
			},
			httpc:    generateSampleResponse(),
			hasError: true,
		},
		{
			name: "failed to unmarshall json",
			expect: testResp{
				resp: nil,
				err:  errors.New("error"),
			},
			httpc:    generateFaultyJsonResp(),
			hasError: true,
		},
	}

	for _, tc := range tCases {

		tt.resp = *tc.httpc
		tt.err = tc.expect.err

		var mcon = make(map[int]model.ContestantInfo)
		mcon[100] = model.ContestantInfo{
			Bio:   "Guadalupe Fierce is an imaginary drag queen.",
			Score: 1000,
		}

		t.Run(tc.name, func(t *testing.T) {
			s := CreateApiClient("", mockextapi{})
			contestants, err := s.FetchBiosAndScores()

			if tc.hasError {
				assert.Error(t, tc.expect.err, err)
			} else {
				assert.EqualValues(t, mcon, contestants)
			}

		})

	}
}

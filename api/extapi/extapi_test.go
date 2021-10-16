package extapi

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/emamex98/academy-go-q32021/model"
	"github.com/stretchr/testify/assert"
)

type testResp struct{}

type mockextapi struct {
	HttpClient mhttpClient
}

type mhttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

func (m mockextapi) Get(url string) (resp *http.Response, err error) {
	json := `{"100": {"extract": "Guadalupe Fierce is an imaginary drag queen.", "ns": 0, "pageid": 100, "score": 1000, "title": "Guadalupe Fierce"}}`
	body := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	return &http.Response{
		StatusCode: 200,
		Body:       body,
	}, nil
}

func TestFetchBiosAndScores(t *testing.T) {

	tCases := []struct {
		name     string
		response testResp
		err      error
		hasError bool
	}{
		{
			name:     "test 1",
			response: testResp{},
			err:      nil,
			hasError: false,
		},
	}

	for _, tc := range tCases {

		var mcon = make(map[int]model.ContestantInfo)
		mcon[100] = model.ContestantInfo{
			Bio:   "Guadalupe Fierce is an imaginary drag queen.",
			Score: 1000,
		}

		t.Run(tc.name, func(t *testing.T) {
			s := CreateApiClient("", mockextapi{})
			contestants, err := s.FetchBiosAndScores()

			if tc.hasError {
				assert.NotEqual(t, 0, err)
			} else {
				assert.EqualValues(t, mcon, contestants)
			}

		})

	}
}

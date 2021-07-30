package gobudins

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type Config struct {
	ClientID  string `mapstructure:"clientID"`
	APIKey    string `mapstructure:"apiKey"`
	APISecret string `mapstructure:"apiSecret"`
}

// Controller is the struct for budget insight controller
type Controller struct {
	config     Config
	httpClient *http.Client
	listeners  WebhooksListeners
}

// New create new budget insight controller
func NewController(config Config, listeners WebhooksListeners) *Controller {
	return &Controller{
		config:     config,
		httpClient: &http.Client{},
		listeners:  listeners,
	}
}

// apiErrors map api errors by errors code
var apiErrors = map[string]error{
	"invalid_client": ErrAPIInvalidClientID,
	"invalid_grant":  ErrAPIInvalidGrant,
}

var (
	ErrAPIUnhandled       = errors.New("unhandled budget insight api error")
	ErrAPIInvalidClientID = errors.New("invalid budget insight client id")
	ErrAPIInvalidGrant    = errors.New("invalid budget insight grant")
)

func (ctrl *Controller) request(method string, route string, queryParams map[string]string, requestData interface{}, token string, responseData interface{}) (err error) {
	URL := ctrl.getURL(route)

	log.Printf("%s at %s", method, URL)

	var req *http.Request
	var requestBody []byte
	if requestData != nil {
		if method != http.MethodPost && method != http.MethodPut {
			return fmt.Errorf("request data can't be sended with %s", method)
		}
		requestBody, err := json.Marshal(requestData)
		if err != nil {
			return errors.Wrap(err, "failed to marshal request data")
		}

		log.Printf("json data: %s", requestBody)

		req, err = http.NewRequest(method, URL, bytes.NewBuffer(requestBody))
		if err != nil {
			return errors.Wrap(err, "failed to create request")
		}
	} else {
		req, err = http.NewRequest(method, URL, nil)
		if err != nil {
			return errors.Wrap(err, "failed to create request")
		}
	}

	timestamp := strconv.Itoa(int(time.Now().Unix()))

	// manage signature as described at https://zabo.com/docs/#signing-requests
	sigHash := hmac.New(sha256.New, []byte(ctrl.config.APISecret))
	sigHash.Write([]byte(timestamp + URL + string(requestBody)))
	sig := hex.EncodeToString(sigHash.Sum(nil))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Zabo-Key", ctrl.config.APIKey)
	req.Header.Set("X-Zabo-Sig", sig)
	req.Header.Set("X-Zabo-Timestamp", ctrl.config.APIKey)

	q := req.URL.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}

	req.URL.RawQuery = q.Encode()

	resp, err := ctrl.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to do request")
	}
	defer resp.Body.Close()

	log.Printf("response status: %s", resp.Status)

	success := resp.StatusCode >= 200 && resp.StatusCode < 300

	if responseData != nil || !success {
		var errData ErrorResponse

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read response body")
		}

		log.Printf("response body : %s", body)

		dataToDecode := responseData
		if !success {
			dataToDecode = &errData
		}

		err = json.Unmarshal(body, dataToDecode)
		if err != nil {
			return errors.Wrap(err, "failed to read response body json")
		}

		if !success {
			// manage error data
			apiErr, errHandled := apiErrors[errData.Code]
			if !errHandled {
				apiErr = errors.Wrap(ErrAPIUnhandled, errData.Code)
			}

			return errors.Wrap(apiErr, errData.Description)
		}
	}

	return err
}

func (ctrl *Controller) getURL(route string) string {
	return fmt.Sprintf("https://%s.biapi.pro/2.0%s", ctrl.config.Domain, route)
}

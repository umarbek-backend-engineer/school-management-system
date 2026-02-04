package middlerwares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func Sanitize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// url cleaning
		sanitizePath, err := clean(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// query cleaning

		params := r.URL.Query()
		sanitizedQuery := make(map[string][]string)
		for k, v := range params {
			sanitizedKey, err := clean(k)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var sanitizedValues []string
			for _, value := range v {
				cleanValue, err := clean(value)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				sanitizedValues = append(sanitizedValues, cleanValue.(string))
			}
			sanitizedQuery[sanitizedKey.(string)] = sanitizedValues
		}

		r.URL.Path = sanitizePath.(string)
		r.URL.RawQuery = url.Values(sanitizedQuery).Encode()

		//sanitize the request body

		if r.Header.Get("Content-Type") == "application/json" {
			if r.Body != nil {
				bodyByte, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Error reading body request", http.StatusBadRequest)
					return
				}
				bodyString := strings.TrimSpace(string(bodyByte))
				//reset the request body
				r.Body = io.NopCloser(bytes.NewReader([]byte(bodyString)))
				if len(bodyString) > 0 {
					var inputData interface{}
					err = json.NewDecoder(bytes.NewReader([]byte(bodyString))).Decode(&inputData)
					if err != nil {
						http.Error(w, "Invalid payload", http.StatusBadRequest)
						return
					}

					// sanitize
					sanitizedJsonBody, err := clean(inputData)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					// marshle back
					sanitizedbody, err := json.Marshal(sanitizedJsonBody)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					r.Body = io.NopCloser(bytes.NewReader(sanitizedbody))
				}
			}

		} else if r.Header.Get("Content-Type") != ""{
			log.Println("Content-type is not application json")
			http.Error(w, "Invalid content type", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// clean sanitizes input data to provent xss attacks
func clean(data interface{}) (interface{}, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			v[key] = sanitizeValue(value)
		}
		return v, nil
	case []interface{}:
		for i, value := range v {
			v[i] = sanitizeValue(value)
		}
		return v, nil
	case string:
		return sanitizaString(v), nil
	default:
		//error
		return nil, fmt.Errorf("Unsupported type: %T", data)
	}
}

func sanitizeValue(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			v[key] = sanitizeValue(value)
		}
		return v
	case []interface{}:
		for i, value := range v {
			v[i] = sanitizeValue(value)
		}
		return v
	case string:
		return sanitizaString(v)
	default:
		return v
	}
}

func sanitizaString(value string) string {
	return bluemonday.UGCPolicy().Sanitize(value)
}

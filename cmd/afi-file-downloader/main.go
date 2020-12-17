package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

var (
	ErrInvalidDuration   = errors.New("invalid duration")
	ErrUnsupportedMethod = errors.New("unsupported method")

	NopResolver = func() {}
)

const (
	MethodUnknown = "UNKNOWN"
	MethodGet     = "GET"
	MethodHead    = "HEAD"
)

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	if val, ok := v.(float64); ok {
		*d = Duration(time.Duration(val))

		return nil
	}

	if val, ok := v.(string); ok {
		t, err := time.ParseDuration(val)
		if err != nil {
			return err
		}

		*d = Duration(t)

		return nil
	}

	return ErrInvalidDuration
}

func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

type Handler func(w io.Writer, r io.Reader) error

type PromiseStringResolver func(resolve func(string), reject func(error))

type PromiseHandlerResolver func(resolve func(Handler), reject func(error))

type Promise func(resolve func(), reject func(error))

func main() {
	var (
		catchFunc = func(err error) {
			if err == nil {
				return
			}

			resp := &struct{
				Success      bool   `json:"success"`
				ErrorMessage string `json:"error-message"`
			}{
				ErrorMessage: err.Error(),
			}

			if encodeErr := json.NewEncoder(os.Stdout).Encode(resp); encodeErr != nil {
				_, _ = fmt.Fprintf(os.Stderr, "encode error: %v \n", encodeErr)
			}
		}
		retrieveMethodFunc = func(r io.Reader) PromiseStringResolver {
			return func(resolve func(string), reject func(error)) {
				req := &struct {
					Method string `json:"method"`
				}{
					Method: MethodGet,
				}

				if err := json.NewDecoder(r).Decode(req); err != nil {
					reject(err)
				}

				resolve(req.Method)
			}
		}
		selectHandlerFunc = func(handlers map[string]Handler, method string) PromiseHandlerResolver {
			return func(resolve func(Handler), reject func(error)) {
				handler, ok := handlers[method]
				if !ok {
					reject(ErrUnsupportedMethod)
				}

				resolve(handler)
			}
		}
		executeHandlerFunc = func(h Handler, w io.Writer, r io.Reader) Promise {
			return func(resolve func(), reject func(error)) {
				if err := h(w, r); err != nil {
					reject(err)
				}

				resolve()
			}
		}

		selectHandlerPromiseResolverFunc = func(h Handler) {
			executeHandlerPromise := executeHandlerFunc(h, os.Stdout, os.Stdin)
			executeHandlerPromise(NopResolver, catchFunc)
		}
		retrieveMethodPromiseResolverFunc = func(method string) {
			handlers := map[string]Handler{
				MethodGet:  findFileHandler,
				MethodHead: getFileStatHandler,
			}

			selectHandlerPromise := selectHandlerFunc(handlers, method)
			selectHandlerPromise(selectHandlerPromiseResolverFunc, catchFunc)
		}

		retrieveMethodPromise = retrieveMethodFunc(os.Stdin)
	)

	retrieveMethodPromise(retrieveMethodPromiseResolverFunc, catchFunc)
}

func findFileHandler(w io.Writer, r io.Reader) error {
	return nil
}

func getFileStatHandler(w io.Writer, r io.Reader) error {
	return nil
}

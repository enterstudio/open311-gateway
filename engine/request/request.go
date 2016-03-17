package request

import (
	"net/http"

	"Gateway311/engine/logs"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/davecgh/go-spew/spew"
)

const (
	debugRecover = false
)

var (
	log = logs.Log
)

// Services looks up the service providers and services for the specified location.
// The URL nust contain query parameters of either:
// LatitudeV and LongitudeV, or a city name.
//
// Examples:
//  http;//xyz.com/api/services?lat=34.236144&lon=-118.604794
//  http;//xyz.com/api/services?city=san+jose
func Services(w rest.ResponseWriter, r *rest.Request) {
	if debugRecover {
		defer func() {
			if rcvr := recover(); rcvr != nil {
				rest.Error(w, rcvr.(error).Error(), http.StatusInternalServerError)
			}
		}()
	}
	response, err := processServices(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&response)
}

// Create creates a new report.
func Create(w rest.ResponseWriter, r *rest.Request) {
	log.Debug("Create request: \n%s\n", spew.Sdump(r))
	if debugRecover {
		defer func() {
			if rcvr := recover(); rcvr != nil {
				rest.Error(w, rcvr.(error).Error(), http.StatusInternalServerError)
			}
		}()
	}
	response, err := processCreate(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&response)
}

// Search searches for Reports.
func Search(w rest.ResponseWriter, r *rest.Request) {
	if debugRecover {
		defer func() {
			if rcvr := recover(); rcvr != nil {
				rest.Error(w, rcvr.(error).Error(), http.StatusInternalServerError)
			}
		}()
	}
	response, err := processSearch(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&response)
}

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	v1beta1 "k8s.io/api/admission/v1beta1"
)

func TestHandleMutateErrors(t *testing.T) {

	// returns a new server
	ts := httptest.NewServer(http.HandlerFunc(handleMutate))
	// close the server when done, smart go thingy
	defer ts.Close()

	// default GET on the handle should throw an error trying to convert from empty JSON
	// The request should work though :)
	respDefault, errDefault := http.Get(ts.URL)
	assert.NoError(t, errDefault)

	respMutate, errMutate := http.Get(ts.URL + "/mutate")
	assert.NoError(t, errMutate)

	// Read the body in the response
	_, errMutate = ioutil.ReadAll(respMutate.Body)
	// Keep open until done, smart go thingy
	defer respMutate.Body.Close()
	assert.NoError(t, errMutate)

	// Read the body in the response
	body, err := ioutil.ReadAll(respDefault.Body)
	// Keep open until done, smart go thingy
	defer respDefault.Body.Close()
	assert.NoError(t, err)

	admReview := v1beta1.AdmissionReview{}
	assert.Errorf(t, json.Unmarshal(body, &admReview), "body: %s", string(body))
	assert.Empty(t, admReview.Response)

}

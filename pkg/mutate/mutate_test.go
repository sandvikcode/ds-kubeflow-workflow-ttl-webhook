package mutate

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	v1beta1 "k8s.io/api/admission/v1beta1"
)

func TestMutateJSONWithoutTTL(t *testing.T) {
	rawJSON, err := os.Open("test_data/workflow_without_ttl.json")
	byteValue, _ := ioutil.ReadAll(rawJSON)

	response, err := Mutate([]byte(byteValue))
	if err != nil {
		t.Errorf("failed to mutate AdmissionRequest %s with error %s", string(response), err)
	}

	r := v1beta1.AdmissionReview{}
	err = json.Unmarshal(response, &r)
	assert.NoError(t, err, "failed to unmarshal with error %s", err)

	rr := r.Response
	assert.Equal(t, `[{"op":"add","path":"/spec/ttlSecondsAfterFinished","value":36000}]`, string(rr.Patch))
	assert.Contains(t, rr.AuditAnnotations, "mutateme")
}

func TestMutateJSONWithTTL(t *testing.T) {
	rawJSON, err := os.Open("test_data/workflow_with_ttl.json")
	byteValue, _ := ioutil.ReadAll(rawJSON)

	response, err := Mutate([]byte(byteValue))
	if err != nil {
		t.Errorf("failed to mutate AdmissionRequest %s with error %s", string(response), err)
	}

	r := v1beta1.AdmissionReview{}
	err = json.Unmarshal(response, &r)
	assert.NoError(t, err, "failed to unmarshal with error %s", err)

	rr := r.Response
	assert.Equal(t, `[]`, string(rr.Patch))
	assert.Contains(t, rr.AuditAnnotations, "mutateme")
}

func TestMutateJSONWithTTLToLong(t *testing.T) {
	rawJSON, err := os.Open("test_data/workflow_with_to_long_ttl.json")
	byteValue, _ := ioutil.ReadAll(rawJSON)
	response, err := Mutate([]byte(byteValue))
	if err != nil {
		t.Errorf("failed to mutate AdmissionRequest %s with error %s", string(response), err)
	}

	r := v1beta1.AdmissionReview{}
	err = json.Unmarshal(response, &r)
	assert.NoError(t, err, "failed to unmarshal with error %s", err)

	rr := r.Response
	assert.Equal(t, `[{"op":"replace","path":"/spec/ttlSecondsAfterFinished","value":36000}]`, string(rr.Patch))
	assert.Contains(t, rr.AuditAnnotations, "mutateme")
}

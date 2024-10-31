package statuses

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testStatusStruct struct {
	Status    ProcessStatus `json:"status"`
	SomeValue int           `json:"some_value"`
}

func TestStatusToJson(t *testing.T) {
	test := testStatusStruct{
		Status:    PROCESSED,
		SomeValue: 1,
	}

	data, err := json.Marshal(&test)

	assert.NoError(t, err)
	assert.JSONEq(t, string(data), `{"status":"PROCESSED","some_value":1}`)
}

func TestStatusFromJson(t *testing.T) {
	test := testStatusStruct{
		Status:    PROCESSED,
		SomeValue: 1,
	}

	var test2 testStatusStruct

	err := json.Unmarshal([]byte(`{"status":"PROCESSED","some_value":1}`), &test2)
	assert.NoError(t, err)
	assert.Equal(t, test, test2)
}

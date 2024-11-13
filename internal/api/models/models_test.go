package models

import (
	"encoding/json"
	"github.com/IgorViskov/33_loyalty/internal/appmath"
	"github.com/IgorViskov/33_loyalty/internal/domain/constants"
	"github.com/IgorViskov/33_loyalty/internal/domain/statuses"
	"github.com/govalues/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStatusToJson(t *testing.T) {
	date, _ := time.Parse(constants.LayoutDate, "01.01.2025 18:05:00 +03:00")
	test := OrdersResponse{
		Status:     statuses.PROCESSED,
		Accrual:    appmath.ToNullable(decimal.MustParse("200.89")),
		Number:     "132456789",
		UploadedAt: date,
	}

	data, err := json.Marshal(&test)

	assert.NoError(t, err)
	assert.JSONEq(t, string(data), `{"accrual":200.89,"number": "132456789","status":"PROCESSED","uploaded_at":"2025-01-01T18:05:00+03:00"}`)
}

func TestPointsCalculationResponseFromJson(t *testing.T) {

	tests := []struct {
		name     string
		json     string
		want     AccrualResponse
		positive bool
	}{
		{
			name: "Normal response positive #1",
			json: `{
					  "order": "9278923470",
					  "status": "PROCESSED",
					  "accrual": 200.89
				  }`,
			want: AccrualResponse{
				Status:  statuses.PROCESSED,
				Accrual: appmath.ToNullable(decimal.MustParse("200.89")),
				Order:   "9278923470",
			},
			positive: true,
		},
		{
			name: "Normal response positive accrual as string #2",
			json: `{
					  "order": "9278923470",
					  "status": "PROCESSED",
					  "accrual": "200.89"
				  }`,
			want: AccrualResponse{
				Status:  statuses.PROCESSED,
				Accrual: appmath.ToNullable(decimal.MustParse("200.89")),
				Order:   "9278923470",
			},
			positive: true,
		},
		{
			name: "Normal response positive accrual as int #3",
			json: `{
					  "order": "9278923470",
					  "status": "PROCESSED",
					  "accrual": 200
				  }`,
			want: AccrualResponse{
				Status:  statuses.PROCESSED,
				Accrual: appmath.ToNullable(decimal.MustParse("200")),
				Order:   "9278923470",
			},
			positive: true,
		},
		{
			name: "Invalid response positive #4",
			json: `{
					  "order": "9278923470",
					  "status": "INVALID"
				  }`,
			want: AccrualResponse{
				Status:  statuses.INVALID,
				Accrual: decimal.NullDecimal{},
				Order:   "9278923470",
			},
			positive: true,
		},
		{
			name: "Bad response negative wrong status #6",
			json: `{
					  "order": "9278923470",
					  "status": 1
				  }`,
			want:     AccrualResponse{},
			positive: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var response AccrualResponse
			err := json.Unmarshal([]byte(tt.json), &response)
			if tt.positive {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, response)
			} else {
				assert.Error(t, err)
			}

		})
	}
}

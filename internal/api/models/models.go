package models

import (
	"encoding/json"
	"fmt"
	"github.com/IgorViskov/33_loyalty/internal/apperrors"
	"github.com/IgorViskov/33_loyalty/internal/appmath"
	"github.com/IgorViskov/33_loyalty/internal/domain/statuses"
	"github.com/golang-jwt/jwt/v4"
	"github.com/govalues/decimal"
	"strconv"
	"time"
)

type AuthModel struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uint64
}

type OrdersResponse struct {
	Number     string                 `json:"number"`
	Status     statuses.ProcessStatus `json:"status"`
	Accrual    decimal.NullDecimal    `json:"accrual"`
	UploadedAt time.Time              `json:"uploaded_at"`
}

type AccrualResponse struct {
	Order   string                 `json:"order"`
	Accrual decimal.NullDecimal    `json:"accrual"`
	Status  statuses.ProcessStatus `json:"status"`
}

type BalanceResponse struct {
	Current   decimal.Decimal `json:"current"`
	Withdrawn decimal.Decimal `json:"withdrawn"`
}

type WithdrawRequest struct {
	Order string          `json:"order"`
	Sum   decimal.Decimal `json:"sum"`
}

type WithdrawResponse struct {
	Order       string          `json:"order"`
	Sum         decimal.Decimal `json:"sum"`
	ProcessedAt time.Time       `json:"processed_at"`
}

func (r *OrdersResponse) MarshalJSON() ([]byte, error) {
	middle := make(map[string]interface{})
	middle["number"] = r.Number
	middle["status"] = &r.Status
	if r.Accrual.Valid {
		acc, _ := r.Accrual.Decimal.Float64()
		middle["accrual"] = acc
	}
	middle["uploaded_at"] = r.UploadedAt.Format(time.RFC3339)
	return json.Marshal(middle)
}

func (r *AccrualResponse) MarshalJSON() ([]byte, error) {
	middle := make(map[string]interface{})
	middle["order"] = r.Order
	middle["status"] = &r.Status
	acc, _ := r.Accrual.Decimal.Float64()
	middle["accrual"] = acc
	return json.Marshal(middle)
}

func (r *BalanceResponse) MarshalJSON() ([]byte, error) {
	middle := make(map[string]interface{})
	cur, _ := r.Current.Float64()
	middle["current"] = cur
	wit, _ := r.Withdrawn.Float64()
	middle["withdrawn"] = wit
	return json.Marshal(middle)
}

func (r *WithdrawResponse) MarshalJSON() ([]byte, error) {
	middle := make(map[string]interface{})
	middle["order"] = r.Order
	sum, _ := r.Sum.Float64()
	middle["sum"] = sum
	middle["processed_at"] = r.ProcessedAt.Format(time.RFC3339)
	return json.Marshal(middle)
}

func (r *AccrualResponse) UnmarshalJSON(b []byte) error {
	var dat map[string]interface{}
	var err error
	if err = json.Unmarshal(b, &dat); err != nil {
		return err
	}

	var d interface{}
	var ok bool

	d, ok = dat["order"]
	if !ok {
		return apperrors.ErrNotValidJSON
	}
	r.Order, err = parseString(d)

	d, ok = dat["status"]
	if !ok {
		return apperrors.ErrNotValidJSON
	}
	r.Status, err = statuses.ParseProcessStatus(d)
	if err != nil {
		return err
	}

	d, ok = dat["accrual"]
	if !ok && r.Status == statuses.PROCESSED {
		return apperrors.ErrNotValidJSON
	} else if r.Status == statuses.PROCESSED {
		var accrual decimal.Decimal
		accrual, e := appmath.ParseDecimal(d)
		if e != nil {
			return e
		}
		r.Accrual = appmath.ToNullable(accrual)
	}

	return nil
}

func parseString(val interface{}) (string, error) {
	switch v := val.(type) {
	case int:
		return strconv.Itoa(v), nil
	case string:
		return v, nil
	default:
		return "0", fmt.Errorf("cannot parse int of type %T", v)
	}
}

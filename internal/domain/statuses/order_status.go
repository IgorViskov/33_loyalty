package statuses

import (
	"encoding/json"
	"github.com/IgorViskov/33_loyalty/internal/apperrors"
	"strings"
)

const (
	NEW ProcessStatus = iota
	PROCESSED
	PROCESSING
	INVALID
)

type ProcessStatus uint8

var (
	orderStatusNames = map[uint8]string{
		0: "NEW", //Для GET /api/user/orders отдаем название NEW
		1: "PROCESSED",
		2: "PROCESSING",
		3: "INVALID",
	}
	orderStatusValues = map[string]uint8{
		"REGISTERED": 0, //GET /api/orders/{number} отдает нам тот же самый статус в виде строки REGISTERED
		"PROCESSED":  1,
		"PROCESSING": 2,
		"INVALID":    3,
	}
)

func (s *ProcessStatus) String() string {
	return orderStatusNames[uint8(*s)]
}

func ParseProcessStatus(val interface{}) (ProcessStatus, error) {
	s, ok := val.(string)
	if !ok {
		return ProcessStatus(0), apperrors.ErrProcessStatusShouldBeString
	}
	s = strings.TrimSpace(strings.ToUpper(s))
	value, ok := orderStatusValues[s]
	if !ok {
		return ProcessStatus(0), apperrors.ErrNotValidJSON
	}
	return ProcessStatus(value), nil
}

func (s *ProcessStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *ProcessStatus) UnmarshalJSON(data []byte) error {
	var suits string
	var err error
	if err = json.Unmarshal(data, &suits); err != nil {
		return err
	}
	if *s, err = ParseProcessStatus(suits); err != nil {
		return err
	}
	return nil
}

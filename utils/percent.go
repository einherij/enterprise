package utils

import (
	"strconv"

	json "github.com/json-iterator/go"
)

const (
	OnePercent Percent = 0.01
)

type Percent float64

func PercentFrom[percentT int | int32 | int64 | float32 | float64](percent percentT) Percent {
	switch any(percent).(type) {
	case int, int32, int64:
		return Percent(percent) / 100.
	case float32, float64:
		return Percent(percent)
	}
	return 0.
}

func (p Percent) String() string {
	return strconv.FormatFloat(100.*float64(p), 'f', -1, 64) + "%"
}

func (p Percent) Float64() float64 {
	return float64(p)
}

// Scale returns p% imply that val is 100%
func (p Percent) Scale(val float64) float64 {
	return float64(p) * val
}

// ScaleReversed returns 100% imply that val is p%
func (p Percent) ScaleReversed(val float64) float64 {
	if p == 0. {
		return 0.
	}
	return 1. / float64(p) * val
}

type Commission Percent

func (e Commission) AddTo(value float64) float64 {
	return (100*OnePercent + Percent(e)).Scale(value)
}

func (e Commission) SubtractFrom(value float64) float64 {
	return (100*OnePercent + Percent(e)).ScaleReversed(value)
}

func (p *Percent) UnmarshalJSON(data []byte) error {
	var val int
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	*p = PercentFrom(val)
	return nil
}

func (p *Percent) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(100 * *p))
}

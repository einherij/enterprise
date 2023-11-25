package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Percent_String(t *testing.T) {
	a := assert.New(t)
	p := PercentFrom(0.34)
	str := fmt.Sprintf("%v", p)
	a.Equal("34%", str)
}

func Test_Commission(t *testing.T) {
	a := assert.New(t)
	cms := Commission(0.22)
	bidFloor := 2.5
	a.Equal(3.05, cms.AddTo(bidFloor))
	a.Equal(bidFloor, cms.SubtractFrom(3.05))
}

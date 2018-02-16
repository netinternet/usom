package Usom

import (
	"testing"
)

func TestScandaily(t *testing.T) {
	masks := []string{"89.43.28.0/22", "89.43.26.0/22"}
	list := Scandaily(masks, 1)
	for _, v := range list {
		t.Log(v)
	}
}

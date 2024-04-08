package main

import (
	"github.com/moutend/go-hook/pkg/types"
	"testing"
)

func Test_isLegalCode(t *testing.T) {
	t.Log(isLegalCode(types.VK_OEM_3))
	t.Log(isLegalCode(types.VK_2))
}

package snippet

import (
	"testing"
)

func Test_readSnippetYaml(t *testing.T) {
	sn, err := readSnippetYaml("D:\\project\\Golang\\doDash\\snippet\\mysql\\show_variables.yaml")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sn)
}

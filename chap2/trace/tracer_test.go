package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Newからの戻り値がnil")
	} else {
		tracer.Trace("こんにちは")
		if buf.String() != "こんにちは\n" {
			t.Errorf("%sという誤った文字列が返されました。", buf.String())
		}
	}
}

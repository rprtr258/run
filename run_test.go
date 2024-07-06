package run

import (
	"context"
	"testing"
)

var res = ""

func newCloser(s string) func() {
	return func() {
		res += s
	}
}

func newA(lc *Lifecycle) {
	// component A uses nothing
	res += "<A>"
	lc.Add(newCloser("</A>"))
}

func newB(lc *Lifecycle) {
	// component B uses nothing
	res += "<B>"
	lc.Add(newCloser("</B>"))
}

func newC(lc *Lifecycle) {
	// component C uses component B
	newB(lc)
	// use B

	res += "<C>"
	lc.Add(newCloser("</C>"))
}

func useLC(t *testing.T, f func(lc *Lifecycle)) {
	lc := New()
	f(lc)
	if err := lc.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestLifecycle(t *testing.T) {
	useLC(t, func(lc *Lifecycle) {
		// program uses components A and C
		newA(lc)
		newC(lc)
	})
	// so we close C first, then A
	// before closing C though, we close B first
	// thus result closing order must be CBA
	if res != "<A><B><C></C></B></A>" {
		t.Fatal(res)
	}
}

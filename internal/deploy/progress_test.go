package deploy

import "testing"

func TestUnmappedProgress(t *testing.T) {
	p := Progress("deploy progress output is unmapped")
	if p != ProgressUnmapped {
		t.Errorf("Expected progress to be %d but was %d", ProgressUnmapped, p)
	}
}

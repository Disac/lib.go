// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package writer

import (
	"os"
	"testing"

	"github.com/caixw/lib.go/term"
)

func TestConsole(t *testing.T) {
	c := NewConsole(os.Stderr, term.FCyan)
	c.Write([]byte("is cyan\n"))

	c.SetColor(term.FBlue)
	c.Write([]byte("is blue\n"))

	os.Stderr.WriteString("Reset\n")
}

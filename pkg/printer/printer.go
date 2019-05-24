package printer

import (
	goast "go/ast"
	"io"

	"github.com/davecgh/go-spew/spew"
)

func FprintFile(w io.Writer, file goast.File) {
	spew.Fdump(w, file)
}

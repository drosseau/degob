package degob

import (
	"errors"
	"fmt"
	"io"
)

// Gob is a more concrete representation of a gob. It has all of the found
// types and the decoded value.
type Gob struct {
	Types map[typeId]*WireType
	Value
}

// WriteTypes writes the Gob's types to Writer
func (g *Gob) WriteTypes(w io.Writer) error {
	if g.Types == nil {
		return errors.New("gob has no defined types")
	}
	for _, t := range g.Types {
		_, err := fmt.Fprintf(w, "// type ID: %d\n", t.Id())
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "%s\n\n", t.String())
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteValues writes the Gob's values to a writer with a given style
func (g *Gob) WriteValue(w io.Writer, sty style) error {
	if g.Value == nil {
		return errors.New("gob has no value")
	}
	_, err := fmt.Fprintf(w, "%s\n", g.Value.Display(sty))
	return err
}

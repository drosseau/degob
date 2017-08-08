package degob

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// Gob is a more concrete representation of a gob, but it is still lower level.
type Gob struct {
	Types map[typeId]*WireType
	Value Value
}

// WriteTypes writes the Gob's types to Writer
func (g *Gob) WriteTypes(w io.Writer) error {
	_, err := fmt.Fprintln(w, "// Types:")
	if err != nil {
		return err
	}
	for _, t := range g.Types {
		s := t.String()
		if strings.HasPrefix(s, "type struct") {
			_, err = fmt.Fprintln(w, "// Anonymous struct")
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprintf(w, "%s\n\n", s)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteValues writes the Gob's values to a writer with a given style
func (g *Gob) WriteValue(w io.Writer, sty style) error {
	if g.Value == nil {
		return errors.New("attempted to write nil Value gob")
	}
	_, err := fmt.Fprintln(w, "// Values:")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s\n", g.Value.Display(sty))
	return err
}

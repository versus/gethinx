package gethinx

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
)

func PrintVars(w io.Writer, writePre bool, vars ...interface{}) {
	if writePre {
		io.WriteString(w, "<pre>\n")
	}
	for i, v := range vars {
		fmt.Fprintf(w, "» item %d type %T:\n", i, v)
		j, err := json.MarshalIndent(v, "", "    ")
		switch {
		case err != nil:
			fmt.Fprintf(w, "error: %v", err)
		case len(j) < 3: // {}, empty struct maybe or empty string, usually mean unexported struct fields
			w.Write([]byte(html.EscapeString(fmt.Sprintf("%+v", v))))
		default:
			w.Write(j)
		}
		w.Write([]byte("\n\n"))
	}
	if writePre {
		io.WriteString(w, "</pre>\n")
	}
}

/*
type dummy struct {
	s string
	d *dummy
}
type Dummy struct {
	S string
	D *Dummy
}

func main() {
	d1 := dummy{"test", &dummy{"child", &dummy{"grand", nil}}}
	d2 := Dummy{"test", &Dummy{"child", &Dummy{"grand", nil}}}
	m := map[string]string{"stuff": "values"}
	printVars(os.Stdout, false, d1, d2, m)
}
*/

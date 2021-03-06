package envsubst

import (
	"testing"

	"github.com/google/gofuzz"
)

// test cases sourced from tldp.org
// http://www.tldp.org/LDP/abs/html/parameter-substitution.html

func TestExpand(t *testing.T) {
	var expressions = []struct {
		params map[string]string
		input  string
		output string
	}{
		// text-only
		{
			params: map[string]string{},
			input:  "abcdEFGH28ij",
			output: "abcdEFGH28ij",
		},
		// length
		{
			params: map[string]string{"var01": "abcdEFGH28ij"},
			input:  "${#var01}",
			output: "12",
		},
		// uppercase first
		{
			params: map[string]string{"var01": "abcdEFGH28ij"},
			input:  "${var01^}",
			output: "AbcdEFGH28ij",
		},
		// uppercase
		{
			params: map[string]string{"var01": "abcdEFGH28ij"},
			input:  "${var01^^}",
			output: "ABCDEFGH28IJ",
		},
		// lowercase first
		{
			params: map[string]string{"var01": "ABCDEFGH28IJ"},
			input:  "${var01,}",
			output: "aBCDEFGH28IJ",
		},
		// lowercase
		{
			params: map[string]string{"var01": "ABCDEFGH28IJ"},
			input:  "${var01,,}",
			output: "abcdefgh28ij",
		},
		// substring with position
		{
			params: map[string]string{"path_name": "/home/bozo/ideas/thoughts.for.today"},
			input:  "${path_name:11}",
			output: "ideas/thoughts.for.today",
		},
		// substring with position and length
		{
			params: map[string]string{"path_name": "/home/bozo/ideas/thoughts.for.today"},
			input:  "${path_name:11:5}",
			output: "ideas",
		},
		// default not used
		{
			params: map[string]string{"var": "abc"},
			input:  "${var=abc}",
			output: "abc",
		},
		// default used
		{
			params: map[string]string{},
			input:  "${var=xyz}",
			output: "xyz",
		},
		{
			params: map[string]string{},
			input:  "${var:=xyz}",
			output: "xyz",
		},
		// replace suffix
		{
			params: map[string]string{"stringZ": "abcABC123ABCabc"},
			input:  "${stringZ/%abc/XYZ}",
			output: "abcABC123ABCXYZ",
		},
		// replace prefix
		{
			params: map[string]string{"stringZ": "abcABC123ABCabc"},
			input:  "${stringZ/#abc/XYZ}",
			output: "XYZABC123ABCabc",
		},
		// replace all
		{
			params: map[string]string{"stringZ": "abcABC123ABCabc"},
			input:  "${stringZ//abc/xyz}",
			output: "xyzABC123ABCxyz",
		},
		// replace first
		{
			params: map[string]string{"stringZ": "abcABC123ABCabc"},
			input:  "${stringZ/abc/xyz}",
			output: "xyzABC123ABCabc",
		},
		// // delete shortest match prefix
		// {
		// 	params: map[string]string{"filename": "bash.string.txt"},
		// 	input:  "${filename#*.}",
		// 	output: "txt",
		// },
		// // delete longest match prefix
		// {
		// 	params: map[string]string{"filename": "bash.string.txt"},
		// 	input:  "${filename#*.}",
		// 	output: "string.txt",
		// },
		// // delete shortest match suffix
		// {
		// 	params: map[string]string{"filename": "bash.string.txt"},
		// 	input:  "${filename%.*}",
		// 	output: "bash.string",
		// },
		// // delete longest match suffix
		// {
		// 	params: map[string]string{"filename": "bash.string.txt"},
		// 	input:  "${filename%%.*}",
		// 	output: "bash",
		// },

		// nested parameters
		{
			params: map[string]string{"var01": "abcdEFGH28ij"},
			input:  "${var=${var01^^}}",
			output: "ABCDEFGH28IJ",
		},
	}

	for _, expr := range expressions {
		output, err := Eval(expr.input, func(s string) string {
			return expr.params[s]
		})
		if err != nil {
			t.Errorf("Want %q expanded but got error %q", expr.input, err)
		}

		if output != expr.output {
			t.Errorf("Want %q expanded to %q, got %q",
				expr.input,
				expr.output,
				output)
		}
	}
}

func TestFuzz(t *testing.T) {
	noopFunc := func(s string) string {
		return ""
	}

	f := fuzz.New()
	for i := 0; i < 1000; i++ {
		var in string
		f.Fuzz(&in)
		_, err := Eval(in, noopFunc)
		if err != nil {
			t.Errorf("got error parsing input %q", in)
		}
	}
}

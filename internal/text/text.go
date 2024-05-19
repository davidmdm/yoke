package text

import (
	"bytes"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
	"gopkg.in/yaml.v3"

	"github.com/davidmdm/ansi"
)

type DiffFunc func(expected, actual File, context int) string

type File struct {
	Name    string
	Content string
}

func Diff(expected, actual File, context int) string {
	diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(expected.Content),
		B:        difflib.SplitLines(actual.Content),
		FromFile: expected.Name,
		ToFile:   actual.Name,
		Context:  context,
	})
	return diff
}

func DiffColorized(expected, actual File, context int) string {
	return colorize(Diff(expected, actual, context))
}

var (
	green = ansi.MakeStyle(ansi.FgGreen)
	red   = ansi.MakeStyle(ansi.FgRed)
)

func colorize(value string) string {
	lines := strings.Split(value, "\n")
	colorized := make([]string, len(lines))
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		switch line[0] {
		case '-':
			colorized[i] = green.Sprint(line)
		case '+':
			colorized[i] = red.Sprint(line)
		default:
			colorized[i] = line
		}
	}

	return strings.Join(colorized, "\n")
}

func ToYamlFile(name string, value any) (File, error) {
	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)
	err := encoder.Encode(value)
	return File{Name: name, Content: buffer.String()}, err
}

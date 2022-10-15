package pkg

import "strings"

type Fields struct {
	fields map[string]bool
}

func ParseFields(fields string) Fields {
	f := Fields{fields: map[string]bool{}}

	parts := strings.Split(fields, ",")
	for _, part := range parts {
		part = strings.Trim(part, " \r\n")
		f.fields[part] = true
	}

	return f
}

func (f *Fields) Has(field string) bool {
	return f.fields[field]
}

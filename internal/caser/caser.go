package caser

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var initialisms = []string{"ACL", "API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SQL", "SSH", "SVG", "TCP", "TLS", "TTL", "UDP", "UI", "UID", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XMLNS", "XMPP", "XSRF", "XSS"}

func init() {
	for _, a := range initialisms {
		strcase.ConfigureAcronym(a, a)
		strcase.ConfigureAcronym(cases.Title(language.Und).String(a), a)
		strcase.ConfigureAcronym(strings.ToLower(a), a)
	}
}

func GoPascal(input string) string {
	s := input

	toReplace := []string{"-", ":", "/", "."}
	sep := "_"
	for _, old := range toReplace {
		s = strings.ReplaceAll(s, old, sep)
	}

	// Format each part to match acronyms
	parts := strings.Split(s, sep)

	if len(parts) == 0 {
		panic(fmt.Sprintf("no parts for input: %s", input))
	}

	for i := range parts {
		if parts[i] == "" {
			continue
		} else {
			parts[i] = strcase.ToCamel(parts[i])
		}
		// Check and replace acronyms
		for _, initialism := range initialisms {
			if strings.HasPrefix(parts[i], cases.Title(language.Und).String(initialism)) {
				result := initialism + parts[i][len(initialism):]
				parts[i] = result
			}
		}
	}

	out := strings.Join(parts, "")

	// Check for unformatted acronyms in the output
	for _, a := range initialisms {
		wrong := cases.Title(language.Und).String(a)
		if strings.Contains(out, wrong) {
			panic(fmt.Sprintf("error in goPascal func: wrong acronym formatting: %s -> %s -> %s", input, s, out))
		}
	}

	return out
}

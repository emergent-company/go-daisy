package ui

import "strings"

// disclosureChevronClass derives the Tailwind "open" rotate class from a group
// name. A named group "group/cap" becomes "group-open/cap:rotate-180"; the
// plain "group" becomes "group-open:rotate-180". Inserting "-open" after
// "group" preserves the "/name" suffix.
func disclosureChevronClass(group string) string {
	return strings.Replace(group, "group", "group-open", 1) + ":rotate-180"
}

package schemaform

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Walk generates ordered form fields from a JSON Schema object subset.
//
// Supported schema keywords: properties, required, type, title, description, enum, items.
// Unsupported keywords are silently ignored. No validation against a specific
// JSON Schema draft.
//
// schema: decoded JSON Schema object (map[string]any)
// defaults: default values keyed by path
// values: current submitted values (override defaults on re-render)
// allowList: dotted paths → AllowMode (managed or disabled)
func Walk(schema, defaults, values map[string]any, allowList map[string]AllowMode) []Field {
	props, _ := schema["properties"].(map[string]any)
	if props == nil {
		return nil
	}
	required := requiredSet(schema)

	keys := sortedKeys(props)
	fields := make([]Field, 0, len(keys))
	for _, key := range keys {
		prop, _ := props[key].(map[string]any)
		if prop == nil {
			continue
		}
		field := buildField(key, "", prop, defaults, values, allowList, required)
		if field != nil {
			fields = append(fields, *field)
		}
	}
	return fields
}

// FallbackFromDefaults infers form fields from default values when no JSON
// Schema is available. Supports string, float64, int, bool, []any, and
// map[string]any defaults.
func FallbackFromDefaults(defaults, values map[string]any, allowList map[string]AllowMode) []Field {
	keys := sortedKeys(defaults)
	fields := make([]Field, 0, len(keys))
	for _, key := range keys {
		field := fieldFromDefault(key, "", defaults[key], values, allowList)
		if field != nil {
			fields = append(fields, *field)
		}
	}
	return fields
}

// PruneDisabled removes paths marked AllowDisabled from values before save.
func PruneDisabled(values map[string]any, allowList map[string]AllowMode) {
	for path, mode := range allowList {
		if mode == AllowDisabled {
			delete(values, path)
		}
	}
}

func buildField(key, prefix string, prop, defaults, values map[string]any, allowList map[string]AllowMode, required map[string]bool) *Field {
	path := key
	if prefix != "" {
		path = prefix + "." + key
	}
	mode := allowList[path]
	if mode == AllowDisabled {
		return nil
	}

	kind, enumVals := propKind(prop)
	managed := mode == AllowManaged

	def := valToString(lookup(defaults, key))
	val := def
	if v, ok := lookupV(values, key); ok {
		val = valToString(v)
	}

	label := propLabel(prop, key)
	helper := prop["description"]

	f := &Field{
		Path:       path,
		Label:      label,
		HelperText: fmt.Sprint(helper),
		Kind:       kind,
		Managed:    managed,
		Default:    def,
		Value:      val,
		Options:    enumVals,
	}

	if kind == KindObject {
		nested, _ := prop["properties"].(map[string]any)
		if nested != nil {
			nk := sortedKeys(nested)
			for _, k := range nk {
				np, _ := nested[k].(map[string]any)
				if np != nil {
					child := buildField(k, path, np, nestedDefaults(defaults, key), nestedDefaults(values, key), allowList, nil)
					if child != nil {
						f.Children = append(f.Children, *child)
					}
				}
			}
		}
	}

	if _, isReq := required[key]; isReq {
		if f.HelperText == "" {
			f.HelperText = "Required"
		}
	}

	return f
}

func fieldFromDefault(key, prefix string, defVal any, values map[string]any, allowList map[string]AllowMode) *Field {
	path := key
	if prefix != "" {
		path = prefix + "." + key
	}
	mode := allowList[path]
	if mode == AllowDisabled {
		return nil
	}

	kind := KindUnknown
	var opts []string
	var children []Field
	def := valToString(defVal)
	val := def
	if v, ok := lookupV(values, key); ok {
		val = valToString(v)
	}

	switch v := defVal.(type) {
	case string:
		kind = KindString
	case float64:
		if v == float64(int64(v)) {
			kind = KindInteger
		} else {
			kind = KindNumber
		}
	case int:
		kind = KindInteger
	case int64:
		kind = KindInteger
	case bool:
		kind = KindBoolean
	case []any:
		kind = KindArray
		var strOpts []string
		for _, item := range v {
			strOpts = append(strOpts, valToString(item))
		}
		opts = strOpts
	case []string:
		kind = KindArray
		opts = v
	case map[string]any:
		kind = KindObject
		nk := sortedKeys(v)
		for _, k := range nk {
			child := fieldFromDefault(k, path, v[k], values, allowList)
			if child != nil {
				children = append(children, *child)
			}
		}
	}

	f := &Field{
		Path:     path,
		Label:    humanize(key),
		Kind:     kind,
		Managed:  mode == AllowManaged,
		Default:  def,
		Value:    val,
		Options:  opts,
		Children: children,
	}
	return f
}

func propKind(prop map[string]any) (Kind, []string) {
	if enum, ok := prop["enum"].([]any); ok {
		strs := make([]string, len(enum))
		for i, v := range enum {
			strs[i] = fmt.Sprint(v)
		}
		return KindEnum, strs
	}

	t, _ := prop["type"].(string)
	switch t {
	case "boolean":
		return KindBoolean, nil
	case "integer":
		return KindInteger, nil
	case "number":
		return KindNumber, nil
	case "object":
		return KindObject, nil
	case "array":
		return KindArray, nil
	default:
		return KindString, nil
	}
}

func propLabel(prop map[string]any, key string) string {
	if t, ok := prop["title"].(string); ok && t != "" {
		return t
	}
	return humanize(key)
}

func requiredSet(schema map[string]any) map[string]bool {
	set := make(map[string]bool)
	if req, ok := schema["required"].([]any); ok {
		for _, r := range req {
			if s, ok := r.(string); ok {
				set[s] = true
			}
		}
	}
	return set
}

func nestedDefaults(m map[string]any, key string) map[string]any {
	if v, ok := m[key].(map[string]any); ok {
		return v
	}
	return nil
}

func lookup(m map[string]any, key string) any {
	if m == nil {
		return nil
	}
	return m[key]
}

func lookupV(m map[string]any, key string) (any, bool) {
	if m == nil {
		return nil, false
	}
	v, ok := m[key]
	return v, ok
}

func valToString(v any) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case float64:
		if x == float64(int64(x)) {
			return strconv.FormatInt(int64(x), 10)
		}
		return strconv.FormatFloat(x, 'f', -1, 64)
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	case bool:
		return strconv.FormatBool(x)
	default:
		return fmt.Sprint(v)
	}
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func humanize(s string) string {
	if s == "" {
		return ""
	}
	r := strings.NewReplacer("_", " ", "-", " ", ".", " ")
	s = r.Replace(s)
	return strings.ToUpper(s[:1]) + s[1:]
}

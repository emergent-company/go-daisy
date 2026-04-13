// boundarytoken generates galleryruntime.DesignToken slices from WithBoundary
// functions that carry a structured doc-comment annotation.
//
// Usage (typically via go:generate):
//
//	go run github.com/emergent-company/go-daisy/cmd/boundarytoken \
//	    -pkg github.com/emergent-company/go-daisy/components/ui \
//	    -out components/ui/boundary_tokens_gen.go \
//	    components/ui/boundary.go
//
// Annotation format (on the WithBoundary func's doc comment):
//
//	// gallery:token param1,param2,param3
//
// Type mapping rules:
//   - named string type with package consts  → TokenTypeSelect (options auto-populated)
//   - plain string                           → TokenTypeText
//   - bool                                   → TokenTypeSelect {false/true}
//   - int                                    → TokenTypeRange (default min=0,max=100,step=1)
//   - []T slice                              → 3× TokenTypeText item label tokens
//
// Override int range or slice count with extra annotations:
//
//	// gallery:hint rows:range(2,10,1)
//	// gallery:hint items:slice(3)
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

func main() {
	pkgImport := flag.String("pkg", "", "import path of the package being parsed (for the generated file header)")
	outFile := flag.String("out", "", "output file path (default: <input dir>/boundary_tokens_gen.go)")
	outPkg := flag.String("out-pkg", "", "package name for the output file (default: same as input package)")
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("usage: boundarytoken [flags] <boundary.go>")
	}
	inFile := flag.Arg(0)
	if *outFile == "" {
		*outFile = filepath.Join(filepath.Dir(inFile), "boundary_tokens_gen.go")
	}

	fset := token.NewFileSet()
	// Parse the target boundary.go
	f, err := parser.ParseFile(fset, inFile, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("parse %s: %v", inFile, err)
	}

	// Also parse the whole package dir to collect const declarations for named types.
	pkgDir := filepath.Dir(inFile)
	pkgFiles, err := parser.ParseDir(fset, pkgDir, nil, 0)
	if err != nil {
		log.Fatalf("parse package dir %s: %v", pkgDir, err)
	}
	// Collect all const values grouped by type name: map[typeName][]constEntry
	typeConsts := collectTypeConsts(pkgFiles)

	pkgName := f.Name.Name
	if *outPkg != "" {
		pkgName = *outPkg
	}

	var entries []funcEntry
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name == nil {
			continue
		}
		if !strings.HasSuffix(fn.Name.Name, "WithBoundary") {
			continue
		}
		if fn.Doc == nil {
			continue
		}
		tokens, hints := parseAnnotations(fn.Doc.Text())
		if tokens == nil {
			continue // no gallery:token annotation
		}

		params := extractParams(fn)
		entry := funcEntry{
			FuncName:  fn.Name.Name,
			TokenFunc: strings.TrimSuffix(fn.Name.Name, "WithBoundary") + "Tokens",
			Tokens:    buildTokens(tokens, hints, params, typeConsts),
		}
		entries = append(entries, entry)
	}

	if len(entries) == 0 {
		fmt.Fprintf(os.Stderr, "boundarytoken: no gallery:token annotations found in %s\n", inFile)
		os.Exit(0)
	}

	var importPkg string
	if *pkgImport != "" {
		importPkg = *pkgImport
	}

	out, err := render(pkgName, importPkg, entries)
	if err != nil {
		log.Fatalf("render: %v", err)
	}
	if err := os.WriteFile(*outFile, out, 0644); err != nil {
		log.Fatalf("write %s: %v", *outFile, err)
	}
	fmt.Printf("boundarytoken: wrote %s (%d functions)\n", *outFile, len(entries))
}

// ── annotation parsing ────────────────────────────────────────────────────────

// parseAnnotations extracts the list of token param names and any hints from
// the doc comment text. Returns (nil, nil) when no gallery:token line is found.
func parseAnnotations(doc string) (tokens []string, hints map[string]hintVal) {
	hints = map[string]hintVal{}
	for _, line := range strings.Split(doc, "\n") {
		line = strings.TrimSpace(line)
		if rest, ok := strings.CutPrefix(line, "gallery:token "); ok {
			for _, t := range strings.Split(rest, ",") {
				t = strings.TrimSpace(t)
				if t != "" {
					tokens = append(tokens, t)
				}
			}
		}
		if rest, ok := strings.CutPrefix(line, "gallery:hint "); ok {
			name, spec, found := strings.Cut(rest, ":")
			if !found {
				continue
			}
			name = strings.TrimSpace(name)
			spec = strings.TrimSpace(spec)
			existing := hints[name]
			if strings.HasPrefix(spec, "range(") {
				inner := strings.TrimPrefix(strings.TrimSuffix(spec, ")"), "range(")
				parts := strings.Split(inner, ",")
				if len(parts) == 3 {
					min, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
					max, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
					step, _ := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
					// Preserve any existing defaultVal
					existing.kind = "range"
					existing.min = min
					existing.max = max
					existing.step = step
					hints[name] = existing
				}
			} else if strings.HasPrefix(spec, "slice(") {
				inner := strings.TrimPrefix(strings.TrimSuffix(spec, ")"), "slice(")
				n, _ := strconv.Atoi(strings.TrimSpace(inner))
				existing.kind = "slice"
				existing.sliceN = n
				hints[name] = existing
			} else if strings.HasPrefix(spec, "default(") {
				inner := strings.TrimPrefix(strings.TrimSuffix(spec, ")"), "default(")
				// strip surrounding quotes if present
				inner = strings.Trim(inner, `"`)
				// Merge: preserve existing kind/range/slice, only set defaultVal
				existing.defaultVal = inner
				hints[name] = existing
			}
		}
	}
	if tokens == nil {
		return nil, nil
	}
	return tokens, hints
}

type hintVal struct {
	kind           string // "range" | "slice" | "default"
	min, max, step float64
	sliceN         int
	defaultVal     string
}

// ── param extraction ──────────────────────────────────────────────────────────

type paramInfo struct {
	name     string
	typeName string // resolved simple name: "string","bool","int","[]ActionMenuItem","ButtonVariant"…
}

func extractParams(fn *ast.FuncDecl) map[string]paramInfo {
	result := map[string]paramInfo{}
	if fn.Type.Params == nil {
		return result
	}
	for _, field := range fn.Type.Params.List {
		typeName := exprString(field.Type)
		for _, ident := range field.Names {
			result[ident.Name] = paramInfo{name: ident.Name, typeName: typeName}
		}
	}
	return result
}

func exprString(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.ArrayType:
		return "[]" + exprString(t.Elt)
	case *ast.StarExpr:
		return "*" + exprString(t.X)
	case *ast.SelectorExpr:
		return exprString(t.X) + "." + t.Sel.Name
	case *ast.Ellipsis:
		return "..." + exprString(t.Elt)
	default:
		return fmt.Sprintf("%T", e)
	}
}

// ── const collection ──────────────────────────────────────────────────────────

type constEntry struct {
	name  string // const identifier, e.g. "ButtonPrimary"
	value string // string value, e.g. "btn-primary"
}

func collectTypeConsts(pkgFiles map[string]*ast.Package) map[string][]constEntry {
	result := map[string][]constEntry{}
	for _, pkg := range pkgFiles {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				gd, ok := decl.(*ast.GenDecl)
				if !ok || gd.Tok != token.CONST {
					continue
				}
				for _, spec := range gd.Specs {
					vs, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					typeName := ""
					if vs.Type != nil {
						typeName = exprString(vs.Type)
					}
					if typeName == "" {
						continue
					}
					for i, name := range vs.Names {
						if i >= len(vs.Values) {
							continue
						}
						lit, ok := vs.Values[i].(*ast.BasicLit)
						if !ok || lit.Kind != token.STRING {
							continue
						}
						val, _ := strconv.Unquote(lit.Value)
						result[typeName] = append(result[typeName], constEntry{name: name.Name, value: val})
					}
				}
			}
		}
	}
	return result
}

// ── token building ────────────────────────────────────────────────────────────

type tokenDef struct {
	Label      string
	Group      string
	Type       string // "TokenTypeText" | "TokenTypeSelect" | "TokenTypeRange"
	Default    string
	QueryParam string
	// range
	Min, Max, Step float64
	// select
	Options []optionDef
}

type optionDef struct {
	Value string
	Label string
}

func buildTokens(tokenParams []string, hints map[string]hintVal, params map[string]paramInfo, typeConsts map[string][]constEntry) []tokenDef {
	var result []tokenDef
	for _, pname := range tokenParams {
		info, ok := params[pname]
		if !ok {
			// param not in signature — annotated but wrong name, skip
			continue
		}
		hint := hints[pname]
		result = append(result, infer(pname, info.typeName, hint, typeConsts)...)
	}
	return result
}

func infer(pname, typeName string, hint hintVal, typeConsts map[string][]constEntry) []tokenDef {
	label := camelToLabel(pname)

	// slice → N text tokens for item labels
	if strings.HasPrefix(typeName, "[]") || typeName == "..."+strings.TrimPrefix(typeName, "[]") {
		n := 3
		if hint.kind == "slice" && hint.sliceN > 0 {
			n = hint.sliceN
		}
		var toks []tokenDef
		for i := 1; i <= n; i++ {
			qp := fmt.Sprintf("%s%d", pname, i)
			toks = append(toks, tokenDef{
				Label:      fmt.Sprintf("%s %d", label, i),
				Group:      "Component",
				Type:       "galleryruntime.TokenTypeText",
				Default:    fmt.Sprintf("Item %d", i),
				QueryParam: qp,
			})
		}
		return toks
	}

	switch typeName {
	case "bool":
		return []tokenDef{{
			Label:      label,
			Group:      "Component",
			Type:       "galleryruntime.TokenTypeSelect",
			Default:    "false",
			QueryParam: pname,
			Options: []optionDef{
				{Value: "false", Label: "No"},
				{Value: "true", Label: "Yes"},
			},
		}}

	case "int":
		min, max, step := 0.0, 100.0, 1.0
		if hint.kind == "range" {
			min, max, step = hint.min, hint.max, hint.step
		}
		def := fmt.Sprintf("%g", min)
		if hint.defaultVal != "" {
			def = hint.defaultVal
		}
		return []tokenDef{{
			Label:      label,
			Group:      "Component",
			Type:       "galleryruntime.TokenTypeRange",
			Default:    def,
			QueryParam: pname,
			Min:        min,
			Max:        max,
			Step:       step,
		}}

	case "string":
		def := ""
		if hint.defaultVal != "" {
			def = hint.defaultVal
		}
		return []tokenDef{{
			Label:      label,
			Group:      "Component",
			Type:       "galleryruntime.TokenTypeText",
			Default:    def,
			QueryParam: pname,
		}}

	default:
		// named type — look up consts
		consts := typeConsts[typeName]
		if len(consts) == 0 {
			// unknown named type, fall back to text
			fallbackDef := ""
			if hint.defaultVal != "" {
				fallbackDef = hint.defaultVal
			}
			return []tokenDef{{
				Label:      label,
				Group:      "Component",
				Type:       "galleryruntime.TokenTypeText",
				Default:    fallbackDef,
				QueryParam: pname,
			}}
		}
		// Sort consts deterministically by name
		sort.Slice(consts, func(i, j int) bool { return consts[i].name < consts[j].name })
		var opts []optionDef
		for _, c := range consts {
			opts = append(opts, optionDef{
				Value: c.value,
				Label: constLabel(c.name, typeName),
			})
		}
		// Pick default: explicit hint wins, then natural keyword priority, then first option.
		def := ""
		if len(consts) > 0 {
			def = consts[0].value
		}
		if hint.defaultVal != "" {
			def = hint.defaultVal
			goto foundDefault
		}
		for _, priority := range []string{"Primary", "Default", "MD", "Medium", "Info"} {
			for _, c := range consts {
				if strings.Contains(c.name, priority) {
					def = c.value
					goto foundDefault
				}
			}
		}
	foundDefault:
		return []tokenDef{{
			Label:      label,
			Group:      "Component",
			Type:       "galleryruntime.TokenTypeSelect",
			Default:    def,
			QueryParam: pname,
			Options:    opts,
		}}
	}
}

// camelToLabel converts a camelCase name to a Title Case label.
// e.g. "currentPage" → "Current Page", "loading" → "Loading", "typ" → "Type"
func camelToLabel(s string) string {
	return titleCase(s)
}

func titleCase(s string) string {
	if s == "typ" {
		return "Type"
	}
	var buf strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		isUpper := r >= 'A' && r <= 'Z'
		if i == 0 {
			if r >= 'a' && r <= 'z' {
				buf.WriteRune(r - 32)
			} else {
				buf.WriteRune(r)
			}
			continue
		}
		if isUpper {
			// peek ahead: if next char is lowercase, this starts a new word
			nextLower := i+1 < len(runes) && runes[i+1] >= 'a' && runes[i+1] <= 'z'
			prevLower := runes[i-1] >= 'a' && runes[i-1] <= 'z'
			if prevLower || nextLower {
				buf.WriteByte(' ')
			}
			buf.WriteRune(r)
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

// constLabel strips the type prefix from a const name for display.
// e.g. "ButtonPrimary" with type "ButtonVariant" → "Primary"
//
//	"AlertSuccess" with type "AlertType" → "Success"
func constLabel(constName, typeName string) string {
	// First try stripping the full type name as a prefix (e.g. "ButtonType" from "ButtonTypeButton" → "Button").
	if stripped := strings.TrimPrefix(constName, typeName); stripped != "" && stripped != constName {
		return stripped
	}
	// Fall back: strip the base part of the type name (e.g. "Button" from "ButtonVariant")
	// by removing common type-name suffixes first.
	prefix := strings.TrimSuffix(typeName, "Type")
	prefix = strings.TrimSuffix(prefix, "Variant")
	prefix = strings.TrimSuffix(prefix, "Size")
	prefix = strings.TrimSuffix(prefix, "Intent")
	prefix = strings.TrimSuffix(prefix, "Kind")
	prefix = strings.TrimSuffix(prefix, "Shape")
	stripped := strings.TrimPrefix(constName, prefix)
	if stripped == "" || stripped == constName {
		return constName
	}
	return stripped
}

// ── code generation ───────────────────────────────────────────────────────────

type funcEntry struct {
	FuncName  string
	TokenFunc string
	Tokens    []tokenDef
}

const tmplSrc = `// Code generated by boundarytoken. DO NOT EDIT.
// Source: {{ .Pkg }}

package {{ .PkgName }}

import "github.com/emergent-company/go-daisy/galleryruntime"
{{ range .Entries }}
// {{ .TokenFunc }} returns the DesignToken slice for {{ .FuncName }}.
func {{ .TokenFunc }}() []galleryruntime.DesignToken {
	return []galleryruntime.DesignToken{
		{{- range .Tokens }}
		{
			Label:      {{ printf "%q" .Label }},
			Group:      "Component",
			Type:       {{ .Type }},
			Default:    {{ printf "%q" .Default }},
			QueryParam: {{ printf "%q" .QueryParam }},
			{{- if eq .Type "galleryruntime.TokenTypeRange" }}
			Min:        {{ printf "%g" .Min }},
			Max:        {{ printf "%g" .Max }},
			Step:       {{ printf "%g" .Step }},
			{{- end }}
			{{- if .Options }}
			Options: []galleryruntime.TokenOption{
				{{- range .Options }}
				{Value: {{ printf "%q" .Value }}, Label: {{ printf "%q" .Label }}},
				{{- end }}
			},
			{{- end }}
		},
		{{- end }}
	}
}
{{ end }}`

func render(pkgName, pkgImport string, entries []funcEntry) ([]byte, error) {
	t, err := template.New("gen").Parse(tmplSrc)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, map[string]any{
		"PkgName": pkgName,
		"Pkg":     pkgImport,
		"Entries": entries,
	}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

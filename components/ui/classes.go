package ui

// Reusable structural class strings.
//
// These replace the hand-typed inline-flex stacks that used to appear
// throughout the library (`flex items-center gap-2` and friends). They are
// plain strings so they compose with template expressions:
//
//	<div class={ Row, "justify-between" }>
//	<div class={ RowLoose, "mt-1" }>
//
// They carry no colors or typography — use the tone tokens
// (text-body … text-faint) and .eyebrow for those.
const (
	// Row is a horizontal, vertically-centred flex row with a 0.5rem gap.
	Row = "flex items-center gap-2"
	// RowLoose is Row with a 0.75rem gap.
	RowLoose = "flex items-center gap-3"
	// RowTight is Row with a 0.25rem gap.
	RowTight = "flex items-center gap-1"
	// RowBetween is Row with the children pushed to the edges.
	RowBetween = "flex items-center justify-between gap-2"
	// ColStack is a vertical flex column with a 0.5rem gap.
	ColStack = "flex flex-col gap-2"
	// ColStackLoose is ColStack with a 0.75rem gap.
	ColStackLoose = "flex flex-col gap-3"
)

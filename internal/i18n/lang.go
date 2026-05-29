// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package i18n

// Lang is a Feishu locale (e.g. "zh_cn"); "" means unset.
type Lang string

const (
	LangZhCN Lang = "zh_cn"
	LangEnUS Lang = "en_us"
	LangJaJP Lang = "ja_jp"
)

type langEntry struct {
	Code  Lang   // canonical Feishu locale
	Short string // ISO 639-1 code, also accepted as input shorthand
}

// catalog is the single source of truth; order drives --help and error listing.
// Locked to {zh, en, ja} as of 2026-05-28: TUI bundles only ship for zh/en
// (ja falls back to the zh bundle), and Lark API client code only branches on
// these three for localization. Adding more entries here is meaningful only
// after the downstream codepaths (mail signature locale, TUI bundle) gain
// branches for them.
var catalog = []langEntry{
	{LangZhCN, "zh"}, {LangEnUS, "en"}, {LangJaJP, "ja"},
}

// find matches a short code or Feishu locale against the catalog (case-sensitive).
func find(s string) (langEntry, bool) {
	for _, e := range catalog {
		if string(e.Code) == s || e.Short == s {
			return e, true
		}
	}
	return langEntry{}, false
}

// Parse resolves a short code or Feishu locale to its canonical Lang.
// "" and unrecognized values return ("", false).
func Parse(s string) (Lang, bool) {
	e, ok := find(s)
	return e.Code, ok
}

// IsEnglish reports whether l uses the English TUI bundle (robust to "en_us"
// and legacy "en").
func (l Lang) IsEnglish() bool {
	e, _ := find(string(l))
	return e.Code == LangEnUS
}

// Base returns the ISO 639-1 short code ("en_us" → "en"), or "" if unknown.
func (l Lang) Base() string {
	e, _ := find(string(l))
	return e.Short
}

// Codes lists the canonical locales, for --help and error messages.
func Codes() []string {
	out := make([]string, len(catalog))
	for i, e := range catalog {
		out[i] = string(e.Code)
	}
	return out
}

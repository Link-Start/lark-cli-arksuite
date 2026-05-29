// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package common

import (
	"testing"

	"github.com/larksuite/cli/internal/core"
	"github.com/larksuite/cli/internal/i18n"
)

func TestRuntimeContext_Lang(t *testing.T) {
	tests := []struct {
		name   string
		stored i18n.Lang
		want   i18n.Lang
	}{
		{"canonical locale", i18n.LangJaJP, i18n.LangJaJP},
		{"legacy short value normalizes", "ja", i18n.LangJaJP},
		{"legacy short zh normalizes", "zh", i18n.LangZhCN},
		{"unset stays empty", "", ""},
		// Flipped semantics: unrecognized non-empty values are now treated
		// as legacy storage from the pre-2026-05-28 14-language catalog
		// and silently coerced to LangZhCN, not left empty.
		{"unrecognized garbage coerces to zh", "klingon", i18n.LangZhCN},
		{"legacy ko_kr coerces to zh", "ko_kr", i18n.LangZhCN},
		{"legacy fr_fr coerces to zh", "fr_fr", i18n.LangZhCN},
		{"legacy short ko coerces to zh", "ko", i18n.LangZhCN},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &RuntimeContext{Config: &core.CliConfig{Lang: tt.stored}}
			if got := ctx.Lang(); got != tt.want {
				t.Errorf("Lang() with stored %q = %q, want %q", tt.stored, got, tt.want)
			}
		})
	}
}

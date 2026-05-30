// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package mail

import (
	"encoding/json"

	"github.com/larksuite/cli/internal/event"
)

// matchMailbox compares the V2-envelope payload's event.mail_address against
// the normalized params.mailbox. Drops events whose mail_address doesn't match.
//
// Fail-open policy: if params.mailbox is empty (no filter), the payload can't
// be parsed, or the payload omits/moves event.mail_address (defensive —
// upstream schema may evolve), accept the event rather than silently dropping
// legitimate traffic. Only a present-but-mismatched mail_address drops.
//
// IMPORTANT: caller must ensure params.mailbox is already normalized to a real
// email (not "me"). normalizeMailParams handles this.
func matchMailbox(raw *event.RawEvent, params map[string]string) bool {
	target := params["mailbox"]
	if target == "" {
		return true
	}
	var env struct {
		Event struct {
			MailAddress string `json:"mail_address"`
		} `json:"event"`
	}
	if err := json.Unmarshal(raw.Payload, &env); err != nil {
		return true // fail-open on unparseable payload
	}
	if env.Event.MailAddress == "" {
		return true // fail-open on shape drift (field absent/moved); let Process decide
	}
	return env.Event.MailAddress == target
}

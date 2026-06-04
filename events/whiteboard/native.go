// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package whiteboard

type BoardWhiteboardUpdatedV1Data struct {
	WhiteboardID string       `json:"whiteboard_id"`
	OperatorIDs  []OperatorID `json:"operator_ids"`
}

type OperatorID struct {
	OpenID  string `json:"open_id"`
	UnionID string `json:"union_id"`
	UserID  string `json:"user_id"`
}

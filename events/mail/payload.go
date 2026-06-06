// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

// Package mail registers Mail-domain EventKeys and supporting types.
package mail

// MailReceivedPayload is the unified output schema for
// mail.user_mailbox.event.message_received_v1. Fields are populated
// conditionally based on the msg-format param; readers should treat
// absent fields as null/zero:
//
//   - msg-format=event           : MessageID, MailAddress, MailboxType, Subscriber
//   - msg-format=metadata        : event fields + From, Subject, Snippet, FolderID, LabelIDs
//   - msg-format=plain_text_full : metadata fields + BodyText
//   - msg-format=full            : metadata + BodyText + BodyHTML + Attachments
type MailReceivedPayload struct {
	// Always present (msg-format=event and above)
	MessageID   string     `json:"message_id"   desc:"Unique message identifier"`
	MailAddress string     `json:"mail_address" desc:"Recipient mailbox address (matches the subscribed mailbox)"`
	MailboxType int        `json:"mailbox_type" desc:"Mailbox type enum: 1=primary, 2=shared"`
	Subscriber  Subscriber `json:"subscriber"   desc:"Subscribers of the event — the users whose mailbox received this message"`

	// Populated when msg-format >= metadata
	From     string   `json:"from,omitempty"      desc:"Sender email address (msg-format >= metadata)"`
	Subject  string   `json:"subject,omitempty"   desc:"Mail subject (msg-format >= metadata)"`
	Snippet  string   `json:"snippet,omitempty"   desc:"Body preview, first ~100 chars (msg-format >= metadata)"`
	FolderID string   `json:"folder_id,omitempty" desc:"Folder ID containing this message (msg-format >= metadata)"`
	LabelIDs []string `json:"label_ids,omitempty" desc:"Label IDs attached (msg-format >= metadata)"`

	// Populated when msg-format >= plain_text_full
	BodyText string `json:"body_text,omitempty" desc:"Plain-text body (msg-format >= plain_text_full)"`

	// Populated when msg-format=full only
	BodyHTML    string           `json:"body_html,omitempty"   desc:"HTML body (msg-format=full only)"`
	Attachments []MailAttachment `json:"attachments,omitempty" desc:"Attachment metadata (msg-format=full only)"`
}

// MailAttachment mirrors the message API's attachment object. The fetch
// endpoint returns id / filename / is_inline (no size or MIME type); resolve a
// download URL separately via the attachment download_url API. Inline images
// embedded in the HTML body (signatures, logos) are returned here too, flagged
// by is_inline so consumers can skip them.
// Verified via: lark-cli api GET /open-apis/mail/v1/user_mailboxes/{mailbox}/messages/{id}
type MailAttachment struct {
	AttachmentID string `json:"attachment_id" desc:"Attachment ID (resolve a download URL via the attachment download_url API)"`
	Filename     string `json:"filename"      desc:"Original filename"`
	IsInline     bool   `json:"is_inline"     desc:"True for images embedded in the HTML body (e.g. signatures/logos); false for real file attachments"`
}

// Subscriber is the raw event-envelope `subscriber` block: the set of users
// whose mailbox received the message. Each element carries the three Feishu
// user identifier forms (user_id, open_id, union_id); fields are omitempty
// because in practice only the IDs the app is scoped to are populated.
type Subscriber struct {
	UserIDs []SubscriberUserID `json:"user_ids,omitempty" desc:"Recipients of the mail event (mailbox owners)"`
}

type SubscriberUserID struct {
	UserID  string `json:"user_id,omitempty"  desc:"Tenant-scoped user_id"`
	OpenID  string `json:"open_id,omitempty"  desc:"App-scoped open_id"`
	UnionID string `json:"union_id,omitempty" desc:"Cross-tenant union_id"`
}

// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package consume

// NOTE: TestNormalizeParams_ErrorIsWrappedWithEventKey drives the real Run()
// path — NormalizeParams fails before EnsureBus, so no bus/transport is
// actually exercised, yet the assertion covers the production error-wrapping
// code (not a reconstruction). TestDoHello_PassesSubscriptionIDToWire covers
// the Hello wire encoding. The cleanup-error WARN format is verified
// end-to-end by the sandbox E2E (TestEventConsume_Mail_ReadyAndTimeout),
// which asserts the real stderr contract rather than a duplicated literal.

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net"
	"strings"
	"testing"

	"github.com/larksuite/cli/internal/event"
	"github.com/larksuite/cli/internal/event/protocol"
	"github.com/larksuite/cli/internal/event/transport"
)

// fakeRT is a minimal event.APIClient mock.
type fakeRT struct {
	err error
}

func (f *fakeRT) CallAPI(_ context.Context, _, _ string, _ interface{}) (json.RawMessage, error) {
	return nil, f.err
}

func TestNormalizeParams_ErrorIsWrappedWithEventKey(t *testing.T) {
	// Drive the real Run() path. NormalizeParams failure returns before
	// EnsureBus, so the bus/transport is never actually contacted, but the
	// error-wrapping under test (`fmt.Errorf("normalize params for %s: %w")`)
	// is the genuine production code path — if Run() ever stops wrapping, this
	// test fails.
	const key = "test.evt_normalize_fail"
	event.RegisterKey(event.KeyDefinition{
		Key:       key,
		EventType: key,
		Schema:    event.SchemaDef{Custom: &event.SchemaSpec{Raw: json.RawMessage(`{"type":"object"}`)}},
		NormalizeParams: func(_ context.Context, _ event.APIClient, _ map[string]string) error {
			return errors.New("simulated normalize failure")
		},
	})
	defer event.UnregisterKeyForTest(key)

	err := Run(context.Background(), transport.New(), "app", "", "", Options{
		EventKey: key,
		Runtime:  &fakeRT{},
		Quiet:    true,
	})
	if err == nil {
		t.Fatal("expected Run to fail when NormalizeParams errors")
	}
	if !strings.Contains(err.Error(), "normalize params for "+key+":") {
		t.Errorf("error not wrapped with EventKey prefix: %v", err)
	}
	if !strings.Contains(err.Error(), "simulated normalize failure") {
		t.Errorf("underlying error not propagated: %v", err)
	}
}

func TestDoHello_PassesSubscriptionIDToWire(t *testing.T) {
	a, b := net.Pipe()
	defer a.Close()
	defer b.Close()

	// Server-side: read Hello, decode, assert SubscriptionID, send ack
	done := make(chan string, 1)
	go func() {
		br := bufio.NewReader(b)
		line, err := protocol.ReadFrame(br)
		if err != nil {
			done <- "READ_ERR:" + err.Error()
			return
		}
		msg, err := protocol.Decode(bytes.TrimRight(line, "\n"))
		if err != nil {
			done <- "DECODE_ERR:" + err.Error()
			return
		}
		if hello, ok := msg.(*protocol.Hello); ok {
			done <- hello.SubscriptionID
			// send ack so client can return
			ack := protocol.NewHelloAck("v1", true)
			_ = protocol.EncodeWithDeadline(b, ack, protocol.WriteTimeout)
		} else {
			done <- "WRONG_TYPE"
		}
	}()

	ack, _, err := doHello(a, "mail.x", []string{"mail.x"}, "mail.x:alice")
	if err != nil {
		t.Fatalf("doHello error: %v", err)
	}
	if ack == nil {
		t.Fatal("got nil ack")
	}
	got := <-done
	if got != "mail.x:alice" {
		t.Errorf("Hello.SubscriptionID on wire = %q, want %q", got, "mail.x:alice")
	}
}

// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"net/http"
	"strings"
	"testing"
)

func TestBillingSessionRedeemVoucherReadsTheRedemption(t *testing.T) {
	t.Parallel()

	// Kaiten answers a dry run 200 and a redemption 201. The spec declares only the 201,
	// so the generated client leaves the 200 undecoded and the wrapper has to read it.
	for name, status := range map[string]int{"redeemed": http.StatusCreated, "dry run": http.StatusOK} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client, sent := givenClient(t, respondJSON(status, `{"name":"Launch"}`))

			redemption, err := client.Billing.Session.RedeemVoucher(t.Context(), SessionRedeem{Code: "SUMMER"})
			if err != nil {
				t.Fatalf("RedeemVoucher() error = %v", err)
			}
			if redemption.Name != "Launch" {
				t.Errorf("redemption = %+v, want the decoded body", redemption)
			}
			if !strings.HasSuffix(sent.only(t).path, "/public/session/vouchers/redeem") {
				t.Errorf("path = %s", sent.only(t).path)
			}
		})
	}
}

func TestBillingPublicConfirmCheckoutSessionTargetsTheSession(t *testing.T) {
	t.Parallel()

	client, sent := givenClient(t, respondJSON(http.StatusOK, `{"status":"completed"}`), WithBearerToken("ksc_token"))

	confirmation, err := client.Billing.Public.ConfirmCheckoutSession(t.Context(), "ses_1")
	if err != nil {
		t.Fatalf("ConfirmCheckoutSession() error = %v", err)
	}

	request := sent.only(t)
	if request.method != http.MethodPost || !strings.HasSuffix(request.path, "/public/checkout-sessions/ses_1/confirm") {
		t.Errorf("request = %s %s", request.method, request.path)
	}
	if got := request.header.Get("Authorization"); got != "Bearer ksc_token" {
		t.Errorf("Authorization = %q", got)
	}
	if confirmation.Status != CheckoutConfirmationStatus("completed") {
		t.Errorf("status = %q", confirmation.Status)
	}
}

func TestBillingCheckoutSessionsListWalksEveryPage(t *testing.T) {
	t.Parallel()

	client, sent := givenClient(t, func(req *http.Request) (*http.Response, error) {
		if req.URL.Query().Get("cursor") == "" {
			return jsonResponse(req, http.StatusOK, `{"items":[{"id":"a"}],"hasMore":true,"nextCursor":"next"}`), nil
		}
		return jsonResponse(req, http.StatusOK, `{"items":[{"id":"b"}],"hasMore":false}`), nil
	})

	items, err := client.Billing.CheckoutSessions.List(t.Context(), nil)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(items) != 2 || sent.count() != 2 {
		t.Errorf("items = %d, requests = %d, want 2 and 2", len(items), sent.count())
	}
}

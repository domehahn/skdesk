package ticket_test

import (
	"testing"

	"github.com/domehahn/skdesk/pkg/signer"
	"github.com/domehahn/skdesk/pkg/ticket"
)

func TestTicketLifecycle(t *testing.T) {
	mgr := ticket.NewManager()
	tick := mgr.Create("prod-deployer", "sha256:abcd", "production", "Network egress to unknown domain")
	if tick.Status != ticket.StatusPending {
		t.Fatalf("expected PENDING status")
	}

	kp, err := signer.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	approved, err := mgr.Approve(tick.ID, "security-lead", kp.PrivateKeyHex)
	if err != nil {
		t.Fatalf("approval failed: %v", err)
	}
	if approved.Status != ticket.StatusApproved {
		t.Errorf("expected APPROVED status")
	}
	if approved.Signature == "" {
		t.Errorf("expected non-empty signature")
	}
}

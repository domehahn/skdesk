package signer_test

import (
	"testing"

	"github.com/domehahn/skdesk/pkg/signer"
)

func TestSigner(t *testing.T) {
	kp, err := signer.GenerateKeyPair()
	if err != nil {
		t.Fatalf("failed to generate keypair: %v", err)
	}

	payload := []byte("approval_ticket_12345")
	sig, err := signer.SignPayload(payload, kp.PrivateKeyHex)
	if err != nil {
		t.Fatalf("failed to sign payload: %v", err)
	}

	valid := signer.VerifySignature(payload, sig, kp.PublicKeyHex)
	if !valid {
		t.Errorf("expected signature to be valid")
	}

	invalid := signer.VerifySignature([]byte("tampered"), sig, kp.PublicKeyHex)
	if invalid {
		t.Errorf("expected signature verification to fail on tampered payload")
	}
}

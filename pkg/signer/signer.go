package signer

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type KeyPair struct {
	PublicKeyHex  string `json:"public_key"`
	PrivateKeyHex string `json:"private_key"`
}

func GenerateKeyPair() (*KeyPair, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &KeyPair{
		PublicKeyHex:  hex.EncodeToString(pub),
		PrivateKeyHex: hex.EncodeToString(priv),
	}, nil
}

func SignPayload(payload []byte, privateKeyHex string) (string, error) {
	privBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key hex: %w", err)
	}
	if len(privBytes) != ed25519.PrivateKeySize {
		return "", fmt.Errorf("invalid private key size: %d", len(privBytes))
	}

	h := sha256.Sum256(payload)
	sig := ed25519.Sign(privBytes, h[:])
	return hex.EncodeToString(sig), nil
}

func VerifySignature(payload []byte, signatureHex, publicKeyHex string) bool {
	pubBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil || len(pubBytes) != ed25519.PublicKeySize {
		return false
	}
	sigBytes, err := hex.DecodeString(signatureHex)
	if err != nil || len(sigBytes) != ed25519.SignatureSize {
		return false
	}

	h := sha256.Sum256(payload)
	return ed25519.Verify(pubBytes, h[:], sigBytes)
}

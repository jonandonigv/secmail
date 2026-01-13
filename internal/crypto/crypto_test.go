package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecryptPassphrase(t *testing.T) {
	// Generate RSA key pair
	publicKey, privateKey, err := GenerateRSAKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Test passphrase
	passphrase := "test-passphrase-123"

	// Encrypt with public key
	encrypted, err := EncryptPassphrase(passphrase, publicKey)
	if err != nil {
		t.Fatalf("Failed to encrypt passphrase: %v", err)
	}

	// Decrypt with private key
	decrypted, err := DecryptPassphrase(encrypted, privateKey)
	if err != nil {
		t.Fatalf("Failed to decrypt passphrase: %v", err)
	}

	// Verify it matches
	if decrypted != passphrase {
		t.Errorf("Decrypted passphrase does not match original: got %s, want %s", decrypted, passphrase)
	}
}

func TestEncryptDecryptBody(t *testing.T) {
	// Test plaintext
	plaintext := []byte("This is a test message for encryption.")

	// Encrypt
	ciphertext, passphrase, err := EncryptBody(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt body: %v", err)
	}

	// Ensure ciphertext is not empty and different from plaintext
	if len(ciphertext) == 0 {
		t.Error("Ciphertext is empty")
	}
	if string(ciphertext) == string(plaintext) {
		t.Error("Ciphertext is the same as plaintext")
	}

	// Decrypt
	decrypted, err := DecryptBody(ciphertext, passphrase)
	if err != nil {
		t.Fatalf("Failed to decrypt body: %v", err)
	}

	// Verify it matches
	if string(decrypted) != string(plaintext) {
		t.Errorf("Decrypted body does not match original: got %s, want %s", string(decrypted), string(plaintext))
	}
}

func TestGenerateRSAKeyPair(t *testing.T) {
	publicKey, privateKey, err := GenerateRSAKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Check that PEMs are not empty
	if len(publicKey) == 0 {
		t.Error("Public key PEM is empty")
	}
	if len(privateKey) == 0 {
		t.Error("Private key PEM is empty")
	}

	// Check that public key contains expected header
	if !bytes.Contains(publicKey, []byte("-----BEGIN PUBLIC KEY-----")) {
		t.Error("Public key PEM does not contain expected header")
	}

	// Check that private key contains expected header
	if !bytes.Contains(privateKey, []byte("-----BEGIN PRIVATE KEY-----")) {
		t.Error("Private key PEM does not contain expected header")
	}
}

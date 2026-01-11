package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"

	"filippo.io/age"
)

// EncryptBody encrypts the plaintext using age with a randomly generated passphrase.
// Returns the ciphertext, the passphrase (for encrypting with recipient keys), and any error.
func EncryptBody(plaintext []byte) (ciphertext []byte, passphrase string, err error) {
	// Generate a random 32-byte passphrase
	passphraseBytes := make([]byte, 32)
	if _, err = rand.Read(passphraseBytes); err != nil {
		return nil, "", err
	}
	passphrase = base64.StdEncoding.EncodeToString(passphraseBytes)

	// Create age recipient
	recipient, err := age.NewScryptRecipient(passphrase)
	if err != nil {
		return nil, "", err
	}

	// Encrypt
	var buf bytes.Buffer
	w, err := age.Encrypt(&buf, recipient)
	if err != nil {
		return nil, "", err
	}
	if _, err = w.Write(plaintext); err != nil {
		return nil, "", err
	}
	if err = w.Close(); err != nil {
		return nil, "", err
	}

	ciphertext = buf.Bytes()
	return ciphertext, passphrase, nil
}

// DecryptBody decrypts the ciphertext using age with the provided passphrase.
func DecryptBody(ciphertext []byte, passphrase string) (plaintext []byte, err error) {
	// Create age identity
	identity, err := age.NewScryptIdentity(passphrase)
	if err != nil {
		return nil, err
	}

	// Decrypt
	r, err := age.Decrypt(bytes.NewReader(ciphertext), identity)
	if err != nil {
		return nil, err
	}

	plaintext, err = io.ReadAll(r)
	return plaintext, err
}

// EncryptPassphrase encrypts the passphrase using RSA OAEP with the recipient's public key.
func EncryptPassphrase(passphrase string, publicKeyPEM []byte) ([]byte, error) {
	// Parse public key
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("not an RSA public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	// Encrypt
	encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPub, []byte(passphrase), nil)
	return encrypted, err
}

// DecryptPassphrase decrypts the passphrase using RSA OAEP with the recipient's private key.
func DecryptPassphrase(encrypted []byte, privateKeyPEM []byte) (string, error) {
	// Parse private key
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return "", errors.New("not an RSA private key")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// Decrypt
	decrypted, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, encrypted, nil)
	return string(decrypted), err
}

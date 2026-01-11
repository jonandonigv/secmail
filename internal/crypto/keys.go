package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

// GenerateRSAKeyPair generates a new RSA key pair (2048 bits) and returns PEM-encoded public and private keys.
func GenerateRSAKeyPair() (publicKeyPEM, privateKeyPEM []byte, err error) {
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	// Encode private key to PEM
	privateKeyDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyDER,
	}
	privateKeyPEM = pem.EncodeToMemory(privateKeyBlock)

	// Encode public key to PEM
	publicKeyDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyDER,
	}
	publicKeyPEM = pem.EncodeToMemory(publicKeyBlock)

	return publicKeyPEM, privateKeyPEM, nil
}

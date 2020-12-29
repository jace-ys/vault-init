package encryption

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type LocalEncryption struct {
	SecretKey string
}

func NewLocalEncryption(secretKey string) (*LocalEncryption, error) {
	if secretKey == "" {
		return nil, fmt.Errorf("no secret key provided")
	}

	if len(secretKey) != 32 {
		return nil, fmt.Errorf("provided secret key is %d-bytes, expected 32", len(secretKey))
	}

	return &LocalEncryption{
		SecretKey: secretKey,
	}, nil
}

func (e *LocalEncryption) Encrypt(ctx context.Context, plaintext string) (string, error) {
	block, err := aes.NewCipher([]byte(e.SecretKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (e *LocalEncryption) Decrypt(ctx context.Context, data string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", nil
	}

	block, err := aes.NewCipher([]byte(e.SecretKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

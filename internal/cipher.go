package whitearchive

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

type Cipher struct {
	gcm cipher.AEAD
}

func NewCipher(keyRaw []byte) (*Cipher, error) {

	hash := sha256.Sum256(keyRaw)
	key := hash[:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &Cipher{
		gcm: gcm,
	}, nil
}

func (c *Cipher) Encrypt(data []byte) ([]byte, error) {
	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return c.gcm.Seal(nonce, nonce, data, nil), nil
}

func (c *Cipher) Decrypt(data []byte) ([]byte, error) {
	nonce, ciphertext := data[:c.gcm.NonceSize()], data[c.gcm.NonceSize():]
	return c.gcm.Open(nil, nonce, ciphertext, nil)
}

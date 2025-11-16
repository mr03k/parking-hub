package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

type AESEncryptor struct {
	secret []byte
	iv     []byte
}

func NewAESEncryptor(secret, iv []byte) *AESEncryptor {
	return &AESEncryptor{
		secret: secret,
		iv:     iv,
	}
}

func (a *AESEncryptor) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.secret)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}
	cfb := cipher.NewCFBEncrypter(block, a.iv)
	cipherText := make([]byte, len(data))
	cfb.XORKeyStream(cipherText, data)
	encodeData, err := EncodeBase64(cipherText)
	if err != nil {
		return nil, err
	}
	return encodeData, nil
}

func (a *AESEncryptor) Decrypt(data []byte) ([]byte, error) {
	data, err := DecodeBase64(data)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(a.secret)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}
	cfb := cipher.NewCFBDecrypter(block, a.iv)
	plainText := make([]byte, len(data))
	cfb.XORKeyStream(plainText, data)
	return plainText, nil
}

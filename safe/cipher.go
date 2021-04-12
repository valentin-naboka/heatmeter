package safe

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"heatmeter/generated"
)

func Decrypt(data string) (string, error) {
	key, err := hex.DecodeString(generated.Data1)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv, err := hex.DecodeString(generated.Data2)
	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), err
}

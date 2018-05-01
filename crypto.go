package gethinx

import (
	"crypto/rand"
	"encoding/hex"
)

// Key return random string
func Key(key int) string {
	buf := make([]byte, key)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err) // out of randomness, should never happen
	}

	return hex.EncodeToString(buf)
	// or hex.EncodeToString(buf)
	// or base64.StdEncoding.EncodeToString(buf)
}

/*
func main() {
	text := []byte("My name is Astaxie")
	key := []byte("the-key-has-to-be-32-bytes-long!")

	ciphertext, err := encrypt(text, key)
	if err != nil {
		// TODO: Properly handle error
		log.Fatal(err)
	}
	fmt.Printf("%s => %x\n", text, ciphertext)

	plaintext, err := decrypt(ciphertext, key)
	if err != nil {
		// TODO: Properly handle error
		log.Fatal(err)
	}
	fmt.Printf("%x => %s\n", ciphertext, plaintext)
}


func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}


*/

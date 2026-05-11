package bcryptx

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh"
	"io"
	"strings"
)

// BcryptHash Encrypt passwords using bcrypt
func BcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck Compare plaintext passwords with database hash values
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateAESKeyBase64 generates a new cryptographically secure 32-byte (256-bit) AES key
// and returns its Base64 encoded string representation.
// This key can then be securely stored or transmitted.
func GenerateAESKeyBase64() (string, error) {
	key := make([]byte, 32) // 32 bytes for AES-256
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", fmt.Errorf("failed to generate random AES key: %w", err)
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// EncryptAESGCM performs AES-256-GCM encryption.
// plaintext: The data to be encrypted.
// key: 32-byte (256-bit) encryption key.
// additionalData: Optional authenticated data (AAD) that is not encrypted but authenticated.
// Returns: Nonce || Ciphertext || Tag, or error.
func EncryptAESGCM(plaintext []byte, key []byte, additionalData []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key length: expected 32 bytes for AES-256, got %d", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	nonce := make([]byte, nonceSize)
	// IMPORTANT: Use crypto/rand for cryptographically secure random Nonce
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate Nonce: %w", err)
	}

	// gcm.Seal prepends the Nonce to the ciphertext.
	// The order is: Nonce || Ciphertext || Tag
	ciphertext := gcm.Seal(nonce, nonce, plaintext, additionalData)
	return ciphertext, nil
}

// DecryptAESGCM performs AES-256-GCM decryption.
// encryptedData: Nonce || Ciphertext || Tag
// key: 32-byte (256-bit) encryption key.
// additionalData: Optional authenticated data (AAD) that was used during encryption.
// Returns: Plaintext, or error if decryption or authentication fails.
func DecryptAESGCM(encryptedData []byte, key []byte, additionalData []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key length: expected 32 bytes for AES-256, got %d", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("encrypted data is too short to contain Nonce")
	}

	nonce, ciphertextWithTag := encryptedData[:nonceSize], encryptedData[nonceSize:]

	// gcm.Open validates the tag and decrypts the ciphertext.
	// If authentication fails, an error is returned.
	plaintext, err := gcm.Open(nil, nonce, ciphertextWithTag, additionalData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt or authenticate data: %w", err)
	}
	return plaintext, nil
}

// GenerateED25519KeyBase64 generates a new Ed25519 public/private key pair.
// It returns the Base64 encoded string representations of both the public and private keys.
// The private key includes the public key for convenience in Ed25519.
func GenerateED25519KeyBase64() (publicKeyBase64 string, privateKeyBase64 string, err error) {

	publicKey, privateKey, err := GenerateED25519KeyOpenSSH()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate Ed25519 Open ssh key pair: %w", err)
	}
	// Base64 encode the public key
	publicKeyBase64 = base64.StdEncoding.EncodeToString(publicKey)

	// Base64 encode the private key
	// Note: Ed25519 private keys inherently contain enough information to derive the public key.
	// The `privateKey` returned by GenerateKey is typically 64 bytes:
	// first 32 bytes are the seed, last 32 bytes are the public key.
	privateKeyBase64 = base64.StdEncoding.EncodeToString(privateKey)

	return publicKeyBase64, privateKeyBase64, nil
}

// GenerateED25519KeyOpenSSH generates an Ed25519 key pair
// and returns them in formats compatible with OpenSSH.
//
// publicKeyString will be the OpenSSH authorized_keys format string,
// e.g., "ssh-ed25519 AAAA... [comment]".
//
// privateKeyString will be the OpenSSH PKCS#8 PEM format string,
// e.g., "-----BEGIN OPENSSH PRIVATE KEY-----...-----END OPENSSH PRIVATE KEY-----".
//
// Note: The meaning of 'Base64' in the function name has been adjusted here. For OpenSSH format,
// these strings contain Base64 encoded parts, but they are complete, OpenSSH-compliant
// formatted strings, not just raw Base64 encodings.
func GenerateED25519KeyOpenSSH() (publicKeyByte, privateKeyByte []byte, err error) {
	// ed25519.GenerateKey uses crypto/rand.Reader for cryptographic randomness by default.
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate Ed25519 key pair: %w", err)
	}

	// --- Public Key: Convert to OpenSSH authorized_keys format ---
	// Convert the raw Ed25519 public key to an ssh.PublicKey object.
	sshPublicKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create SSH public key from Ed25519 public key: %w", err)
	}

	// Encode the ssh.PublicKey into the standard OpenSSH authorized_keys line format.
	// This format includes the key type prefix ("ssh-ed25519"), Base64 encoded key data, and an optional comment.
	publicKeyByte = ssh.MarshalAuthorizedKey(sshPublicKey)

	// --- Private Key: Convert to OpenSSH PKCS#8 PEM format ---
	// Key modification here: Directly cast ed25519.PrivateKey to crypto.Signer interface
	// then use ssh.MarshalPrivateKey.
	// This avoids certain internal types that ssh.NewSignerFromKey might return, which could cause MarshalPrivateKey to error.
	signer := crypto.Signer(privateKey) // Cast ed25519.PrivateKey to crypto.Signer

	// Encode the private key into a PEM block.
	// Passing nil means the private key will not be encrypted. If you want to encrypt the private key, you can provide a passphrase of type "".
	pemBlock, err := ssh.MarshalPrivateKey(signer, "")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal SSH private key to PEM format: %w", err)
	}

	// Encode the PEM block into bytes and convert to a
	privateKeyByte = pem.EncodeToMemory(pemBlock)

	return publicKeyByte, privateKeyByte, nil
}

// Ed25519PrivateKeyFromBase64 decodes a Base64 encoded Ed25519 private key string to ed25519.PrivateKey.
func Ed25519PrivateKeyFromBase64(base64Key string) (ed25519.PrivateKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key base64 string: %w", err)
	}

	privateKey, err := parsePrivateKey(string(keyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key : %w", err)
	}
	// Ed25519 private keys are 64 bytes long
	if len(privateKey) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid private key length: expected %d bytes, got %d", ed25519.PrivateKeySize, len(keyBytes))
	}
	return privateKey, nil
	//return ed25519.PrivateKey(keyBytes), nil
}

// parsePrivateKey parses an Ed25519 private key in SSH format.
func parsePrivateKey(privKey string) (ed25519.PrivateKey, error) {

	// Decode the PEM format private key
	block, _ := pem.Decode([]byte(privKey))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Parse the SSH private key
	priv, err := ssh.ParseRawPrivateKey([]byte(privKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH private key: %w", err)
	}

	// Convert to ed25519.PrivateKey
	edPriv, ok := priv.(*ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an Ed25519 private key")
	}

	return *edPriv, nil
}

// Ed25519PublicKeyFromBase64 decodes a Base64 encoded Ed25519 public key string to ed25519.PublicKey.
func Ed25519PublicKeyFromBase64(base64Key string) (ed25519.PublicKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key base64 string: %w", err)
	}

	publicKey, err := parsePublicKey(string(keyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key : %w", err)
	}

	// Ed25519 public keys are 32 bytes long
	if len(publicKey) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid public key length: expected %d bytes, got %d", ed25519.PublicKeySize, len(keyBytes))
	}
	return publicKey, nil
	//return ed25519.PublicKey(keyBytes), nil
}

// parsePublicKey parses an Ed25519 public key in SSH format.
func parsePublicKey(pubKey string) (ed25519.PublicKey, error) {

	// SSH public key format: "ssh-ed25519 <base64> [comment]"
	parts := strings.Fields(pubKey)
	if len(parts) < 2 || parts[0] != "ssh-ed25519" {
		return nil, fmt.Errorf("invalid SSH public key format")
	}

	// Parse the SSH public key
	pub, _, _, _, err := ssh.ParseAuthorizedKey([]byte(pubKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH public key: %w", err)
	}

	// Convert to ed25519.PublicKey
	edPub, ok := pub.(ssh.CryptoPublicKey).CryptoPublicKey().(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an Ed25519 public key")
	}

	return edPub, nil
}

// Ed25519DataSignGenerate signs data using the Ed25519 algorithm.
// The signing process involves hashing and cryptographic operations on the message
// using the private key to produce a unique digital signature.
//
// Parameters:
//
//	message []byte: The original message data to be signed. This can be a byte slice
//	  of arbitrary length.
//	privateKey ed25519.PrivateKey: The Ed25519 private key used to generate the signature.
//	  This key is typically obtained via `ed25519.GenerateKey` or loaded from secure storage.
//
// Returns:
//
//	string: The Base64 encoded string representation of the generated signature.
//	  The binary signature data is encoded for convenient network transmission or storage.
//
// Security Note:
//
//	The private key must always be kept confidential. Disclosure of the private key
//	will allow anyone who possesses it to forge your signatures.
//	This function does not explicitly handle errors (e.g., invalid private key)
//	because `ed25519.Sign` generally does not fail if the private key type is correct.
//	However, in more robust implementations, pre-checks for private key validity might be added.
func Ed25519DataSignGenerate(message []byte, privateKey ed25519.PrivateKey) string {
	// Use the ed25519 library's Sign function to generate the signature.
	// This function takes the private key and message as input and returns
	// a fixed-length digital signature.
	signature := ed25519.Sign(privateKey, message)

	// Encode the binary signature data into a Base64 string for easier
	// transmission and storage.
	return base64.StdEncoding.EncodeToString(signature)
}

// Ed25519DataSignVerify verifies the validity of an Ed25519 digital signature.
// It uses the public key, the original message, and the received signature
// to check the integrity and authenticity of the signature.
//
// Parameters:
//
//	baseSign string: The Base64 encoded string representation of the signature to be verified.
//	message []byte: The original message data that was signed. This *must* be
//	  identical to the message used during the signing process.
//	publicKey ed25519.PublicKey: The Ed25519 public key used to verify the signature.
//	  This key is typically derived from the corresponding private key or obtained
//	  through other secure channels.
//
// Returns:
//
//	bool: Returns true if the signature verification is successful (i.e., the signature
//	  is valid and the message has not been tampered with); otherwise, returns false.
//	error: Returns an error if the Base64 decoding fails (e.g., `baseSign` is not
//	  a valid Base64 string). The `ed25519.Verify` function itself returns false on
//	  verification failure, not an error.
//
// Usage Note:
//
//	Successful signature verification means that the message has not been altered
//	since it was signed, and that the signature was indeed generated by the holder
//	of the corresponding private key. In network communication, a verification
//	failure typically indicates data tampering or a forged signature.
func Ed25519DataSignVerify(baseSign string, message []byte, publicKey ed25519.PublicKey) (bool, error) {
	// Decode the Base64 encoded signature string back into its original binary byte slice.
	// If decoding fails, it indicates that the input `baseSign` is not in a correct format.
	signature, err := base64.StdEncoding.DecodeString(baseSign)
	if err != nil {
		// Return a formatted error message, wrapping the original decoding error.
		return false, fmt.Errorf("failed to decode signature base64 string: %w", err)
	}

	// Use the ed25519 library's Verify function to validate the signature.
	// This function checks:
	// 1. If the public key, message, and signature are consistent.
	// 2. If the signature was indeed generated by the corresponding private key.
	// 3. If the message has been tampered with after signing.
	// It returns true if all checks pass; otherwise, it returns false.
	return ed25519.Verify(publicKey, message, signature), nil
}

func PrivateKeyFromBase64(base64Key string) (ed25519.PrivateKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key base64 string: %w", err)
	}
	// Ed25519 private keys are 64 bytes long
	if len(keyBytes) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid private key length: expected %d bytes, got %d", ed25519.PrivateKeySize, len(keyBytes))
	}
	return ed25519.PrivateKey(keyBytes), nil
}

// PublicKeyFromBase64 decodes a Base64 encoded Ed25519 public key string to ed25519.PublicKey.
func PublicKeyFromBase64(base64Key string) (ed25519.PublicKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key base64 string: %w", err)
	}
	// Ed25519 public keys are 32 bytes long
	if len(keyBytes) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid public key length: expected %d bytes, got %d", ed25519.PublicKeySize, len(keyBytes))
	}
	return ed25519.PublicKey(keyBytes), nil
}

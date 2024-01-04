walkdir_get_files:
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
)

func main() {
	var SystemDir string

	flag.StringVar(&SystemDir, "dir", "", " give the dir for the root ")
	flag.Parse()
	if SystemDir == "" {

		panic("SystemDir is empty")
		return
	}

	files, _ := ioutil.ReadDir(SystemDir)
	for _, file := range files {
		if file.IsDir() {

			dir := filepath.Join(SystemDir, file.Name())
			filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
				if d.IsDir() {
					return nil
				}
				fmt.Println(path)
				return nil
			})
		}
	}

}
-------------------------------------------------------------------------------------------
sign data and verity sign data

package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	publicKeyStr := "MIGfMA0GCSqGSIb3DQEBAQU8YqokIXr4AB7bw1BGnb7ErStR+SSXdsuOtC5prOfEdKrdU8mfbcf5n64o3b9qGwFzgSsMHnhpO6najdCXCxqLGX6ZF2hov2tiOqFqJ4ckjRvOmQIDAQAB"
	publicKeyPem := `-----BEGIN RSA PUBLIC KEY-----` + "\n" + publicKeyStr + "\n" + `-----END RSA PUBLIC KEY-----`

	privateKeyStr := "MIICeAIBADANBgkKpBxkgbLxiqiQhevgAHtvDUEadvsStK1H5JJd2y460Lmms58R0qt1TyZ9tx/mfrijdv2obAXOBKwweeGk7qdqN0JcLGosZfpkXaGi/a2I6oWonhySNG86ZAgMBAAECgYEAkdQTj3r8vq4R8+/9RdDJ4uL1yAjcIWsCH2w7WHkHmkrIb/qFc47TqT3yD9wYiHVcMBrZyG2zuc53eJeAR83d8wRscocj2GIzsNjZzUEhYkoItrLMOH/I8dKb2Z85x/HrkbdYTf1qCXpxhvAUsdGKGfIbSyjhymgeGCWoUQ8KZ3ECQQD3c4BLYhFPTb00gNv7yMvGxIXYmecOsApoqJG+9vUVuS5bHR79ToV+E3/7Uyr17lItt3yJVhaxKJAEOB5yj5RdAkEA5v8Ls2EekqAAh0pfZYf1PaXEIv5KF8mLF0fsRoNrn6GGL3qfNCvQk3ASZ6Vwyc8RRN5KyTgmY+cyqyTUkuX/bQJBAKUJteGRMLZRxQWFhDL0A2U4oYSLcR3Mr8SJ2VsiXuf0MES4sXiErGggHVXEbHzGTK0NGdSHRG83/IWz4CrMNEkCQQCACmmK8c+HiOciFuiQF++pT0RL/VZGnzHZIsXmRByY7Gi70qWCvrKrtxiMmRjO1FeHLAyaQuSMxe/BC/ZEwvZ1AkAg4FgPWMu6PrpYMfpFH6U6sOpSI5wtw3ac6Rmu1rPAwewy8wBI0iIoJJ+yeVfqkK02AOWEY0XaasFmPdFItb/j"
	privateKeyPem := `-----BEGIN RSA PRIVATE KEY-----` + "\n" + privateKeyStr + "\n" + `-----END RSA PRIVATE KEY-----`

	requestData := "your_request_data_here"

	// Step 1: Load the RSA Private Key
	block, _ := pem.Decode([]byte(privateKeyPem))
	if block == nil {
		fmt.Println("29 failed to parse PEM block containing the key")
		return
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("failed to parse private key:", err)
		return
	}

	// Hash the requestData using SHA-256
	hashed := sha256.Sum256([]byte(requestData))

	// Sign the hashed data with the private key
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Println("error signing:", err)
		return
	}

	// The signature is now in the 'signature' variable.
	fmt.Println("Signature:", signature)

	// Step 3: Base64 Encode the Encrypted Data
	encodedEncryptedData := base64.StdEncoding.EncodeToString(signature)
	fmt.Println("Encoded encrypted data:", encodedEncryptedData)

	// Step 4: Sign with the Private Key
	// Signature (assuming it's base64 encoded)
	encodedSignature := encodedEncryptedData

	// Original message that was signed
	originalMessage := requestData

	// Decode the base64 encoded signature
	decodedSignature, err := base64.StdEncoding.DecodeString(encodedSignature)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode signature: %v\n", err)
		return
	}

	// Decode the PEM-formatted public key
	block, _ = pem.Decode([]byte(publicKeyPem))
	if block == nil {
		fmt.Fprintln(os.Stderr, "76 failed to parse PEM block containing the key")
		return
	}

	// Parse the public key
	pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse public key: %v\n", err)
		return
	}
	publicKey, ok := pubKeyInterface.(*rsa.PublicKey)
	if !ok {
		fmt.Fprintln(os.Stderr, "not an RSA public key")
		return
	}

	// Hash the original message
	hashed = sha256.Sum256([]byte(originalMessage))

	// Verify the signature
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], decodedSignature)
	if err != nil {
		fmt.Fprintf(os.Stderr, "signature verification failed: %v\n", err)
		return
	}

	fmt.Println("Signature verified successfully!")

}


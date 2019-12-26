package cipher

import (
    "crypto/rand"
    "crypto/rsa"
	"crypto/sha512"
    "crypto/x509"
    "encoding/pem"
	"io/ioutil"
	"fmt"
)

// GenerateKeyPair generates a new key pair
func GenerateKeyPair(bits int) (privkey *rsa.PrivateKey,pubkey *rsa.PublicKey, err error) {
	privkey, err = rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return
	}
	pubkey = &privkey.PublicKey
	return 
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) (privBytes []byte, err error) {
	privBytes = pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return 
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) (pubBytes []byte, err error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return
	}

	pubBytes = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})
	return 
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) (key *rsa.PrivateKey, err error) {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return
			
		}
	}
	key, err = x509.ParsePKCS1PrivateKey(b)
	return 
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) (key *rsa.PublicKey, err error) {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		err =  fmt.Errorf("Invalid pubkey type assertion")
	}
	return 
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) (ciphertext []byte, err error) {
	hash := sha512.New()
	ciphertext, err = rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	return 
}

// EncryptWithPublicKeyPKCS1v15 encrypts data with public key
func EncryptWithPublicKeyPKCS1v15(msg []byte, pub *rsa.PublicKey) (ciphertext []byte, err error) {
	ciphertext, err = rsa.EncryptPKCS1v15(rand.Reader,pub,msg)
	return 
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) (plaintext []byte, err error) {
	hash := sha512.New()
	plaintext, err = rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	return
} 

//LoadX509Certificate load cert
func LoadX509Certificate(certFile string) (cert *x509.Certificate, err error) {
    cf, err := ioutil.ReadFile(certFile)
    if err != nil {
		return
    }
    cpb, _ := pem.Decode(cf)
	cert, err = x509.ParseCertificate(cpb.Bytes)
    return 
}

//PubKeyFromCertFile get pubkey from certfile
func PubKeyFromCertFile(certFile string) (pubKey *rsa.PublicKey, err error){
	cert,err := LoadX509Certificate(certFile)
	pubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok{
		err = fmt.Errorf("Cannot cert convert to public key")
	}
	return
}

//EncryptPKCS1v15FromCert encrypt with rsa PKCS1v15 From Cert
func EncryptPKCS1v15FromCert(certFile string, msg []byte)(ciphertext []byte, err error) {
	pubKey, err := PubKeyFromCertFile(certFile)
	if err != nil{
		return
	}
	ciphertext,err = EncryptWithPublicKeyPKCS1v15(msg,pubKey)
	return
}
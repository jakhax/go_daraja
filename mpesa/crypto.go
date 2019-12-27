package mpesa

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
	"fmt"
)

//EncryptPassword encrypts password using daraja cert
// uses rsa PKCS1v15 as described in the daraja documentation
func EncryptPassword(password, environment string)(cipherText []byte,err error){
	var certB []byte
	switch environment{
		case SandBox:
			certB = SandBoxCert
			break
		case Production:
			certB = ProductionCert
			break
		default:
			err = fmt.Errorf("Invalid environment")
			return
	}
	cpb,_ := pem.Decode(certB)
	cert,err := x509.ParseCertificate(cpb.Bytes)
	pubKey ,ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok{
		err = fmt.Errorf("Cannot retrieve public key from cert")
		return
	}
	cipherText, err = rsa.EncryptPKCS1v15(rand.Reader,pubKey,[]byte(password))
	return
}
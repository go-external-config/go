package env

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"github.com/go-external-config/go/lang"
	"github.com/go-external-config/go/util/optional"
)

// Custom property source as an additional logic for decrypting properties using RSA private key, like pass=RSA:m+WQ5zMBqwMmEEP...
//
// Initialize at the beginning of the main package:
//
//	var environment = env.Instance().AddPropertySource(env.NewRsaPropertySource())
//
// Provide rsa.privateKey.path property to a PEM file to decrypt properties on a target environment.
//
// Encrypt property using RSA public key
//
//	echo -n "dbSecret123" | openssl pkeyutl -encrypt \
//	-pubin -inkey public2048.pem \
//	-pkeyopt rsa_padding_mode:oaep \
//	-pkeyopt rsa_oaep_md:sha256 | base64
//
// echo -n "dbSecret123" - supplies plaintext to stdin (no trailing newline).
//
// openssl pkeyutl -encrypt - tells OpenSSL to encrypt with a public key.
//
// -pubin -inkey public2048.pem - uses your RSA public key.
//
// -pkeyopt rsa_padding_mode:oaep - selects OAEP padding (modern, secure).
//
// -pkeyopt rsa_oaep_md:sha256 - sets both OAEP digest and (implicitly) MGF1 digest to SHA-256.
// (OpenSSL ≥1.0.2 automatically uses the same digest for MGF1 if not overridden).
//
// | base64 - encodes the ciphertext to text format. It is safe to remove any line breaks.
type RsaPropertySource struct {
	environment *Environment
}

func NewRsaPropertySource() *RsaPropertySource {
	return &RsaPropertySource{
		environment: Instance()}
}

func (s *RsaPropertySource) Name() string {
	return "RsaPropertySource"
}

func (s *RsaPropertySource) HasProperty(key string) bool {
	for _, source := range environment.PropertySources() {
		if source.Properties() != nil && source.HasProperty(key) {
			return strings.HasPrefix(source.Property(key), "RSA:")
		}
	}
	return false
}

func (s *RsaPropertySource) Property(key string) string {
	for _, source := range environment.PropertySources() {
		if source.Properties() != nil && source.HasProperty(key) {
			value := source.Property(key)[4:]
			rsaPrivateKeyPath := environment.Property("rsa.privateKey.path")
			return decryptWithPrivateKey(key, value, rsaPrivateKeyPath)
		}
	}
	panic("No value present for " + key)
}

func decryptWithPrivateKey(key, value, privateKeyPath string) string {
	data, _ := os.ReadFile(privateKeyPath)
	block, _ := pem.Decode(data)
	lang.AssertState(block != nil, "No PEM block found in %s", privateKeyPath)

	var priv *rsa.PrivateKey
	switch block.Type {
	case "RSA PRIVATE KEY":
		priv = optional.OfCommaErr(x509.ParsePKCS1PrivateKey(block.Bytes)).OrElsePanic("Cannot parse private key from %s", privateKeyPath)
	case "PRIVATE KEY":
		priv = optional.OfCommaErr(x509.ParsePKCS8PrivateKey(block.Bytes)).OrElsePanic("Cannot parse private key from %s", privateKeyPath).(*rsa.PrivateKey)
	default:
		panic(fmt.Sprintf("Unsupported key type %s", block.Type))
	}
	cipher, _ := base64.StdEncoding.DecodeString(value)
	decrypted := optional.OfCommaErr(rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, cipher, nil)).OrElsePanic("Cannot decrypt %s=%s", key, value)
	return string(decrypted)
}

func (s *RsaPropertySource) Properties() map[string]string {
	return nil
}

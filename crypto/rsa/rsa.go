package rsa

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"

	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/timing"
	"github.com/neurafuse/tools-go/vars"
)

func GenerateKeys(module, certPath string, printKeys bool) (string, string) {
	logging.Log([]string{"\n", vars.EmojiCrypto, vars.EmojiProcess}, "Creating selfsigned TLS certs for module: "+module, 0)
	// priv, err := rsa.GenerateKey(rand.Reader, *rsaBits)
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Failed during ecdsa.GenerateKey!", false, true, true)
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      getPk(module),
		NotBefore:    timing.GetCurrentTime(),
		NotAfter:     timing.GetCurrentTime().Add(timing.GetTimeDuration(24, "h") * 180),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	/*
	   hosts := strings.Split(*host, ",")
	   for _, h := range hosts {
	   	if ip := net.ParseIP(h); ip != nil {
	   		template.IPAddresses = append(template.IPAddresses, ip)
	   	} else {
	   		template.DNSNames = append(template.DNSNames, h)
	   	}
	   }

	   if *isCA {
	   	template.IsCA = true
	   	template.KeyUsage |= x509.KeyUsageCertSign
	   }
	*/

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Failed to create certificate!", false, true, true)
	out := &bytes.Buffer{}

	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if printKeys {
		fmt.Println(out.String())
	}
	if !filesystem.Exists(certPath) {
		filesystem.CreateDir(certPath, false)
	}
	publicKeyFileName := "public.crt"
	publicKeyFilePath := certPath + publicKeyFileName
	if filesystem.Exists(publicKeyFilePath) {
		filesystem.Delete(publicKeyFilePath, false)
	}
	filesystem.SaveByteArrayToFile(out.Bytes(), publicKeyFilePath)

	out.Reset()

	pem.Encode(out, pemBlockForKey(priv))
	if printKeys {
		fmt.Println(out.String())
	}
	privateKeyFileName := "private.key"
	privateKeyFilePath := certPath + privateKeyFileName
	if filesystem.Exists(privateKeyFilePath) {
		filesystem.Delete(privateKeyFilePath, false)
	}
	filesystem.SaveByteArrayToFile(out.Bytes(), privateKeyFilePath)
	logging.Log([]string{"", vars.EmojiCrypto, vars.EmojiSuccess}, "TLS certificates created.\n", 0)
	return publicKeyFilePath, privateKeyFilePath
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

func getPk(module string) pkix.Name {
	var pk pkix.Name
	if module == vars.NeuraKubeName {
		pk = pkix.Name{Organization: []string{vars.NeuraKubeName + " | " + vars.OrganizationName},
			Country:       []string{"Internet"},
			Province:      []string{"Internet"},
			Locality:      []string{"Internet"},
			StreetAddress: []string{"Internet"},
			PostalCode:    []string{"99999999"}}
	}
	return pk
}

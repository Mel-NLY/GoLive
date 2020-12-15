package main

/*
	Package ssl allows the creation of both cert.pem and key.pem files. Both of which are important for when developing a self-signed
	certificate for development usage. These files would eventually allow the server to run more securely by using HTTPS.
*/

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"
)

func main() {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)
	subject := pkix.Name{
		Organization:       []string{"BikeTransport Co."},
		OrganizationalUnit: []string{"BikeTransport"},
		CommonName:         "BikeTransport",
	}

	template := x509.Certificate{ //crypto/x509 is used to create certificate + Certificate struct is instantiated
		SerialNumber: serialNumber, //Randomly generated very large integer (Not CA assigned)
		Subject:      subject,      //Distinguished name
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour), // Validity period is 1 year from date cert is created
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}, //Indicate that X.509 cert is used for server authentication
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},             //Running from localhost
	}

	pk, _ := rsa.GenerateKey(rand.Reader, 2048) //Create RSA private key using cypto/rsa
	//RSA private key has a private key, used by x509.CreateCertificate func to create SSL cert
	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk) //Takes in Certificate struct, public and private keys to create a slice of DER-formatted bytes
	certOut, _ := os.Create("cert.pem")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}) //encoding/pem to encode the certificate into cert.pem file
	certOut.Close()

	keyOut, _ := os.Create("key.pem") //Use PEM encode and save key generated earlier into key.pem file
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	keyOut.Close()
}

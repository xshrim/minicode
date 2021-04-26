package pki

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	rd "math/rand"
	"net"
	"os"
	"path/filepath"
	"time"
)

type Attribute struct {
	CommonName         string
	Country            string
	Province           string
	Locality           string
	Organization       string
	OrganizationalUnit string
}

func NewAttr(cn, c, p, l, o, ou string) Attribute {
	return Attribute{
		CommonName:         cn,
		Country:            c,
		Province:           p,
		Locality:           l,
		Organization:       o,
		OrganizationalUnit: ou,
	}
}

type SSL struct {
	Host   string
	Expire time.Duration
	Attr   Attribute
	AltIP  []net.IP
	AltDNS []string
	Dir    string
}

func New(dir, host string, expire time.Duration, attr Attribute, altIP []net.IP, altDNS []string) *SSL {
	return &SSL{
		Host:   host,
		Expire: expire,
		Attr:   attr,
		AltIP:  altIP,
		AltDNS: altDNS,
		Dir:    dir,
	}
}

func (ssl *SSL) GenerateCertKey(parentCert *x509.Certificate, parentKey *rsa.PrivateKey, opts ...string) ([]byte, []byte, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	commonName := ssl.Attr.CommonName
	organization := ssl.Attr.Organization
	organizationalUnit := ssl.Attr.OrganizationalUnit
	expire := ssl.Expire

	if len(opts) > 0 {
		commonName = opts[0]
	}
	if len(opts) > 1 {
		organization = opts[1]
	}
	if len(opts) > 2 {
		organizationalUnit = opts[2]
	}
	if expire == 0 {
		expire = time.Hour * 24 * 365
	}

	if commonName == "" {
		return nil, nil, fmt.Errorf("CommonName can not be empty")
	}

	isCA := true
	sn := big.NewInt(0)
	usage := x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
	subject := pkix.Name{
		CommonName: commonName,
	}

	if ssl.Attr.Country != "" {
		subject.Country = []string{ssl.Attr.Country}
	}
	if ssl.Attr.Province != "" {
		subject.Province = []string{ssl.Attr.Province}
	}
	if ssl.Attr.Locality != "" {
		subject.Locality = []string{ssl.Attr.Locality}
	}
	if organization != "" {
		subject.Organization = []string{organization}
	}
	if organizationalUnit != "" {
		subject.OrganizationalUnit = []string{organizationalUnit}
	}

	if parentCert != nil {
		isCA = false
		sn = big.NewInt(rd.Int63())
		usage ^= x509.KeyUsageCertSign

		// subject.Country = []string{ssl.Attr.Country}
		// subject.Province = []string{ssl.Attr.Province}
		// subject.Locality = []string{ssl.Attr.Locality}
		// subject.Organization = []string{organization}
		// subject.OrganizationalUnit = []string{organizationalUnit}
	}

	template := x509.Certificate{
		SerialNumber: sn,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(expire),

		KeyUsage:              usage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  isCA,
	}

	if ssl.Host != "" {
		if ip := net.ParseIP(ssl.Host); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, ssl.Host)
		}
	}

	template.IPAddresses = append(template.IPAddresses, ssl.AltIP...)
	template.DNSNames = append(template.DNSNames, ssl.AltDNS...)

	if isCA {
		parentCert = &template
		parentKey = priv
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, parentCert, &priv.PublicKey, parentKey)
	if err != nil {
		return nil, nil, err
	}

	// Generate cert
	certBuffer := bytes.Buffer{}
	if err := pem.Encode(&certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return nil, nil, err
	}

	// Generate key
	keyBuffer := bytes.Buffer{}
	if err := pem.Encode(&keyBuffer, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return nil, nil, err
	}

	return certBuffer.Bytes(), keyBuffer.Bytes(), nil
}

func (ssl *SSL) GenerateSelfSignedCertKey(opts ...string) ([]byte, []byte, error) {
	return ssl.GenerateCertKey(nil, nil, opts...)
}

func ExtractCertKey(certBytes, keyBytes []byte) (*x509.Certificate, *rsa.PrivateKey, error) {
	certBlock, _ := pem.Decode(certBytes)

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	keyBlock, _ := pem.Decode(keyBytes)
	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	return cert, key, nil
}

func ExtractCertKeyFromFile(certFile, keyFile string) (*x509.Certificate, *rsa.PrivateKey, error) {
	certBytes, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, nil, err
	}

	keyBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, nil, err
	}

	return ExtractCertKey(certBytes, keyBytes)
}

func SaveCertKey(dir, name string, certBytes, keyBytes []byte) error {

	if err := ioutil.WriteFile(filepath.Join(dir, fmt.Sprintf("%s.crt", name)), certBytes, 0644); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(dir, fmt.Sprintf("%s.key", name)), keyBytes, 0644); err != nil {
		return err
	}

	return nil
}

func (ssl *SSL) GenerateKubernetesCertKeys(nodes, etcds []string, ca ...string) error {

	// https://feisky.gitbooks.io/kubernetes/content/deploy/kubernetes-the-hard-way/04-certificate-authority.html
	var err error
	var certBytes, keyBytes []byte
	if len(ca) > 1 {
		certBytes, err = ioutil.ReadFile(ca[0])
		if err != nil {
			return err
		}

		keyBytes, err = ioutil.ReadFile(ca[1])
		if err != nil {
			return err
		}
	} else {
		certBytes, keyBytes, err = ssl.GenerateSelfSignedCertKey("kubernetes", "Kubernetes", "System")
		if err != nil {
			return err
		}
	}

	_, err = os.Stat(ssl.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(ssl.Dir, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	// ca
	cacert, cakey, err := ExtractCertKey(certBytes, keyBytes)
	if err != nil {
		return err
	}
	if err := SaveCertKey(ssl.Dir, "ca", certBytes, keyBytes); err != nil {
		return err
	}

	// apiserver
	certBytes, keyBytes, err = ssl.GenerateCertKey(cacert, cakey, "apiserver", "system:apiserver", "System")
	if err != nil {
		return err
	}
	if err := SaveCertKey(ssl.Dir, "apiserver", certBytes, keyBytes); err != nil {
		return err
	}

	// admin
	certBytes, keyBytes, err = ssl.GenerateCertKey(cacert, cakey, "master", "system:master", "System")
	if err != nil {
		return err
	}
	if err := SaveCertKey(ssl.Dir, "admin", certBytes, keyBytes); err != nil {
		return err
	}

	// controller-manager
	certBytes, keyBytes, err = ssl.GenerateCertKey(cacert, cakey, "kube-controller-manager", "system:kube-controller-manager", "System")
	if err != nil {
		return err
	}
	if err := SaveCertKey(ssl.Dir, "controller-manager", certBytes, keyBytes); err != nil {
		return err
	}

	// scheduler
	certBytes, keyBytes, err = ssl.GenerateCertKey(cacert, cakey, "kube-scheduler", "system:kube-scheduler", "System")
	if err != nil {
		return err
	}
	if err := SaveCertKey(ssl.Dir, "scheduler", certBytes, keyBytes); err != nil {
		return err
	}

	// kube-proxy
	certBytes, keyBytes, err = ssl.GenerateCertKey(cacert, cakey, "kube-proxy", "system:kube-proxy", "System")
	if err != nil {
		return err
	}
	if err := SaveCertKey(ssl.Dir, "proxy", certBytes, keyBytes); err != nil {
		return err
	}

	// client
	certBytes, keyBytes, err = ssl.GenerateCertKey(cacert, cakey, "client", "system:client", "System")
	if err != nil {
		return err
	}
	if err := SaveCertKey(ssl.Dir, "client", certBytes, keyBytes); err != nil {
		return err
	}

	// kubelet
	host := ssl.Host
	for _, node := range nodes {
		ssl.Host = node
		certBytes, keyBytes, err = ssl.GenerateCertKey(cacert, cakey, "kubelet", "system:kubelet", "System")
		if err != nil {
			return err
		}

		if err := SaveCertKey(ssl.Dir, fmt.Sprintf("kubelet-%s", node), certBytes, keyBytes); err != nil {
			return err
		}
	}
	ssl.Host = host
	// TODO generate etcd credientals
	return nil
}

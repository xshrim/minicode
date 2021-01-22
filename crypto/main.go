package main

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/common/tools/cryptogen/ca"
)

const key = `
-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQguzroFhWe7egyDsVY
+1xe5Q+aELcsMX/ICG+nb8MLh0GhRANCAAR2iYiC/tG7NjOg39P3p3lF1j73fWbu
YmGYRrV+X/3TnSH7/+rLtXkfXsbppIkpMd6PfUwP2q/9l0AI1Qyk/4R8
-----END PRIVATE KEY-----
`

const rootPem = `
-----BEGIN CERTIFICATE-----
MIIEBDCCAuygAwIBAgIDAjppMA0GCSqGSIb3DQEBBQUAMEIxCzAJBgNVBAYTAlVT
MRYwFAYDVQQKEw1HZW9UcnVzdCBJbmMuMRswGQYDVQQDExJHZW9UcnVzdCBHbG9i
YWwgQ0EwHhcNMTMwNDA1MTUxNTU1WhcNMTUwNDA0MTUxNTU1WjBJMQswCQYDVQQG
EwJVUzETMBEGA1UEChMKR29vZ2xlIEluYzElMCMGA1UEAxMcR29vZ2xlIEludGVy
bmV0IEF1dGhvcml0eSBHMjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
AJwqBHdc2FCROgajguDYUEi8iT/xGXAaiEZ+4I/F8YnOIe5a/mENtzJEiaB0C1NP
VaTOgmKV7utZX8bhBYASxF6UP7xbSDj0U/ck5vuR6RXEz/RTDfRK/J9U3n2+oGtv
h8DQUB8oMANA2ghzUWx//zo8pzcGjr1LEQTrfSTe5vn8MXH7lNVg8y5Kr0LSy+rE
ahqyzFPdFUuLH8gZYR/Nnag+YyuENWllhMgZxUYi+FOVvuOAShDGKuy6lyARxzmZ
EASg8GF6lSWMTlJ14rbtCMoU/M4iarNOz0YDl5cDfsCx3nuvRTPPuj5xt970JSXC
DTWJnZ37DhF5iR43xa+OcmkCAwEAAaOB+zCB+DAfBgNVHSMEGDAWgBTAephojYn7
qwVkDBF9qn1luMrMTjAdBgNVHQ4EFgQUSt0GFhu89mi1dvWBtrtiGrpagS8wEgYD
VR0TAQH/BAgwBgEB/wIBADAOBgNVHQ8BAf8EBAMCAQYwOgYDVR0fBDMwMTAvoC2g
K4YpaHR0cDovL2NybC5nZW90cnVzdC5jb20vY3Jscy9ndGdsb2JhbC5jcmwwPQYI
KwYBBQUHAQEEMTAvMC0GCCsGAQUFBzABhiFodHRwOi8vZ3RnbG9iYWwtb2NzcC5n
ZW90cnVzdC5jb20wFwYDVR0gBBAwDjAMBgorBgEEAdZ5AgUBMA0GCSqGSIb3DQEB
BQUAA4IBAQA21waAESetKhSbOHezI6B1WLuxfoNCunLaHtiONgaX4PCVOzf9G0JY
/iLIa704XtE7JW4S615ndkZAkNoUyHgN7ZVm2o6Gb4ChulYylYbc3GrKBIxbf/a/
zG+FA1jDaFETzf3I93k9mTXwVqO94FntT0QJo544evZG0R0SnU++0ED8Vf4GXjza
HFa9llF7b1cq26KqltyMdMKVvvBulRP/F/A8rLIQjcxz++iPAsbw+zOzlTvjwsto
WHPbqCRiOwY1nQ2pM714A5AuTHhdUDqB1O6gyHA43LL5Z/qHQF1hwFGPa4NrzQU6
yuGnBXj8ytqU0CwIPX4WecigUCAkVDNx
-----END CERTIFICATE-----
`

const capem = `
LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNNekNDQWRtZ0F3SUJBZ0lRSnlneUdlZCsrV2NMdU1DOTY0N2kyVEFLQmdncWhrak9QUVFEQWpCa01Rc3cKQ1FZRFZRUUdFd0pEVGpFUU1BNEdBMVVFQ0JNSFFtVnBhbWx1WnpFUU1BNEdBMVVFQnhNSFFtVnBhbWx1WnpFVwpNQlFHQTFVRUNoTU5ZV1J0YVc1MFpYTjBiM0puTVRFWk1CY0dBMVVFQXhNUVkyRXVZV1J0YVc1MFpYTjBiM0puCk1UQWVGdzB4T1RFeU1ETXdPVE16TURCYUZ3MHlPVEV4TXpBd09UTXpNREJhTUdReEN6QUpCZ05WQkFZVEFrTk8KTVJBd0RnWURWUVFJRXdkQ1pXbHFhVzVuTVJBd0RnWURWUVFIRXdkQ1pXbHFhVzVuTVJZd0ZBWURWUVFLRXcxaApaRzFwYm5SbGMzUnZjbWN4TVJrd0Z3WURWUVFERXhCallTNWhaRzFwYm5SbGMzUnZjbWN4TUZrd0V3WUhLb1pJCnpqMENBUVlJS29aSXpqMERBUWNEUWdBRStxaWI2TDI3b0hKOWNtUytXZFhXWm1wc3YyaVROcU9iV2g5VExmYlgKWDRLazhsUHlGV0gwMXNqV0I3NUdUbWp0U0dERzQ3SERCek9MRUVkeVQzc3ZXYU50TUdzd0RnWURWUjBQQVFILwpCQVFEQWdHbU1CMEdBMVVkSlFRV01CUUdDQ3NHQVFVRkJ3TUNCZ2dyQmdFRkJRY0RBVEFQQmdOVkhSTUJBZjhFCkJUQURBUUgvTUNrR0ExVWREZ1FpQkNBeENVUmFzdUtoeHh0SlNoRXdvWHQ0Z1loK0hYYzF2TmRjTTdTVFl6MXYKTERBS0JnZ3Foa2pPUFFRREFnTklBREJGQWlFQTdKbUNlN0lOYjA1Y3NZUlRjOVkxUnJMdUJDbEZWWFY2ckhSVApoUzY2dkw0Q0lDZjJWZ2kxemoraytWNGhpdTZ1bTdsaDdvMGVTRjljZ3gzMlJ2QTVUZnp1Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
`
const adminpem = `
LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNEekNDQWJXZ0F3SUJBZ0lRYVZBRVJ4NkFDM3RHTHF2UlRHa3lGekFLQmdncWhrak9QUVFEQWpCa01Rc3cKQ1FZRFZRUUdFd0pEVGpFUU1BNEdBMVVFQ0JNSFFtVnBhbWx1WnpFUU1BNEdBMVVFQnhNSFFtVnBhbWx1WnpFVwpNQlFHQTFVRUNoTU5ZV1J0YVc1MFpYTjBiM0puTVRFWk1CY0dBMVVFQXhNUVkyRXVZV1J0YVc1MFpYTjBiM0puCk1UQWVGdzB4T1RFeU1ETXdPVE16TURCYUZ3MHlPVEV4TXpBd09UTXpNREJhTUdBeEN6QUpCZ05WQkFZVEFrTk8KTVJBd0RnWURWUVFJRXdkQ1pXbHFhVzVuTVJBd0RnWURWUVFIRXdkQ1pXbHFhVzVuTVE4d0RRWURWUVFMRXdaagpiR2xsYm5ReEhEQWFCZ05WQkFNTUUwRmtiV2x1UUdGa2JXbHVkR1Z6ZEc5eVp6RXdXVEFUQmdjcWhrak9QUUlCCkJnZ3Foa2pPUFFNQkJ3TkNBQVFQd1BSU25jMDl0ZGd6Tm4zSUVvMkh0SzFVNEpDU01nRnVHc3EwclBsK1BzNkMKZlBrMml0WjBsa3BlZlJUb1N5L3F2bVI3TlA4clk2OTlOejdERnVXSm8wMHdTekFPQmdOVkhROEJBZjhFQkFNQwpCNEF3REFZRFZSMFRBUUgvQkFJd0FEQXJCZ05WSFNNRUpEQWlnQ0F4Q1VSYXN1S2h4eHRKU2hFd29YdDRnWWgrCkhYYzF2TmRjTTdTVFl6MXZMREFLQmdncWhrak9QUVFEQWdOSUFEQkZBaUVBenJXMWxmL0RFMkdKcEY5WGMvWUIKTGZvbytaeU9VTGFWL2hoYVlnRmplNHdDSUIyMHNTWnRrOTNCYWg5R3VGR1NkVEFJV0JZWkJVNzdhQ2ZJdkticgpPQm9rCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
`

func getPem(fpath string) string {
	bytes, _ := ioutil.ReadFile(fpath)

	return string(bytes)
}

func getCA(fpath string) *ca.CA {
	block, _ := pem.Decode([]byte(getPem(fpath)))

	cert, err := x509.ParseCertificate(block.Bytes)
	fmt.Println(err)

	fmt.Println(cert.Subject)
	signca := &ca.CA{
		Name:               cert.Subject.CommonName,
		Country:            strings.Join(cert.Subject.Country, " "),
		Province:           strings.Join(cert.Subject.Province, " "),
		Locality:           strings.Join(cert.Subject.Locality, " "),
		OrganizationalUnit: strings.Join(cert.Subject.OrganizationalUnit, " "),
		StreetAddress:      strings.Join(cert.Subject.StreetAddress, " "),
		PostalCode:         strings.Join(cert.Subject.PostalCode, " "),
		SignCert:           cert,
	}

	return signca
}

func verfityCert(certPem string, rootPems []string) bool {
	roots := x509.NewCertPool()
	for _, pem := range rootPems {
		roots.AppendCertsFromPEM([]byte(pem))
	}

	block, _ := pem.Decode([]byte(certPem))
	if block == nil {
		return false
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false
	}

	opts := x509.VerifyOptions{
		//DNSName: "mail.google.com",
		Roots: roots,
	}

	if _, err := cert.Verify(opts); err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func main() {
	// signca := getCA("./msp/cacerts/ca.adminechainorg1-cert.pem")
	// fmt.Println(signca)

	// tlsca := getCA("./msp/tlscacerts/tlsca.adminechainorg1-cert.pem")

	// admin := getPem("./msp/admincerts/Admin@adminechainorg1-cert.pem")

	// fmt.Println(tlsca.SignCert.Subject.Organization)

	// fmt.Println(verfityCert(admin, []string{getPem("./msp/cacerts/ca.adminechainorg1-cert.pem")}))
	// fmt.Println(verfityCert(admin, []string{rootPem}))

	// encPem := base64.StdEncoding.EncodeToString([]byte(rootPem))
	// fmt.Println(encPem)
	// decPem, _ := base64.StdEncoding.DecodeString(encPem)

	// fmt.Println(rootPem == string(decPem))

	//msp.GenerateVerifyingMSP("./test", signca, tlsca, true)

	// ca.NewCA("./ca", "org1", signca.Name, signca.Country, signca.Province, signca.Locality, signca.OrganizationalUnit, signca.StreetAddress, signca.PostalCode)

	// msp.GenerateLocalMSP("./peer", "peer0")

	f, err := filepath.Abs("./cc")
	fmt.Println(f, err)

	keyDERBlock, keyPEMBlock := pem.Decode([]byte(key))
	fmt.Println(keyDERBlock, keyPEMBlock)
	// key, _ := x509.ParseECPrivateKey(keyDERBlock.Bytes)
	// pubkey := key.Public()
	// fmt.Println(pubkey)

	opts := &factory.FactoryOpts{
		ProviderName: "SW",
		SwOpts: &factory.SwOpts{
			HashFamily: "SHA2",
			SecLevel:   256,
		},
	}

	csp, err := factory.GetBCCSPFromOpts(opts)
	if err != nil {
		fmt.Println(err)
	}

	priv, err := csp.KeyImport(keyDERBlock.Bytes, &bccsp.ECDSAPrivateKeyImportOpts{Temporary: true})

	pub, _ := priv.PublicKey()

	fmt.Println(hex.EncodeToString(pub.SKI()))

	capbytes, _ := base64.StdEncoding.DecodeString(capem)
	adpbytes, _ := base64.StdEncoding.DecodeString(adminpem)
	fmt.Println(string(capbytes))
	fmt.Println(string(adpbytes))
	fmt.Println(verfityCert(string(adpbytes), []string{string(capbytes)}))

	fmt.Println("==================================================")
	//fmt.Println(verfityCert(getPem("./msp/admincerts/Admin@adminechainorg1-cert.pem"), []string{getPem("./msp/cacerts/ca.adminechainorg1-cert.pem")}))

	fmt.Println(verfityCert(getPem("./admin.pem"), []string{getPem("./ca.pem")}))

}

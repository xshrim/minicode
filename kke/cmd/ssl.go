/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"crypto/rsa"
	"crypto/x509"
	"kke/pkg/pki"
	"kke/pkg/utils"
	"net"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var cacert string
var cakey string
var subject string
var expire string
var host string
var altIP string
var altDNS string
var dir string
var name string
var node string
var etcd string

// sslCmd represents the ssl command
var sslCmd = &cobra.Command{
	Use:   "ssl",
	Short: "generate ssl credentials",
	Args:  cobra.MinimumNArgs(1),
}

var caCmd = &cobra.Command{
	Use:   "ca",
	Short: "generate selfsigned ssl credentials",
	Run: func(cmd *cobra.Command, args []string) {
		ssl := initSSL()
		certBytes, keyBytes, err := ssl.GenerateSelfSignedCertKey()
		if err != nil {
			panic(err)
		}

		if name == "" {
			name = "ca"
		}

		if err := pki.SaveCertKey(ssl.Dir, name, certBytes, keyBytes); err != nil {
			panic(err)
		}
	},
}

var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "generate ssl private key and sign certificate using ca credentials",
	Run: func(cmd *cobra.Command, args []string) {
		ssl := initSSL()

		var err error
		var caCert *x509.Certificate
		var caKey *rsa.PrivateKey
		if cacert != "" && cakey != "" {
			caCert, caKey, err = pki.ExtractCertKeyFromFile(cacert, cakey)
			if err != nil {
				panic(err)
			}
		} else {
			certBytes, keyBytes, err := ssl.GenerateSelfSignedCertKey()
			if err != nil {
				panic(err)
			}

			if err := pki.SaveCertKey(ssl.Dir, "ca", certBytes, keyBytes); err != nil {
				panic(err)
			}

			caCert, caKey, err = pki.ExtractCertKey(certBytes, keyBytes)
			if err != nil {
				panic(err)
			}
		}

		certBytes, keyBytes, err := ssl.GenerateCertKey(caCert, caKey)
		if err != nil {
			panic(err)
		}

		if name == "" {
			name = "server"
		}

		if err := pki.SaveCertKey(ssl.Dir, name, certBytes, keyBytes); err != nil {
			panic(err)
		}
	},
}

var kubeCmd = &cobra.Command{
	Use:   "kube",
	Short: "generate all credentials for kubernetes cluster deployment",
	Run: func(cmd *cobra.Command, args []string) {
		ssl := initSSL()

		nodes := utils.StringSplit(node)
		etcds := utils.StringSplit(etcd)

		var err error
		if cacert != "" && cakey != "" {
			err = ssl.GenerateKubernetesCertKeys(nodes, etcds, cacert, cakey)
		} else {
			err = ssl.GenerateKubernetesCertKeys(nodes, etcds)
		}
		if err != nil {
			panic(err)
		}
	},
}

func initSSL() *pki.SSL {
	expireDuration, err := time.ParseDuration(expire)
	if err != nil {
		panic(err)
	}

	var altIPs []net.IP
	altAddrs := utils.StringSplit(altIP)
	for _, altAddr := range altAddrs {
		altIPs = append(altIPs, net.ParseIP(altAddr))
	}
	altIPs = append(altIPs, net.ParseIP("127.0.0.1"))

	altDNSs := utils.StringSplit(altDNS)

	var attr pki.Attribute

	subjects := utils.StringSplit(subject)
	for _, subj := range subjects {
		kv := strings.Split(subj, "=")
		if len(kv) != 2 {
			panic("subject argument format invalid")
		}

		switch kv[0] {
		case "C":
			attr.Country = kv[1]
		case "P":
			attr.Province = kv[1]
		case "L":
			attr.Locality = kv[1]
		case "O":
			attr.Organization = kv[1]
		case "OU":
			attr.OrganizationalUnit = kv[1]
		case "CN":
			attr.CommonName = kv[1]
		}
	}

	return pki.New(dir, host, expireDuration, attr, altIPs, altDNSs)
}

func init() {
	sslCmd.PersistentFlags().StringVarP(&cacert, "cacert", "C", "", "certificate file path of Certification Authority")
	sslCmd.PersistentFlags().StringVarP(&cakey, "cakey", "K", "", "private key file path of Certification Authority")
	sslCmd.PersistentFlags().StringVarP(&subject, "subject", "S", "C=CN,P=BeiJing,L=BeiJing,O=System,OU=System,CN=Cert", "subject information of certificate")
	sslCmd.PersistentFlags().StringVarP(&expire, "expire", "E", "87600h", "expire time of certificate")
	sslCmd.PersistentFlags().StringVarP(&host, "host", "H", "", "server ip/domain the certificate signed for")
	sslCmd.PersistentFlags().StringVarP(&altIP, "ips", "I", "", "alternative ip addresses of certificate")
	sslCmd.PersistentFlags().StringVarP(&altDNS, "dnss", "N", "", "alternative dns domains of certificate")
	sslCmd.PersistentFlags().StringVarP(&dir, "dir", "D", "./", "save folder of the generated credentials")

	caCmd.Flags().StringVarP(&name, "name", "n", "", "base filename of the generated credentials")
	certCmd.Flags().StringVarP(&name, "name", "n", "", "base filename of the generated credentials")

	kubeCmd.Flags().StringVarP(&node, "nodes", "n", "", "node list in kubernetes cluster")
	kubeCmd.Flags().StringVarP(&etcd, "etcds", "e", "", "etcd member list in kubernetes etcd cluster")

	sslCmd.AddCommand(caCmd)
	sslCmd.AddCommand(certCmd)
	sslCmd.AddCommand(kubeCmd)
	rootCmd.AddCommand(sslCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sslCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sslCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

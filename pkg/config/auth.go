package config

import (
	"crypto/x509"
	"io/ioutil"
)

func LoadToken(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data[:]), nil
}

func LoadCA(caFile string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	if ca, err := ioutil.ReadFile(caFile); err != nil {
		return nil, err
	} else {
		pool.AppendCertsFromPEM(ca)
	}
	return pool, nil
}

package util

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func ParseCert(fullchain []byte) (*x509.Certificate, error) {
	certBlock, _ := pem.Decode(fullchain)
	if certBlock == nil {
		err := errors.New("Decode.ERROR：fullchain")
		return nil, err
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		err = errors.New("ParseCertificate.ERROR：" + err.Error())
		return nil, err
	}
	return cert, nil
}

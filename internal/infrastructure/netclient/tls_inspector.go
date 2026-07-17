package netclient

import (
	"net/http"
	"goaconly/internal/domain"
	"time"
)

// SSLInfo содержит данные о TLS-сертификате сайта.
type SSLInfo struct {
	ExpiresAt time.Time
	DaysLeft  int
	Issuer    string
	Subject   string
}

// ExtractSSLInfo извлекает данные SSL-сертификата из TLS-соединения.
func ExtractSSLInfo(resp *http.Response) (*SSLInfo, error) {
	if resp.TLS == nil {
		return nil, nil
	}
	if len(resp.TLS.PeerCertificates) == 0 {
		return nil, domain.ErrNoCertificate
	}
	cert := resp.TLS.PeerCertificates[0]
	daysLeft := int(time.Until(cert.NotAfter).Hours() / 24)

	return &SSLInfo{
		ExpiresAt: cert.NotAfter,
		DaysLeft:  daysLeft,
		Issuer:    cert.Issuer.CommonName,
		Subject:   cert.Subject.CommonName,
	}, nil
}

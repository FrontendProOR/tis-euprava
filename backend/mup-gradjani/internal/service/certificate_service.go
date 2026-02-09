package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"tis-euprava/mup-gradjani/internal/domain"
	"tis-euprava/mup-gradjani/internal/repository"
)

type CertificateService struct {
	certs    repository.CertificateRepository
	requests repository.RequestRepository
	payments repository.PaymentRepository
}

func NewCertificateService(certs repository.CertificateRepository, reqs repository.RequestRepository, pays repository.PaymentRepository) *CertificateService {
	return &CertificateService{certs: certs, requests: reqs, payments: pays}
}

// GenerateCertificate creates (or returns existing) certificate PDF for a request.
// Rules:
// - request must be APPROVED
// - payment must be PAID
func (s *CertificateService) GenerateCertificate(requestID string) ([]byte, error) {
	req, err := s.requests.FindByID(requestID)
	if err != nil {
		return nil, errors.New("request not found")
	}
	if req.Status != domain.RequestApproved {
		return nil, errors.New("request must be APPROVED")
	}

	pay, err := s.payments.FindByRequestID(requestID)
	if err != nil {
		return nil, errors.New("payment not found")
	}
	if pay == nil || pay.Status != domain.PaymentPaid {
		return nil, errors.New("payment must be PAID")
	}

	// if exists, return it
	if existing, err := s.certs.FindByRequestID(requestID); err == nil && existing != nil {
		return existing.PDF, nil
	}

	id := uuid.NewString()
	issuedAt := time.Now().UTC()

	pdf := buildMinimalPDF([]string{
		"REPUBLIKA SRBIJA",
		"MINISTARSTVO UNUTRASNJIH POSLOVA",
		"",
		"ELEKTRONSKO UVERENJE",
		"",
		"Ime i prezime: " + req.CitizenID, // ako nemaš ime, OK je i citizenID
		"Vrsta zahteva: " + req.Type,
		"Status zahteva: " + string(req.Status),
		"",
		"Datum izdavanja: " + issuedAt.Format("02.01.2006"),
		"Izdavalac: MUP Republike Srbije",
		"",
		"Ovo uverenje je izdato u elektronskom obliku.",
	})

	cert := &domain.ElectronicCertificate{
		ID:        id,
		RequestID: requestID,
		IssuedAt:  issuedAt,
		PDF:       pdf,
	}

	if err := s.certs.Create(cert); err != nil {
		return nil, err
	}
	return pdf, nil
}

// buildMinimalPDF creates a simple, valid PDF with Helvetica text only.
func buildMinimalPDF(lines []string) []byte {
	esc := func(s string) string {
		s = strings.ReplaceAll(s, "\\", "\\\\")
		s = strings.ReplaceAll(s, "(", "\\(")
		s = strings.ReplaceAll(s, ")", "\\)")
		return s
	}

	var content strings.Builder
	content.WriteString("BT\n/F1 12 Tf\n")
	y := 760

	for _, line := range lines {
		// Tm postavlja APSOLUTNU poziciju (x=50, y=y) svaki put -> tekst uvek ostaje na strani
		content.WriteString(fmt.Sprintf("1 0 0 1 50 %d Tm (%s) Tj\n", y, esc(line)))
		y -= 18
	}

	content.WriteString("ET\n")
	stream := content.String()

	objs := make([]string, 0, 6)
	objs = append(objs, "1 0 obj<< /Type /Catalog /Pages 2 0 R >>endobj\n")
	objs = append(objs, "2 0 obj<< /Type /Pages /Kids [3 0 R] /Count 1 >>endobj\n")
	objs = append(objs, "3 0 obj<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595 842] /Resources<< /Font<< /F1 4 0 R >> >> /Contents 5 0 R >>endobj\n")
	objs = append(objs, "4 0 obj<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>endobj\n")
	objs = append(objs, fmt.Sprintf("5 0 obj<< /Length %d >>stream\n%s\nendstream\nendobj\n", len(stream), stream))

	header := "%PDF-1.4\n"
	var out strings.Builder
	out.WriteString(header)

	offsets := make([]int, 0, len(objs)+1)
	offsets = append(offsets, 0)
	cur := len(header)
	for _, o := range objs {
		offsets = append(offsets, cur)
		out.WriteString(o)
		cur += len(o)
	}

	xrefPos := cur
	out.WriteString("xref\n0 6\n")
	out.WriteString("0000000000 65535 f \n")
	for i := 1; i <= 5; i++ {
		out.WriteString(fmt.Sprintf("%010d 00000 n \n", offsets[i]))
	}
	out.WriteString("trailer<< /Size 6 /Root 1 0 R >>\n")
	out.WriteString("startxref\n")
	out.WriteString(fmt.Sprintf("%d\n", xrefPos))
	out.WriteString("%%EOF\n")
	return []byte(out.String())
}

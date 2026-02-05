package testutils

type MockEmailService struct {
	Sent []string
}

func (m *MockEmailService) Send(to, subject, body string) error {
	m.Sent = append(m.Sent, to)
	return nil
}

package utils

import (
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockMailer struct {
	mock.Mock
}

func (m *mockMailer) SendMail(to []string, subject string, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

func (m *mockMailer) SendMailWithTemplate(to []string, subject string, templateFileName string, data interface{}) error {
	args := m.Called(to, subject, templateFileName, data)
	return args.Error(0)
}

func (m *mockMailer) QueueMail(to []string, subject string, body string) {
	m.Called(to, subject, body)
}

func (m *mockMailer) QueueMailWithTemplate(to []string, subject string, templateFileName string, data interface{}) {
	m.Called(to, subject, templateFileName, data)
}

func TestSendMail(t *testing.T) {
	m := new(mockMailer)
	to := []string{faker.Email(), faker.Email()}
	subject := faker.Sentence()
	body := faker.Paragraph()

	m.On("SendMail", to, subject, body).Return(nil)

	err := m.SendMail(to, subject, body)
	m.AssertExpectations(t)
	assert.Nil(t, err, "Error should be null")
}

func TestSendMailWithTemplate(t *testing.T) {
	m := new(mockMailer)
	to := []string{faker.Email(), faker.Email()}
	subject := faker.Sentence()
	templateFileName := faker.Word() + ".html"
	data := map[string]interface{}{"name": faker.Name()}

	m.On("SendMailWithTemplate", to, subject, templateFileName, data).Return(nil)

	err := m.SendMailWithTemplate(to, subject, templateFileName, data)
	m.AssertExpectations(t)
	assert.Nil(t, err, "Error should be null")
}

package email

import (
	"errors"
)

type Implementation int

const (
	Gomail = Implementation(iota)
)

type Email interface {
	Send(sendEmailTo *MailDetail) error
	Close() error
}

type MailDetail struct {
	Attachment []AttachmentDetail `json:"attachment"`
	SendTo     []string           `valid:"Required" json:"send_to"`
	SendToCc   []string           `json:"send_to_cc"`
	Title      string             `valid:"Required" json:"title"`
	Body       string             `valid:"Required" json:"body"`
	ID         string             `json:"id"`
}

type AttachmentDetail struct {
	Filename string `json:"filename"`
	Byte     []byte `json:"file_byte"`
}

type Config struct {
	Host        string `valid:"Required;Host" json:"host"`
	Port        int    `valid:"Required;Port" json:"port"`
	SenderEmail string `valid:"Required;SenderEmail" json:"sender_email"`
	Username    string `valid:"Required;Username" json:"username"`
	Password    string `valid:"Required" json:"password"`
}

// New email return email handler struct
func NewEmail(impl Implementation, emailConfig Config) (Email, error) {
	if Gomail == impl {
		return newGomail(&emailConfig)
	}

	return nil, errors.New("no email implementations found")
}

package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"regexp"

	"github.com/emersion/go-message/mail"
	"github.com/emersion/go-smtp"
)

var referencesRegex = regexp.MustCompile(`<([^>]+)>`)

type HandlerFunc func(*EmailMessage)

type Session struct {
	msg     *EmailMessage
	handler HandlerFunc
}

func NewSession(handler HandlerFunc) *Session {
	return &Session{
		msg:     &EmailMessage{},
		handler: handler,
	}
}

func (s *Session) Reset() {
	s.msg = &EmailMessage{}
}

func (s *Session) Logout() error {
	s.msg = nil
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	fromAddr, err := mail.ParseAddress(from)
	if err != nil {
		return err
	}

	s.msg.Addresses.From = &EmailAddress{
		Name:    fromAddr.Name,
		Address: fromAddr.Address,
	}

	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	toAddr, err := mail.ParseAddress(to)
	if err != nil {
		return err
	}

	s.msg.Addresses.To = &EmailAddress{
		Name:    toAddr.Name,
		Address: toAddr.Address,
	}

	return nil
}

func (s *Session) Data(r io.Reader) error {
	email, err := mail.CreateReader(r)
	if err != nil {
		return err
	}

	if messageId, err := email.Header.MessageID(); err != nil {
		return err
	} else {
		s.msg.ID = messageId
	}

	if subject, err := email.Header.Subject(); err != nil {
		return err
	} else {
		s.msg.Subject = subject
	}

	if date, err := email.Header.Date(); err != nil {
		return err
	} else {
		s.msg.Date = date.Format("2006-01-02T15:04:05Z07:00")
	}

	for _, m := range referencesRegex.FindAllStringSubmatch(email.Header.Get("References"), -1) {
		if len(m) > 1 {
			s.msg.References = append(s.msg.References, m[1])
		}
	}

	if addrs, err := email.Header.AddressList("Reply-To"); err != nil {
		return err
	} else if addrs != nil {
		s.msg.Addresses.ReplyTo = make([]*EmailAddress, len(addrs))
		for i, addr := range addrs {
			s.msg.Addresses.ReplyTo[i] = &EmailAddress{Name: addr.Name, Address: addr.Address}
		}
	}

	if addrs, err := email.Header.AddressList("Cc"); err != nil {
		return err
	} else if addrs != nil {
		s.msg.Addresses.Cc = make([]*EmailAddress, len(addrs))
		for i, addr := range addrs {
			s.msg.Addresses.Cc[i] = &EmailAddress{Name: addr.Name, Address: addr.Address}
		}
	}

	if bccAddrs, err := mail.ParseAddressList(email.Header.Get("Bcc")); err == nil {
		s.msg.Addresses.Bcc = make([]*EmailAddress, len(bccAddrs))
		for i, addr := range bccAddrs {
			s.msg.Addresses.Bcc[i] = &EmailAddress{Name: addr.Name, Address: addr.Address}
		}
	}

	for {
		part, err := email.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		switch h := part.Header.(type) {
		case *mail.InlineHeader:
			body, err := io.ReadAll(part.Body)
			if err != nil {
				return err
			}

			s.msg.Body = string(body)
		case *mail.AttachmentHeader:
			filename, err := h.Filename()
			if err != nil {
				return err
			}

			if filename == "" {
				return fmt.Errorf("missing filename")
			}

			data, err := io.ReadAll(part.Body)
			if err != nil {
				return err
			}

			s.msg.Attachments = append(s.msg.Attachments, &EmailAttachment{
				Filename: filename,
				Data:     base64.StdEncoding.EncodeToString(data),
			})
		}

	}

	s.handler(s.msg)
	return nil
}

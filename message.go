package smtp2http

type EmailAddress struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
}

type EmailAttachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Data        string `json:"data"`
}

type EmailEmbeddedFile struct {
	CID         string `json:"cid"`
	ContentType string `json:"content_type"`
	Data        string `json:"data"`
}

type EmailBody struct {
	Text string `json:"text,omitempty"`
	HTML string `json:"html,omitempty"`
}

type EmailAddresses struct {
	From    *EmailAddress   `json:"from"`
	To      *EmailAddress   `json:"to"`
	ReplyTo []*EmailAddress `json:"reply_to,omitempty"`
	Cc      []*EmailAddress `json:"cc,omitempty"`
	Bcc     []*EmailAddress `json:"bcc,omitempty"`
}

type EmailMessage struct {
	References []string `json:"references,omitempty"`

	ID      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Subject string `json:"subject,omitempty"`

	Body      string         `json:"body"`
	Addresses EmailAddresses `json:"addresses"`

	Attachments   []*EmailAttachment   `json:"attachments,omitempty"`
	EmbeddedFiles []*EmailEmbeddedFile `json:"embedded_files,omitempty"`
}

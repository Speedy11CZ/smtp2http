# SMTP2HTTP
Lightweight SMTP to HTTP bridge written in Go, created primarily for instant email webhooks on N8N.

## Usage
```
Usage: smtp2http [ ... ]

Parameters:

  -smtp.allow-insecure-auth
        Allow insecure authentication mechanisms (default true)
  -smtp.domain string
        SMTP server domain name (default "localhost")
  -smtp.listen-address string
        Address to bind the SMTP server to (default ":25")
  -smtp.max-message-bytes int
        SMTP server maximum message size in bytes (default 1073741824)
  -smtp.max-recipients int
        SMTP server maximum number of recipients per message (default 50)
  -smtp.read-timeout duration
        SMTP server read timeout (default 10s)
  -smtp.write-timeout duration
        SMTP server write timeout (default 10s)
  -webhook.url string
        URL format of the webhook to send emails to
```

# License
(c) Michal Spišak, 2025. Licensed under [MIT](https://github.com/Speedy11CZ/smtp2http/blob/main/LICENSE) license.


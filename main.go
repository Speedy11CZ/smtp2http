package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/rs/zerolog/log"
)

var (
	smtpListenAddress     = flag.String("smtp.listen-address", ":25", "Address to bind the SMTP server to")
	smtpDomain            = flag.String("smtp.domain", "localhost", "SMTP server domain name")
	smtpReadTimeout       = flag.Duration("smtp.read-timeout", 10*time.Second, "SMTP server read timeout")
	smtpWriteTimeout      = flag.Duration("smtp.write-timeout", 10*time.Second, "SMTP server write timeout")
	smtpMaxMessageBytes   = flag.Int64("smtp.max-message-bytes", 1024*1024*1024, "SMTP server maximum message size in bytes")
	smtpMaxRecipients     = flag.Int("smtp.max-recipients", 50, "SMTP server maximum number of recipients per message")
	smtpAllowInsecureAuth = flag.Bool("smtp.allow-insecure-auth", true, "Allow insecure authentication mechanisms")
	webhookUrl            = flag.String("webhook.url", "", "URL format of the webhook to send emails to")
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [ ... ]\n\nParameters:\n", os.Args[0])
		fmt.Println()
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if *webhookUrl == "" {
		fmt.Println("Error: webhook.url parameter is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	smtpServer := smtp.NewServer(smtp.BackendFunc(func(c *smtp.Conn) (smtp.Session, error) {
		return NewSession(func(msg *EmailMessage) {
			pr, pw := io.Pipe()
			encoder := json.NewEncoder(pw)
			encoder.SetEscapeHTML(false)
			encoder.SetIndent("", "  ")

			go func() {
				if err := encoder.Encode(msg); err != nil {
					log.Printf("failed to encode message: %v", err)
				}

				pw.Close()
			}()

			resp, err := http.DefaultClient.Post(fmt.Sprintf(*webhookUrl, msg.Addresses.To.Address), "application/json", pr)
			if err != nil {
				log.Error().Err(err).Msg("failed to send message to webhook")
				return
			}

			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			log.Debug().Msgf("message sent to webhook, response status: %s", resp.Status)
		}), nil
	}))

	smtpServer.Addr = *smtpListenAddress
	smtpServer.Domain = *smtpDomain
	smtpServer.WriteTimeout = *smtpWriteTimeout
	smtpServer.ReadTimeout = *smtpReadTimeout
	smtpServer.MaxMessageBytes = *smtpMaxMessageBytes
	smtpServer.MaxRecipients = *smtpMaxRecipients
	smtpServer.AllowInsecureAuth = *smtpAllowInsecureAuth

	log.Info().Msgf("starting SMTP server on %s", *smtpListenAddress)
	if err := smtpServer.ListenAndServe(); err != nil {
		if err == smtp.ErrServerClosed {
			log.Info().Msg("SMTP server closed")
			return
		}

		log.Fatal().Err(err).Msg("failed to start SMTP server")
	}
}

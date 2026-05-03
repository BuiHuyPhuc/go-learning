package sendto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type MailRequest struct {
	ToEmail     string `json:"toEmail"`
	MessageBody string `json:"messageBody"`
	Subject     string `json:"subject"`
	Attachment  string `json:"attachment"`
}

func SendEmailOtpByAPI(otp, email, purpose string) error {
	postURL := "..."

	mailRequest := MailRequest{
		ToEmail:     email,
		MessageBody: "OTP is " + otp,
		Subject:     "Verify OTP to " + purpose,
		Attachment:  "path/to/email",
	}

	requestBody, err := json.Marshal(mailRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", postURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Printf("Response status: %d\n", resp.StatusCode)

	return nil
}

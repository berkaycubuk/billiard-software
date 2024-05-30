package brevo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

/*
Email notifications use Brevo.
*/
func SendEmail(email string, name string, subject string, content string) error {
				payload := []byte(fmt.Sprintf(`{
								"sender": {
												"name": "",
												"email": ""
								},
								"to": [
												{
																"email": "%s",
																"name": "%s"
												}
								],
								"subject": "%s",
								"htmlContent": "%s"
				}`, email, name, subject, content))

				apiURL := "https://api.brevo.com/v3/smtp/email"
				apiKey := os.Getenv("BREVO_API_KEY")
				req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
				if err != nil {
								return err
				}

				req.Header.Set("content-type", "application/json")
				req.Header.Set("api-key", apiKey)

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
								return err
				}

				responseBody, err := ioutil.ReadAll(resp.Body)
				if err != nil {
								fmt.Println("Error reading response body:", err)
								return err
				}

				fmt.Println(string(responseBody))

				defer resp.Body.Close()

				return nil
}

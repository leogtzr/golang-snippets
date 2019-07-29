// using SendGrid's Go Library
// https://github.com/sendgrid/sendgrid-go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func main() {

	from := mail.NewEmail("Leonidas", "leonidas@root.com")
	subject := "Automated email from laptop ... "
	to := mail.NewEmail("Leo Gtz", "leogutierrezramirez@gmail.com")
	plainTextContent, _ := ioutil.ReadAll(os.Stdin)
	plainTextContent = []byte(strings.ReplaceAll(string(plainTextContent), "\n", "<br>"))
	// &nbsp;
	// &Tab;
	plainTextContent = []byte(strings.ReplaceAll(string(plainTextContent), "\t", "&Tab;"))
	plainTextContent = []byte(strings.ReplaceAll(string(plainTextContent), " ", "&nbsp;"))
	htmlContent := plainTextContent
	message := mail.NewSingleEmail(from, subject, to, string(plainTextContent), string(htmlContent))
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println("Alv ... ")
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

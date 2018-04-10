package models

import (
	"bytes"
	"github.com/rakutentech/dotconf-assist/src/backend/settings"
	"html/template"
	"net/smtp"
)

type TemplateData struct {
	To             string
	Who            string
	Subject        string
	Body           string
	Link           string
	AdminEmail     string
	Data           interface{}
	AdditionalData interface{}
}

func SendEmail(from, to, target, who, function, subject, link string, Data, AdditionalData interface{}) error {
	var templateData TemplateData
	var err error
	templateData.To = to
	templateData.Who = who
	templateData.Subject = subject
	templateData.Link = link
	templateData.AdminEmail = to
	templateData.Data = Data
	templateData.AdditionalData = AdditionalData

	// settings.WriteDebugLog(templateData.Data)

	html := "src/assets/templates/email-" + target + "-" + function + ".html"
	if templateData.Body, err = ParseTemplate(html, templateData); err != nil {
		return err
	}

	var conf = settings.GetConfig()
	c, err := smtp.Dial(conf.MailServer + ":25")
	if err != nil {
		return err
	}

	if err := c.Mail(from); err != nil { // from
		// settings.WriteDebugLog(err.Error())
		return err
	}
	if err := c.Rcpt(to); err != nil { // to
		// settings.WriteDebugLog(err.Error())
		return err
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		// settings.WriteDebugLog(err.Error())
		return err
	}
	defer wc.Close()

	buf := bytes.NewBufferString("To:" + to)
	buf.WriteString("\r\n")
	buf.WriteString("Subject:" + subject)
	buf.WriteString("\r\n")
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	buf.WriteString(mime + "\n" + templateData.Body)
	if _, err = buf.WriteTo(wc); err != nil {
		// settings.WriteDebugLog(err.Error())
		return err
	}

	// Send the QUIT command and close the connection.
	c.Quit()
	// settings.WriteDebugLog("Email sent successfully")
	return nil
}

func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

package mailer

type Recipients struct {
	FirstName *string
	LastName  *string
	Email     string
}

type Sender struct {
	Email string
	Pass  string
	Name  string
}

type EmailData struct {
	Subject        string
	BodyTemplate   string
	FileAttachment []string
}

type EmailStructure struct {
	Recipients []Recipients
	EmailData  EmailData
}

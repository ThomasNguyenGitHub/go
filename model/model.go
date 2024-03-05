package model

import (
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

const (
	TitlePrefix             = "t"
	MessagePrefix           = "m"
	NotificationDataTitle   = "title"
	NotificationDataBody    = "body"
	NotificationDataMessage = "message"
	NotificationDataType    = "type"
	NotificationDataContent = "content"
)

type PdfOption func(*wkhtmltopdf.PDFGenerator, *wkhtmltopdf.PageReader)

type BaseResponse struct {
	Code    string      `json:"code"`
	Message BaseMessage `json:"message"`
	Data    interface{} `json:"data"`
}

type BaseMessage struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type BaseESBResponse struct {
	Result struct {
		StatusCode string      `json:"STATUSCODE"`
		Data       interface{} `json:"DATA"`
	} `json:"Result"`
}

type TemplateInfo struct {
	Language       string
	TemplateKey    string
	TemplateModule string
	Data           interface{}
}

type BasicSMSInfo struct {
	Otp         string
	Language    string
	Mobile      string
	Action      string
	TemplateKey string
}

type OtpData struct {
	Otp string
}

type SMSInfo struct {
	TemplateInfo
	Mobile string
}

type EmailInfo struct {
	TemplateInfo
	CC             string
	SubjectKey     string
	Recipients     string
	Files          []FileInfo
	PdfAttachments []PdfAttachment
}

type PdfAttachment struct {
	WriteFile   bool
	Key         string
	Password    string
	LocalPath   string
	HeaderURL   string // URL to html file
	FooterURL   string // URL to html file
	FileNameKey string
	Data        interface{}
	PdfOptions  []PdfOption
}

type FileInfo struct {
	FileName string
	FilePath string
}

type PushNotificationInfo struct {
	TemplateInfo
	ClientID  string
	TitleKey  string
	ProjectID string
	Data      map[string]string
}

type PushNotificationUrlTogo struct {
	UrlTogo string                         `json:"url_togo"`
	Param   []PushNotificationUrlTogoParam `json:"param"`
}

type PushNotificationUrlTogoParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

package model

const (
	NtfTemplateStatusActive = "A"
)

// NotificationTemplate a model represents for notification_template table
type NotificationTemplate struct {
	ID             string `gorm:"column:notification_template_id;primary_key"`
	TemplateModule string `gorm:"column:template_module"`
	TemplateKey    string `gorm:"column:template_key"`
	TemplateValue  string `gorm:"column:template_value"`
	Language       string `gorm:"column:language"`
	Description    string `gorm:"column:description"`
	Status         string `gorm:"column:status"`
	CreatedBy      string `gorm:"column:create_by"`
	CreatedAt      string `gorm:"column:create_at"`
	UpdatedBy      string `gorm:"column:update_by"`
	UpdatedAt      string `gorm:"column:update_at"`
}

func (NotificationTemplate) TableName() string {
	return "notification_template"
}

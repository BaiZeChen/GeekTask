package domain

type SMSCallBackArgs struct {
	Biz     string   `json:"biz"`
	Args    []string `json:"args"`
	Numbers []string `json:"numbers"`
}

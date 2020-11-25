package function

//SlackResponse - response from Slack
type SlackResponse struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

//Attachment - part of SlackResponse
type Attachment struct {
	Text     string `json:"text"`
	ImageURL string `json:"image_url"`
}

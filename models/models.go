package models

type GetVideoRequest struct {
	ChannelName string `json:"channelName"`
}

type ImprovedTitle struct {
	NewTitle  string `json:"newTitle"`
	PrevTitle string `json:"prevTitle"`
	Reason    string `json:"reason"`
}

type Response struct {
	LlmResponse []*ImprovedTitle `json:"llmResponse"`
}

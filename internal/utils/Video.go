package utils

import "encoding/json"

type PartialVideo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type PartialCreateVideo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	ChannelId   int64 `json:"channelId"`
}

func (p *PartialVideo) Unmarshal(body []byte) error {
	return json.Unmarshal(body, &p)
}

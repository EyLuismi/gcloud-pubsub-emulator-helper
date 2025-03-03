package pubsub

type Project struct {
	Name   string  `json:"name"`
	Topics []Topic `json:"topics"`
}

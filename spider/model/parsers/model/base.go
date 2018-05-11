package model

type Model struct {
	AcceptedRegUrls []string
	WantedRegUrls []string
	Type string
	Rule map[string]string
	AddQueue []string
}
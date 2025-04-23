package combolist

type SocketFile struct {
	Agent    string `json:"agent"`
	AgentKey string `json:"agent_key"`
	Type     string `json:"type"`
	Raw      struct {
		Files `json:"file"`
	} `json:"raw"`
}

type SocketEntrie struct {
	Agent string `json:"agent"`
	Type  string `json:"type"`
	Hash  string `json:"hash"`
	Raw   struct {
		Users `json:"user_entrie"`
		Urls  `json:"url_entrie"`
	} `json:"raw"`
}

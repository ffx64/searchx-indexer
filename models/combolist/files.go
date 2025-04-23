package combolist

type Files struct {
	ID                    int64  `json:"id"`
	AgentKey              string `json:"agent_key"`
	Name                  string `json:"name"`
	Size                  int64  `json:"size"`
	Hash                  string `json:"hash"`
	Status                int    `json:"status"`
	Source                string `json:"source"`
	Type                  string `json:"type"`
	Description           string `json:"description"`
	ProcessedEntriesCount int64  `json:"processed_entries_count"`
	CreatedAt             string `json:"created_at"`
	ProcessedAt           string `json:"processed_at"`
}

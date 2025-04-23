package combolist

type Urls struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	URL       string `json:"url"`
	FileLine  int64  `json:"file_line"`
	CreatedAt string `json:"created_at"`
}

package combolist

type Users struct {
	ID        int64  `json:"id"`
	FileID    int64  `json:"file_id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

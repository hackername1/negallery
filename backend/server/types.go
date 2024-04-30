package server

type Response struct {
	Status   string `json:"status"`
	Error    string `json:"error"`
	Response any    `json:"content"`
}

type ImageData struct {
	Id          int      `json:"id"`
	Title       string   `json:"title"`
	Country     string   `json:"country"`
	Date        string   `json:"date"`
	Resolutions []string `json:"resolutions"`
}

package data

type Item struct {
	Id          int     `json:"id"`
	Type        string  `json:"type"`
	Deleted     *bool   `json:"deleted"`
	By          *string `json:"by"`
	Time        *int    `json:"time"`
	Text        *string `json:"text"`
	Dead        *bool   `json:"dead"`
	Parent      *int    `json:"parent"`
	Poll        *int    `json:"poll"`
	Kids        *[]int  `json:"kids"`
	Url         *string `json:"url"`
	Score       *int    `json:"score"`
	Title       *string `json:"title"`
	Parts       *[]int  `json:"parts"`
	Descendants *int    `json:"descendants"`
}

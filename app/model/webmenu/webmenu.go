package webmenu

type WebMenu struct {
	Id       interface{}    `json:"id"`
	Icon     string         `json:"icon"`
	Title    string         `json:"title"`
	Type     int            `json:"type"`
	Href     string         `json:"href"`
	Children []MenuChildren `json:"children,omitempty"`
	OpenType string         `json:"openType,omitempty"`
}

type MenuChildren struct {
	Id       interface{}    `json:"id"`
	Title    string         `json:"title"`
	Icon     string         `json:"icon"`
	Type     int            `json:"type"`
	OpenType string         `json:"openType"`
	Href     string         `json:"href"`
	Children []MenuChildren `json:"children,omitempty"`
}

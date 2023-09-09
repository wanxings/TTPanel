package response

type UserRoute struct {
	Routes []interface{} `json:"routes"`
	Home   string        `json:"home"`
}
type Route1 struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Meta      *Meta  `json:"meta"`
}
type Route2 struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Component string    `json:"component"`
	Meta      *Meta     `json:"meta"`
	Children  []*Route1 `json:"children"`
}
type Meta struct {
	Title        string `json:"title"`
	Icon         string `json:"icon"`
	RequiresAuth bool   `json:"requiresAuth"`
	SingleLayout string `json:"singleLayout"`
	Permissions  []int  `json:"permissions"`
	Order        int    `json:"order"`
	KeepAlive    bool   `json:"keepAlive"`
	MultiTab     bool   `json:"multiTab"`
	Hide         bool   `json:"hide"`
	ActiveMenu   string `json:"activeMenu"`
}

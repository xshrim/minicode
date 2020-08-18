package auths

type OAuth struct {
	Url          string `json:"url"`
	Client       string `json:"client"`
	Redirect     string `json:"redirect"`
	Scope        string `json:"scope"`
	State        string `json:"state"`
	Code         string `json:"code"`
	Secret       string `json:"secret"`
	TokenType    string `json:"tokenType"`
	TokenCreate  int64  `json:"tokenCreate"`
	TokenExpire  int64  `json:"tokenExpire"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type User struct {
	Name       string `json:"name"`
	GitlabId   string `json:"gitlabId"`
	GitlabName string `json:"gitlabName"`
	MstodoId   string `json:"mstodoId"`
	MstodoName string `json:"mstodoName"`
}

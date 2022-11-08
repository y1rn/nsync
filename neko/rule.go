package neko

type Status int

const (
	Status_Active = 1
)

type RuleResp struct {
	Status Status `json:"status"`
	Data   []Data `json:"data"`
}

type Data struct {
	Rid     string  `json:"rid"`
	Uid     string  `json:"uid"`
	Sid     string  `json:"sid"`
	Host    string  `json:"host"`
	Port    int     `json:"port"`
	Remote  string  `json:"remote"`
	Rport   int     `json:"rport"`
	Type    string  `json:"type"`
	Status  int     `json:"status"`
	Name    string  `json:"name"`
	Traffic int     `json:"traffic"`
	Date    int     `json:"date"`
	Setting Setting `json:"setting"`
}

type Setting struct {
	ProxyProtocol   int64  `json:"proxyProtocol"`
	LoadbalanceMode string `json:"loadbalanceMode"`
}

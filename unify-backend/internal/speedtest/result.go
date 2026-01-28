package speedtest

type SpeedtestResult struct {
	Timestamp string `json:"timestamp"`

	Ping struct {
		Latency float64 `json:"latency"`
	} `json:"ping"`

	Download struct {
		Bandwidth int64 `json:"bandwidth"`
	} `json:"download"`

	Upload struct {
		Bandwidth int64 `json:"bandwidth"`
	} `json:"upload"`

	ISP string `json:"isp"`

	Interface struct {
		InternalIP string `json:"internalIp"`
		Name       string `json:"name"`
		MACAddr    string `json:"macAddr"`
		ExternalIP string `json:"externalIp"`
	} `json:"interface"`

	Server struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Location string `json:"location"`
		Country  string `json:"country"`
	} `json:"server"`

	Result struct {
		URL string `json:"url"`
	} `json:"result"`
}

package response

type TokenLaunchpadInfo struct {
	Address string `json:"address"`
	Launchpad        string `json:"launchpad"`
	LaunchpadStatus  int    `json:"launchpad_status"`
	LaunchpadProgress string `json:"launchpad_progress"`
	Description       string `json:"description"`
}

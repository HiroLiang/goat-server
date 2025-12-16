package agent

type AgentInfoCondition struct {
}

type AgentInfo struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Status   string `json:"status"`
}

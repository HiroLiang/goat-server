package agent

import (
	"net/http"

	"github.com/HiroLiang/goat-server/internal/application/agent"
	"github.com/HiroLiang/goat-server/internal/interface/http/adapter"
	"github.com/gin-gonic/gin"
)

type AgentHandler struct {
	agentUseCase *agent.UseCase
}

func NewAgentHandler(agentUseCase *agent.UseCase) *AgentHandler {
	return &AgentHandler{agentUseCase}
}

// RegisterAgentRoutes registers user-related API routes
func (h *AgentHandler) RegisterAgentRoutes(r *gin.RouterGroup) {
	r.GET("/available", h.getAvailableAgents)
}

// @Summary Available agents info
// @Description query all available agents
// @Tags Agent
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} AgentInfo
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/agent/available [get]
func (h *AgentHandler) getAvailableAgents(c *gin.Context) {
	input := adapter.BuildInput[agent.QueryAvailableAgentsInput](c, agent.QueryAvailableAgentsInput{})

	outputs, err := h.agentUseCase.FindAvailableAgents(c.Request.Context(), input)
	if err != nil {
		HandleError(c, err)
		return
	}

	agents := make([]AgentInfo, len(outputs))
	for i, output := range outputs {
		agents[i] = AgentInfo{
			Name:     output.Name,
			Provider: output.Provider,
			Status:   output.Status.Desc(),
		}
	}

	c.JSON(http.StatusOK, agents)
}

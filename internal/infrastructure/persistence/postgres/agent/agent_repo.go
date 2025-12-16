package agent

import (
	"context"
	"fmt"

	"github.com/HiroLiang/goat-server/internal/domain/agent"
	"github.com/HiroLiang/goat-server/internal/infrastructure/persistence/database"
	"github.com/HiroLiang/goat-server/internal/infrastructure/persistence/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var agentTable = postgres.Table{
	Name: "public.agents",
	Columns: []string{
		"id",
		"name",
		"type",
		"status",
		"engine",
		"created_at",
		"updated_at",
	},
}

type AgentRepository struct {
	db *sqlx.DB
}

var _ agent.Repository = (*AgentRepository)(nil)

func NewAgentRepository(dbName database.DBName) *AgentRepository {
	return &AgentRepository{db: database.GetDB(dbName)}
}

// FindAll returns all agents.
func (r AgentRepository) FindAll(ctx context.Context) ([]*agent.Agent, error) {
	return r.find(ctx, nil)
}

// FindAllByStatus returns all agents by status.
func (r AgentRepository) FindAllByStatus(
	ctx context.Context,
	status agent.Status,
) ([]*agent.Agent, error) {
	return r.find(ctx, squirrel.Eq{"status": status})
}

// Create inserts a new agent.
func (r AgentRepository) Create(ctx context.Context, agent *agent.Agent) error {
	record := toRecord(agent)

	query, args, err := agentTable.Insert().
		Columns("name", "type", "status", "engine", "create_by", "update_by").
		Values(record.Name, record.Type, record.Status, record.Engine, record.CreatedBy, record.UpdatedBy).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert agent: %w", err)
	}

	if err := postgres.Exec(ctx, r.db, query, args...); err != nil {
		return fmt.Errorf("insert agent: %w", err)
	}

	return nil
}

// find returns agents by condition.
func (r AgentRepository) find(
	ctx context.Context,
	cond squirrel.Sqlizer,
) ([]*agent.Agent, error) {

	builder := agentTable.Select(agentTable.Columns...)
	if cond != nil {
		builder = builder.Where(cond)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build sql: %w", err)
	}

	records, err := postgres.ScanAll[AgentRecord](ctx, r.db, query, args...)
	if err != nil {
		return nil, fmt.Errorf("scan agents: %w", err)
	}

	agents := make([]*agent.Agent, 0, len(records))
	for _, rec := range records {
		domain, err := toDomain(&rec)
		if err != nil {
			return nil, fmt.Errorf("convert agent: %w", err)
		}
		agents = append(agents, domain)
	}

	return agents, nil
}

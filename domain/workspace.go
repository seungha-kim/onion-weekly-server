package domain

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/onion-studio/onion-weekly/dto"

	"github.com/onion-studio/onion-weekly/config"
)

type WorkspaceService struct {
	appConf config.AppConf
}

func NewWorkspaceService(appConf config.AppConf) *WorkspaceService {
	return &WorkspaceService{appConf: appConf}
}

func (s *WorkspaceService) createWorkspace(
	tx pgx.Tx,
	actor dto.User,
	input dto.CreateWorkspaceInput,
) (w dto.Workspace, err error) {
	q := `
insert into workspaces (name, created_by)
values ($1, $2)
returning id, name, created_by, created_at;
`
	ctx := context.Background()
	err = tx.
		QueryRow(ctx, q, input.Name, actor.Id).
		Scan(&w.ID, &w.Name, &w.CreatedBy, &w.CreatedAt)
	if err != nil {
		return
	}

	q = `
insert into workspace_members (user_id, workspace_id)
values ($1, $2);
`
	_, err = tx.
		Exec(ctx, q, actor.Id, w.ID)

	return
}

func (s *WorkspaceService) findWorkspacesByMembership(
	tx pgx.Tx,
	user dto.User,
) (workspaces []dto.Workspace, err error) {
	const q = `
select w.id, w.name, w.created_by, w.created_at from workspaces w
join workspace_members wm on w.id = wm.workspace_id
where wm.user_id = $1;
`
	ctx := context.Background()
	rows, err := tx.Query(ctx, q, user.Id)
	if err != nil {
		return
	}

	for rows.Next() {
		w := dto.Workspace{}
		err = rows.Scan(&w.ID, &w.Name, &w.CreatedBy, &w.CreatedAt)
		if err != nil {
			return
		}
		workspaces = append(workspaces, w)
	}

	return
}

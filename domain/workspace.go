package domain

import (
	"context"
	"errors"

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

func (srv *WorkspaceService) createWorkspace(
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
		Scan(&w.Id, &w.Name, &w.CreatedBy, &w.CreatedAt)
	if err != nil {
		return
	}

	q = `
insert into workspace_members (user_id, workspace_id)
values ($1, $2);
`
	_, err = tx.
		Exec(ctx, q, actor.Id, w.Id)

	return
}

func (srv *WorkspaceService) findWorkspacesByMembership(
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
		err = rows.Scan(&w.Id, &w.Name, &w.CreatedBy, &w.CreatedAt)
		if err != nil {
			return
		}
		workspaces = append(workspaces, w)
	}

	return
}

func (srv *WorkspaceService) checkMembership(
	tx pgx.Tx,
	workspace dto.Workspace,
	user dto.User,
) (bool, error) {
	const q = `
select count(*) from workspace_members wm
where wm.workspace_id = $1 and wm.user_id = $2;
`
	var count int64
	row := tx.QueryRow(context.Background(), q, workspace.Id, user.Id)
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (srv *WorkspaceService) invite(
	tx pgx.Tx,
	workspace dto.Workspace,
	actor dto.User,
	invitee dto.User,
) (err error) {
	isMember, err := srv.checkMembership(tx, workspace, actor)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("actor is not a member")
	}
	const q = `
insert into workspace_members (user_id, workspace_id)
values ($1, $2);
`
	if _, err := tx.Exec(context.Background(), q, invitee.Id, workspace.Id); err != nil {
		return err
	}
	return nil
}

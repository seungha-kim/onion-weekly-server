package domain

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/onion-studio/onion-weekly/config"
	"github.com/onion-studio/onion-weekly/dto"
)

type RecurringService struct {
	appConf          config.AppConf
	workspaceService *WorkspaceService
}

func NewRecurringService(appConf config.AppConf, workspaceService *WorkspaceService) *RecurringService {
	return &RecurringService{appConf: appConf, workspaceService: workspaceService}
}

func (srv *RecurringService) CreateRecurring(
	tx pgx.Tx,
	actor dto.User,
	workspace dto.Workspace,
	input dto.CreateRecurringInput,
) (recurring dto.Recurring, err error) {
	ctx := context.Background()

	if err = srv.checkCreatePermission(tx, actor, workspace); err != nil {
		return
	}

	const q = `
insert into recurrings (workspace_id, title, "interval")
values ($1, $2, $3)
returning id, workspace_id, title, "interval", created_at;
`
	err = tx.
		QueryRow(ctx, q, workspace.Id, input.Title, input.Interval).
		Scan(&recurring.Id, &recurring.WorkspaceId, &recurring.Title, &recurring.Interval, &recurring.CreatedAt)
	return
}

func (srv *RecurringService) UpdateRecurring(tx pgx.Tx, user dto.User, recurring dto.Recurring, input dto.UpdateRecurringInput) (result dto.Recurring, err error) {
	ctx := context.Background()
	if err := srv.checkUpdatePermission(tx, user, recurring); err != nil {
		return dto.Recurring{}, err
	}

	q := `
update recurrings
set title = $1, interval = $2
where id = $3
returning id, workspace_id, title, "interval", created_at;
`
	err = tx.
		QueryRow(ctx, q, input.Title, input.Interval, recurring.Id).
		Scan(&result.Id, &result.WorkspaceId, &result.Title, &result.Interval, &result.CreatedAt)
	return
}

func (srv *RecurringService) CreateRecurringRecord(tx pgx.Tx, user dto.User, recurring dto.Recurring, input dto.CreateRecurringRecordInput) (result dto.RecurringRecord, err error) {
	ctx := context.Background()
	if err = srv.checkUpdatePermission(tx, user, recurring); err != nil {
		return
	}

	q := `
insert into recurring_records (description, actor_id, interval_snapshot, recurring_id)
values ($1, $2, $3, $4)
returning id, description, actor_id, interval_snapshot, recurring_id;
`
	err = tx.
		QueryRow(ctx, q, input.Description, user.Id, recurring.Interval, recurring.Id).
		Scan(&result.Id, &result.Description, &result.ActorId, &result.IntervalSnapshot, &result.RecurringId)
	return
}

func (srv *RecurringService) FindRecurringRecordsByWorkspace(tx pgx.Tx, user dto.User, workspace dto.Workspace) (records []dto.RecurringRecord, err error) {
	ctx := context.Background()
	if err = srv.workspaceService.checkReadPermission(tx, user, workspace); err != nil {
		return
	}
	q := `
select rr.id, rr.description, rr.actor_id, rr.interval_snapshot, rr.created_at, rr.recurring_id from recurring_records rr
join recurrings r on rr.recurring_id = r.id
where r.workspace_id = $1;
`
	rows, err := tx.Query(ctx, q, workspace.Id)
	if err != nil {
		return
	}

	for rows.Next() {
		rr := dto.RecurringRecord{}
		err = rows.Scan(&rr.Id, &rr.Description, &rr.ActorId, &rr.IntervalSnapshot, &rr.CreatedAt, &rr.RecurringId)
		if err != nil {
			return
		}
		records = append(records, rr)
	}
	return
}

func (srv *RecurringService) FindRecurringsByWorkspace(tx pgx.Tx, user dto.User, workspace dto.Workspace) (recurrings []dto.Recurring, err error) {
	ctx := context.Background()
	if err = srv.workspaceService.checkReadPermission(tx, user, workspace); err != nil {
		return
	}

	q := `
with last_records as (
    select distinct on (rr.recurring_id) rr.* from recurring_records rr
    join recurrings r on rr.recurring_id = r.id
    where r.workspace_id = $1
    order by rr.recurring_id, rr.created_at desc
)
select r.id, r.workspace_id, r.title, r.interval, r.created_at
from recurrings r
left join last_records lr on lr.recurring_id = r.id
where r.workspace_id = $1
order by (extract(epoch from now()) - coalesce(extract(epoch from lr.created_at), float8 '-infinity')) / r.interval desc
`
	rows, err := tx.Query(ctx, q, workspace.Id)
	if err != nil {
		return
	}

	for rows.Next() {
		r := dto.Recurring{}
		err = rows.Scan(&r.Id, &r.WorkspaceId, &r.Title, &r.Interval, &r.CreatedAt)
		if err != nil {
			return
		}
		recurrings = append(recurrings, r)
	}
	return
}

func (srv *RecurringService) GetRecurringById(tx pgx.Tx, user dto.User, id dto.UUID) (recurring dto.Recurring, err error) {
	ctx := context.Background()
	const q = `
select r.id, r.workspace_id, r.title, r.interval, r.created_at from recurrings r
where r.id = $1;
`
	if err = tx.QueryRow(ctx, q, id).Scan(&recurring.Id, &recurring.WorkspaceId, &recurring.Title, &recurring.Interval, &recurring.CreatedAt); err != nil {
		return
	}
	if err = srv.checkReadPermission(tx, user, recurring); err != nil {
		return
	}
	return
}

func (srv *RecurringService) checkCreatePermission(tx pgx.Tx, user dto.User, workspace dto.Workspace) error {
	return srv.workspaceService.checkReadPermission(tx, user, workspace)
}

func (srv *RecurringService) checkUpdatePermission(tx pgx.Tx, user dto.User, recurring dto.Recurring) error {
	workspace, err := srv.workspaceService.GetWorkspaceById(tx, recurring.WorkspaceId)
	if err != nil {
		return errors.New("forbidden")
	}
	isMember, err := srv.workspaceService.checkMembership(tx, workspace, user)
	if !isMember || err != nil {
		return errors.New("forbidden")
	}
	return nil
}

func (srv *RecurringService) checkReadPermission(tx pgx.Tx, user dto.User, recurring dto.Recurring) error {
	workspace, err := srv.workspaceService.GetWorkspaceById(tx, recurring.WorkspaceId)
	if err != nil {
		return errors.New("forbidden")
	}
	isMember, err := srv.workspaceService.checkMembership(tx, workspace, user)
	if !isMember || err != nil {
		return errors.New("forbidden")
	}
	return nil
}

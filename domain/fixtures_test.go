package domain

import (
	"context"

	"github.com/onion-studio/onion-weekly/dto"

	"github.com/jackc/pgx/v4"
)

func queryRowAndScan(tx pgx.Tx, q string, dest ...interface{}) {
	ctx := context.Background()
	if err := tx.QueryRow(ctx, q).Scan(dest...); err != nil {
		panic(err)
	}
}

func fixtureUser1(tx pgx.Tx) (u dto.User, up dto.UserProfile, ec dto.EmailCredential) {
	var q string
	q = `
insert into users (id)
values ('F380634D-CEB3-421B-9522-8530EBCE05CB')
returning id, created_at;
`
	queryRowAndScan(tx, q, &u.Id, &u.CreatedAt)

	q = `
insert into user_profiles (user_id, full_name)
values ('F380634D-CEB3-421B-9522-8530EBCE05CB', '김승하')
returning user_id, full_name;
`
	queryRowAndScan(tx, q, &up.UserId, &up.FullName)

	// password=onion
	q = `
insert into email_credentials (user_id, email, hashed_password)
values ('F380634D-CEB3-421B-9522-8530EBCE05CB', 'seungha.me@gmail.com', '$2a$14$.5CKkT8/OS5IEJs15ecbxehkLfM5BrBvAU9.d2JBCEmYQE8QwIK1q')
returning user_id, email, hashed_password, created_at;
`
	queryRowAndScan(tx, q, &ec.UserId, &ec.Email, &ec.HashedPassword, &ec.CreatedAt)

	return
}

func fixtureUser2(tx pgx.Tx) (u dto.User, up dto.UserProfile, ec dto.EmailCredential) {
	var q string
	q = `
insert into users (id)
values ('7F044C9A-557D-43EC-91E9-B4D44388CE9B')
returning id, created_at;
`
	queryRowAndScan(tx, q, &u.Id, &u.CreatedAt)

	q = `
insert into user_profiles (user_id, full_name)
values ('7F044C9A-557D-43EC-91E9-B4D44388CE9B', 'Jeesoo Hong')
returning user_id, full_name;
`
	queryRowAndScan(tx, q, &up.UserId, &up.FullName)

	// password=weekly
	q = `
insert into email_credentials (user_id, email, hashed_password)
values ('7F044C9A-557D-43EC-91E9-B4D44388CE9B', 'jeesoo.hong@gmail.com', '$2a$14$QMqZ.1qKqayUR34FaTlX..ytZwUM2U7AnMZ5DEOej6ggP4xfjSEPS')
returning user_id, email, hashed_password, created_at;
`
	queryRowAndScan(tx, q, &ec.UserId, &ec.Email, &ec.HashedPassword, &ec.CreatedAt)

	return
}

func fixtureWorkspace1(tx pgx.Tx, user dto.User) (w dto.Workspace) {
	ctx := context.Background()
	var q string
	q = `
insert into workspaces (id, name, created_by)
values ('EC2A7558-82E9-411A-85FF-55573486119C', 'Onion Studio', $1)
returning id, name, created_by, created_at;
`
	if err := tx.QueryRow(ctx, q, user.Id).Scan(&w.Id, &w.Name, &w.CreatedBy, &w.CreatedAt); err != nil {
		panic(err)
	}

	q = `
insert into workspace_members (user_id, workspace_id)
values ($1, $2);
`
	if _, err := tx.Exec(ctx, q, user.Id, w.Id); err != nil {
		panic(err)
	}
	return
}

func fixtureRecurring1(tx pgx.Tx, workspace dto.Workspace) (r dto.Recurring) {
	ctx := context.Background()
	q := `
insert into recurrings (id, workspace_id, title, interval)
values ('0B8B8599-6E14-4A66-9E5D-3620D3551B62', $1, '우유 사기', 4)
returning id, workspace_id, title, "interval", created_at;
`
	err := tx.
		QueryRow(ctx, q, workspace.Id).
		Scan(&r.Id, &r.WorkspaceId, &r.Title, &r.Interval, &r.CreatedAt)
	if err != nil {
		panic(err)
	}
	return
}

func fixtureRecurring2(tx pgx.Tx, workspace dto.Workspace) (r dto.Recurring) {
	ctx := context.Background()
	q := `
insert into recurrings (id, workspace_id, title, interval)
values ('F062DCD1-5967-4448-AD8F-07AD0EFEB717', $1, '화분에 물 주기', 4)
returning id, workspace_id, title, "interval", created_at;
`
	err := tx.
		QueryRow(ctx, q, workspace.Id).
		Scan(&r.Id, &r.WorkspaceId, &r.Title, &r.Interval, &r.CreatedAt)
	if err != nil {
		panic(err)
	}
	return
}

func fixtureRecurringRecord(tx pgx.Tx, user dto.User, recurring dto.Recurring, createdAt dto.Timestamptz) (rr dto.RecurringRecord) {
	ctx := context.Background()
	q := `
insert into recurring_records (description, actor_id, interval_snapshot, created_at, recurring_id)
values ('test', $1, $2, $3, $4)
returning id, description, actor_id, interval_snapshot, created_at, recurring_id;
`
	err := tx.
		QueryRow(ctx, q, user.Id, recurring.Interval, createdAt, recurring.Id).
		Scan(&rr.Id, &rr.Description, &rr.ActorId, &rr.IntervalSnapshot, &rr.CreatedAt, &rr.RecurringId)
	if err != nil {
		panic(err)
	}
	return
}

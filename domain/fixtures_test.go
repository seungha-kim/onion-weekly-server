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

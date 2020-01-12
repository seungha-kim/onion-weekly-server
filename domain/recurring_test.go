package domain

import (
	"context"
	"testing"
	"time"

	"github.com/onion-studio/onion-weekly/dto"
	"github.com/stretchr/testify/require"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/onion-studio/onion-weekly/db"

	"github.com/stretchr/testify/suite"

	"github.com/onion-studio/onion-weekly/config"
	"go.uber.org/fx"
)

type RecurringTestSuite struct {
	suite.Suite
	appConf          config.AppConf
	pgxPool          *pgxpool.Pool
	recurringService *RecurringService
}

func (s *RecurringTestSuite) SetupTest() {
	app := fx.New(
		fx.Provide(
			config.NewTestAppConf,
			NewWorkspaceService,
			NewRecurringService,
			db.NewPgxPool,
		),
		fx.Populate(&s.appConf, &s.recurringService, &s.pgxPool),
		fx.NopLogger,
	)
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
}

func (s *RecurringTestSuite) TestCreateRecurring() {
	t := s.T()
	t.Run("멤버는 워크스페이스에 챙길 것을 추가할 수 있다", func(t *testing.T) {
		db.RollbackForTest(s.pgxPool, func(tx pgx.Tx) {
			user1, _, _ := fixtureUser1(tx)
			workspace1 := fixtureWorkspace1(tx, user1)
			_, err := s.recurringService.CreateRecurring(
				tx,
				user1,
				workspace1,
				dto.CreateRecurringInput{
					Title:    "Test Recurring",
					Interval: 1,
				},
			)

			require.NoError(t, err)
		})
	})

	t.Run("멤버가 아닌 사용자는 워크스페이스에 챙길 것을 추가할 수 없다", func(t *testing.T) {
		db.RollbackForTest(s.pgxPool, func(tx pgx.Tx) {
			user1, _, _ := fixtureUser1(tx)
			user2, _, _ := fixtureUser2(tx)
			workspace1 := fixtureWorkspace1(tx, user1)
			_, err := s.recurringService.CreateRecurring(
				tx,
				user2,
				workspace1,
				dto.CreateRecurringInput{
					Title:    "Test Recurring",
					Interval: 1,
				},
			)

			require.Error(t, err)
		})
	})
}

func (s *RecurringTestSuite) TestUpdateRecurring() {
	t := s.T()
	t.Run("멤버는 챙길 것의 내용을 변경할 수 있다", func(t *testing.T) {
		db.RollbackForTest(s.pgxPool, func(tx pgx.Tx) {
			user1, _, _ := fixtureUser1(tx)
			workspace1 := fixtureWorkspace1(tx, user1)
			recurring1 := fixtureRecurring1(tx, workspace1)

			updated, err := s.recurringService.UpdateRecurring(
				tx,
				user1,
				recurring1,
				dto.UpdateRecurringInput{
					Title:    "과일 사기",
					Interval: 1,
				},
			)

			require.NoError(t, err)
			require.Equal(t, "우유 사기", recurring1.Title)
			require.Equal(t, "과일 사기", updated.Title)
		})
	})

	t.Run("멤버가 아닌 사용자는 챙길 것의 내용을 변경할 수 없다", func(t *testing.T) {
		db.RollbackForTest(s.pgxPool, func(tx pgx.Tx) {
			user1, _, _ := fixtureUser1(tx)
			user2, _, _ := fixtureUser2(tx)
			workspace1 := fixtureWorkspace1(tx, user1)
			recurring1 := fixtureRecurring1(tx, workspace1)

			_, err := s.recurringService.UpdateRecurring(
				tx,
				user2,
				recurring1,
				dto.UpdateRecurringInput{
					Title:    "과일 사기",
					Interval: 1,
				},
			)

			require.Error(t, err)
		})
	})
}

func (s *RecurringTestSuite) TestCreateRecurringRecord() {
	t := s.T()

	t.Run("멤버는 챙긴 기록을 추가할 수 있다.", func(t *testing.T) {
		db.RollbackForTest(s.pgxPool, func(tx pgx.Tx) {
			user1, _, _ := fixtureUser1(tx)
			workspace1 := fixtureWorkspace1(tx, user1)
			recurring1 := fixtureRecurring1(tx, workspace1)

			rr, err := s.recurringService.CreateRecurringRecord(
				tx,
				user1,
				recurring1,
				dto.CreateRecurringRecordInput{Description: "하하"},
			)

			require.NoError(t, err)
			require.Equal(t, rr.Description, "하하")
		})
	})

	t.Run("멤버가 아닌 사용자는 챙긴 기록을 추가할 수 없다.", func(t *testing.T) {
		db.RollbackForTest(s.pgxPool, func(tx pgx.Tx) {
			user1, _, _ := fixtureUser1(tx)
			user2, _, _ := fixtureUser2(tx)
			workspace1 := fixtureWorkspace1(tx, user1)
			recurring1 := fixtureRecurring1(tx, workspace1)

			_, err := s.recurringService.CreateRecurringRecord(
				tx,
				user2,
				recurring1,
				dto.CreateRecurringRecordInput{Description: "하하"},
			)

			require.Error(t, err)
		})
	})
}

func (s *RecurringTestSuite) TestFindRecurringRecordsByWorkspace() {
	t := s.T()
	t.Run("멤버는 워크스페이스의 챙긴 기록을 시간 순으로 볼 수 있다.", func(t *testing.T) {
		db.RollbackForTest(s.pgxPool, func(tx pgx.Tx) {
			user1, _, _ := fixtureUser1(tx)
			workspace1 := fixtureWorkspace1(tx, user1)

			r1, err := s.recurringService.CreateRecurring(
				tx,
				user1,
				workspace1,
				dto.CreateRecurringInput{
					Title:    "1",
					Interval: 1,
				},
			)
			require.NoError(t, err)
			r2, err := s.recurringService.CreateRecurring(
				tx,
				user1,
				workspace1,
				dto.CreateRecurringInput{
					Title:    "2",
					Interval: 2,
				},
			)
			require.NoError(t, err)

			_, err = s.recurringService.CreateRecurringRecord(
				tx,
				user1,
				r2, // for r2
				dto.CreateRecurringRecordInput{Description: "1"},
			)

			_, err = s.recurringService.CreateRecurringRecord(
				tx,
				user1,
				r1, // for r1
				dto.CreateRecurringRecordInput{Description: "2"},
			)

			records, err := s.recurringService.FindRecurringRecordsByWorkspace(
				tx,
				user1,
				workspace1,
			)

			require.NoError(t, err)
			require.Len(t, records, 2)
			require.Equal(t, records[0].Description, "1")
			require.Equal(t, records[0].RecurringId, r2.Id)
			require.Equal(t, records[0].IntervalSnapshot, r2.Interval)
			require.Equal(t, records[0].ActorId, user1.Id)
			require.Equal(t, records[1].Description, "2")
			require.Equal(t, records[1].RecurringId, r1.Id)
			require.Equal(t, records[1].IntervalSnapshot, r1.Interval)
			require.Equal(t, records[1].ActorId, user1.Id)
		})
	})
}

func (s *RecurringTestSuite) TestFindRecurringsByWorkspace() {
	t := s.T()
	t.Run("멤버는 워크스페이스의 챙길 것들을 시급도 순으로 볼 수 있다.", func(t *testing.T) {
		db.RollbackForTest(s.pgxPool, func(tx pgx.Tx) {
			user1, _, _ := fixtureUser1(tx)
			workspace1 := fixtureWorkspace1(tx, user1)
			recurring1 := fixtureRecurring1(tx, workspace1)
			recurring2 := fixtureRecurring2(tx, workspace1)

			anYearAgo := dto.Timestamptz{}
			if err := anYearAgo.Set(time.Now().Add(time.Hour * 24 * -365)); err != nil {
				panic(err)
			}
			aMonthAgo := dto.Timestamptz{}
			if err := aMonthAgo.Set(time.Now().Add(time.Hour * 24 * -30)); err != nil {
				panic(err)
			}
			// Given
			_ = fixtureRecurringRecord(tx, user1, recurring1, anYearAgo)

			// When
			rs1, err := s.recurringService.FindRecurringsByWorkspace(tx, user1, workspace1)

			// Then
			require.NoError(t, err)
			require.Equal(t, rs1[0].Id, recurring2.Id) // recurring2를 안 챙겼으니 제일 위에

			// Given
			_ = fixtureRecurringRecord(tx, user1, recurring2, aMonthAgo)

			// When
			rs2, err := s.recurringService.FindRecurringsByWorkspace(tx, user1, workspace1)

			// Then
			require.NoError(t, err)
			require.Equal(t, rs2[0].Id, recurring1.Id)
		})
	})
}

// Next

// 멤버는 챙긴 기록을 추가/수정할 때 누가 한 일인지, 언제 한 것인지 설정할 수 있다.
// 멤버는 챙긴 기록을 수정할 수 있다.
// 멤버는 챙긴 기록을 삭제할 수 있다.

func TestRecurringServiceSuite(t *testing.T) {
	suite.Run(t, new(RecurringTestSuite))
}

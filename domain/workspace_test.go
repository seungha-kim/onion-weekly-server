package domain

import (
	"testing"

	"github.com/onion-studio/onion-weekly/dto"
	"github.com/stretchr/testify/require"

	"github.com/jackc/pgx/v4"
	"github.com/onion-studio/onion-weekly/db"

	"github.com/stretchr/testify/suite"

	"github.com/onion-studio/onion-weekly/config"
	"go.uber.org/fx"
)

type WorkspaceTestSuite struct {
	appConf config.AppConf
	suite.Suite
	workspaceService *WorkspaceService
}

func (s *WorkspaceTestSuite) SetupTest() {
	fx.New(
		fx.Provide(
			config.NewTestAppConf,
			NewWorkspaceService,
		),
		fx.Populate(&s.appConf, &s.workspaceService),
		fx.NopLogger,
	)
}

func (s *WorkspaceTestSuite) TestWorkspaceService() {
	t := s.T()
	t.Run("워크스페이스를 생성할 수 있다", func(t *testing.T) {
		db.RollbackForTest(pgxPool, func(tx pgx.Tx) {
			user, _, _ := fixtureUser1(tx)
			input := dto.CreateWorkspaceInput{Name: "Test Workspace"}
			w, err := s.workspaceService.createWorkspace(tx, user, input)
			require.NoError(t, err)
			require.Equal(t, "Test Workspace", w.Name)

			workspaces, err := s.workspaceService.findWorkspacesByMembership(tx, user)
			require.NoError(t, err)
			require.Len(t, workspaces, 1)
		})
	})

	t.Run("본인의 워크스페이스만 반환되어야 한다", func(t *testing.T) {
		db.RollbackForTest(pgxPool, func(tx pgx.Tx) {
			user1, _, _ := fixtureUser1(tx)
			user2, _, _ := fixtureUser2(tx)
			input := dto.CreateWorkspaceInput{Name: "Test Workspace"}
			w, err := s.workspaceService.createWorkspace(tx, user1, input)
			require.NoError(t, err)
			require.Equal(t, "Test Workspace", w.Name)

			workspaces, err := s.workspaceService.findWorkspacesByMembership(tx, user2)
			require.NoError(t, err)
			require.Len(t, workspaces, 0)
		})
	})

	t.Run("다른 사용자를 워크스페이스의 멤버로 가입시킬 수 있다", func(t *testing.T) {
		db.RollbackForTest(pgxPool, func(tx pgx.Tx) {
			user1, _, _ := fixtureUser1(tx)
			user2, _, _ := fixtureUser2(tx)
			workspace1 := fixtureWorkspace1(tx, user1)

			isMember, err := s.workspaceService.checkMembership(tx, workspace1, user2)
			require.NoError(t, err)
			require.False(t, isMember)

			err = s.workspaceService.invite(tx, workspace1, user1, user2)
			require.NoError(t, err)

			isMember, err = s.workspaceService.checkMembership(tx, workspace1, user2)
			require.NoError(t, err)
			require.True(t, isMember)
		})
	})

	t.Run("이미 멤버인 사용자를 재가입시키려고 시도하면 에러", func(t *testing.T) {
		db.RollbackForTest(pgxPool, func(tx pgx.Tx) {
			user1, _, _ := fixtureUser1(tx)
			workspace1 := fixtureWorkspace1(tx, user1)

			err := s.workspaceService.invite(tx, workspace1, user1, user1)
			require.Error(t, err)
		})
	})

	t.Run("멤버가 아닌 사람이 초대 시도를 하면 에러", func(t *testing.T) {
		db.RollbackForTest(pgxPool, func(tx pgx.Tx) {
			user1, _, _ := fixtureUser1(tx)
			user2, _, _ := fixtureUser2(tx)
			workspace1 := fixtureWorkspace1(tx, user1)

			err := s.workspaceService.invite(tx, workspace1, user2, user2)
			require.Error(t, err)
		})
	})

	// Next

	// 관리자는 특정 멤버의 관리자 권한을 부여/철회할 수 있다
	// 관리자는 최소 1명 이상 존재해야 한다

	// 관리자는 다른 멤버를 워크스페이스에서 추방할 수 있다
	// 관리자는 최소 1명 이상 존재해야 한다
}

func TestWorkspaceSuite(t *testing.T) {
	suite.Run(t, new(WorkspaceTestSuite))
}

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
	t.Run("워크스페이스 생성, 불러오기", func(t *testing.T) {
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

	t.Run("본인의 워크스페이스만 반환되어야 함", func(t *testing.T) {
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
}

func TestWorkspaceSuite(t *testing.T) {
	suite.Run(t, new(WorkspaceTestSuite))
}

package web

import (
	"github.com/jackc/pgx/v4"
	em "github.com/labstack/echo/middleware"
	"github.com/onion-studio/onion-weekly/dto"

	"github.com/onion-studio/onion-weekly/config"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/labstack/echo"
	"github.com/onion-studio/onion-weekly/domain"
	"github.com/onion-studio/onion-weekly/web/middleware"
)

type workspaceHandler struct {
	appConf          config.AppConf
	pgxPool          *pgxpool.Pool
	userService      *domain.UserService
	workspaceService *domain.WorkspaceService
}

func NewWorkspaceHandler(
	appConf config.AppConf,
	pgxPool *pgxpool.Pool,
	workspaceService *domain.WorkspaceService,
) *workspaceHandler {
	return &workspaceHandler{
		appConf:          appConf,
		pgxPool:          pgxPool,
		workspaceService: workspaceService,
	}
}

func (h *workspaceHandler) Register(g *echo.Group) {
	g.Use(middleware.Transaction(h.appConf, h.pgxPool))
	g.Use(em.JWT(h.appConf.Secret), middleware.Actor(h.pgxPool, h.userService))
	g.GET("", h.handleGetWorkspaces)
	g.POST("", h.handlePostWorkspace)
	g.GET("/:workspace-id", h.handleGetWorkspace)
	//g.POST("/register", h.handlePostUser)
	//g.GET("/me", h.handleGetTokenPayload, em.JWT(h.appConf.Secret))
}

func (h *workspaceHandler) handlePostWorkspace(c echo.Context) (err error) {
	tx := c.Get("tx").(pgx.Tx)
	actor := c.Get("actor").(dto.User)

	input := dto.CreateWorkspaceInput{}
	if err = c.Bind(&input); err != nil {
		return err
	}

	workspace, err := h.workspaceService.CreateWorkspace(tx, actor, input)
	if err != nil {
		return err
	}
	return c.JSON(200, workspace)
}

func (h *workspaceHandler) handleGetWorkspaces(c echo.Context) (err error) {
	tx := c.Get("tx").(pgx.Tx)
	actor := c.Get("actor").(dto.User)

	workspaces, err := h.workspaceService.FindWorkspacesByMembership(tx, actor)
	if err != nil {
		return err
	}
	return c.JSON(200, workspaces)
}

func (h *workspaceHandler) handleGetWorkspace(c echo.Context) (err error) {
	tx := c.Get("tx").(pgx.Tx)
	//actor := c.Get("actor").(dto.User)
	id := dto.UUID{}
	if err := id.Set(c.Get("workspace-id")); err != nil {
		return err
	}
	// TODO: permission check for actor
	workspace, err := h.workspaceService.GetWorkspaceById(tx, id)
	if err != nil {
		return err
	}
	return c.JSON(200, workspace)
}

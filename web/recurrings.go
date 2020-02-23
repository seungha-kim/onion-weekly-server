package web

import (
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	em "github.com/labstack/echo/v4/middleware"
	"github.com/onion-studio/onion-weekly/config"
	"github.com/onion-studio/onion-weekly/domain"
	"github.com/onion-studio/onion-weekly/dto"
	"github.com/onion-studio/onion-weekly/web/middleware"
)

type recurringHandler struct {
	appConf          config.AppConf
	pgxPool          *pgxpool.Pool
	userService      *domain.UserService
	workspaceService *domain.WorkspaceService
	recurringService *domain.RecurringService
}

func NewRecurringHandler(
	appConf config.AppConf,
	pgxPool *pgxpool.Pool,
	userService *domain.UserService,
	workspaceService *domain.WorkspaceService,
	recurringService *domain.RecurringService,
) *recurringHandler {
	return &recurringHandler{
		appConf:          appConf,
		pgxPool:          pgxPool,
		userService:      userService,
		workspaceService: workspaceService,
		recurringService: recurringService,
	}
}

func (h *recurringHandler) Register(g *echo.Group) {
	g.Use(middleware.Transaction(h.appConf, h.pgxPool))
	g.Use(em.JWT(h.appConf.Secret), middleware.Actor(h.pgxPool, h.userService))
	g.GET("/workspaces/:workspace-id/recurrings", h.handleGetRecurrings)
	g.POST("/workspaces/:workspace-id/recurrings", h.handlePostRecurring)
	g.GET("/workspaces/:workspace-id/recurring-records", h.handleGetRecurringRecords)
	g.POST("/recurrings/:recurring-id/recurring-records", h.handlePostRecurringRecord)
	//g.POST("/register", h.handlePostUser)
	//g.GET("/me", h.handleGetTokenPayload, em.JWT(h.appConf.Secret))
}

func (h *recurringHandler) handlePostRecurring(c echo.Context) (err error) {
	tx := c.Get("tx").(pgx.Tx)
	actor := c.Get("actor").(dto.User)
	input := dto.CreateRecurringInput{}
	workspaceId := dto.UUID{}
	if err = c.Bind(&input); err != nil {
		return err
	}
	if err := workspaceId.Set(c.Param("workspace-id")); err != nil {
		return err
	}
	workspace, err := h.workspaceService.GetWorkspaceById(tx, workspaceId)
	if err != nil {
		return err
	}
	recurring, err := h.recurringService.CreateRecurring(tx, actor, workspace, input)
	if err != nil {
		return err
	}
	return c.JSON(200, recurring)
}

func (h *recurringHandler) handleGetRecurrings(c echo.Context) (err error) {
	tx := c.Get("tx").(pgx.Tx)
	actor := c.Get("actor").(dto.User)
	workspaceId := dto.UUID{}
	if err := workspaceId.Set(c.Param("workspace-id")); err != nil {
		return err
	}
	workspace, err := h.workspaceService.GetWorkspaceById(tx, workspaceId)
	if err != nil {
		return err
	}
	recurrings, err := h.recurringService.FindRecurringsByWorkspace(tx, actor, workspace)
	if err != nil {
		return err
	}
	return c.JSON(200, recurrings)
}

func (h *recurringHandler) handleGetRecurringRecords(c echo.Context) (err error) {
	tx := c.Get("tx").(pgx.Tx)
	actor := c.Get("actor").(dto.User)
	workspaceId := dto.UUID{}
	if err := workspaceId.Set(c.Param("workspace-id")); err != nil {
		return err
	}
	workspace, err := h.workspaceService.GetWorkspaceById(tx, workspaceId)
	if err != nil {
		return err
	}
	records, err := h.recurringService.FindRecurringRecordsByWorkspace(tx, actor, workspace)
	if err != nil {
		return err
	}
	return c.JSON(200, records)
}

func (h *recurringHandler) handlePostRecurringRecord(c echo.Context) (err error) {
	tx := c.Get("tx").(pgx.Tx)
	actor := c.Get("actor").(dto.User)
	input := dto.CreateRecurringRecordInput{}
	recurringId := dto.UUID{}
	if err = c.Bind(&input); err != nil {
		return err
	}
	if err := recurringId.Set(c.Param("recurring-id")); err != nil {
		return err
	}
	recurring, err := h.recurringService.GetRecurringById(tx, actor, recurringId)
	record, err := h.recurringService.CreateRecurringRecord(tx, actor, recurring, input)
	if err != nil {
		return err
	}
	return c.JSON(200, record)
}

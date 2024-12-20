package http

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"gmvr.pw/boombox-web-runner/config"
	"gmvr.pw/boombox-web-runner/pkg/model"
)

type HttpRunnerController struct {
	config        *config.HttpEntrypointConfig
	server        *echo.Echo
	runnerService RunnerService
	logger        *slog.Logger
}

func NewHttpRunnerController(
	config *config.HttpEntrypointConfig,
	logger *slog.Logger,
) (*HttpRunnerController, error) {
	return &HttpRunnerController{config: config, server: echo.New(), logger: logger}, nil
}

func (c *HttpRunnerController) Init(runnerService RunnerService) error {
	c.runnerService = runnerService

	c.server.POST("/runners", c.Create)
	c.server.DELETE("/runners/:id", c.Delete)

	return nil
}

func (c *HttpRunnerController) Serve() error {
	return c.server.Start(fmt.Sprintf(":%d", c.config.Port))
}

func (c *HttpRunnerController) Create(ctx echo.Context) error {
	var err error

	runner := &model.Runner{}
	if err := ctx.Bind(&runner); err != nil {
		c.logger.Error("cannot bind runner", "error", err)
		ctx.String(echo.ErrBadRequest.Code, "")
	}
	c.logger.Info("runner binded")

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ctx.RealIP(), runner.Port))
	if err != nil {
		c.logger.Error("cannot resolve udp address", "error", err)
		return nil
	}
	c.logger.Info("upd address resolved")

	con, err := net.Dial("udp", addr.String())
	if err != nil {
		c.logger.Error("cannot create udp dial", "error", err)
		return nil
	}
	c.logger.Info("upd dial created")

	out := make(chan []byte)
	go func() {
		defer func() {
			err := con.Close()
			if err != nil {
				c.logger.Error("cannot close connection", "error", err)
			}
		}()
		for o := range out {
			_, err := con.Write(o)
			if err != nil {
				c.logger.Error("cannot send music", "error", err)
				return
			}
		}
	}()

	err = c.runnerService.Run(runner, out)
	if err != nil {
		if _, ok := err.(*model.ModuleNotFoundError); ok {
			c.logger.Error("cannot find module for url", "url", runner.Url)
			ctx.String(echo.ErrBadRequest.Code, "")
			close(out)
			return nil
		}
		c.logger.Error("cannot run url", "error", err)
		ctx.String(echo.ErrInternalServerError.Code, "")
		close(out)
		return nil
	}
	c.logger.Info("runner started")

	err = ctx.JSON(http.StatusOK, runner)
	if err != nil {
		c.logger.Info("cannot send response", "error", err)
		return err
	}

	return nil
}

func (c *HttpRunnerController) Delete(ctx echo.Context) error {
	id := ctx.Param("id")

	c.logger.Info("delete requested", "id", id)

	runner, err := c.runnerService.Stop(id)
	if err != nil {
		c.logger.Info("cannot stop runner", "error", err)
		if _, ok := err.(*model.RunnerNotFoundError); ok {
			ctx.JSON(echo.ErrNotFound.Code, &model.Runner{ID: id})
			return err
		}
		ctx.String(echo.ErrInternalServerError.Code, "")
		return err
	}

	return ctx.JSON(http.StatusOK, runner)
}

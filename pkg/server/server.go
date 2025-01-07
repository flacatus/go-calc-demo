package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/liornabat-sealights/go-calc-demo/lib/types"
	"github.com/liornabat-sealights/go-calc-demo/pkg/calc"
)

type Server struct {
	echoWebServer *echo.Echo
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start(ctx context.Context) error {
	s.echoWebServer = echo.New()
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	fmt.Println("ABC")
	e.GET("/add", func(c echo.Context) error {
		return s.handelAdd(c)
	})
	e.GET("/sub", func(c echo.Context) error {
		return s.handelSub(c)
	})
	e.GET("/mul", func(c echo.Context) error {
		return s.handelMul(c)
	})
	e.GET("/div", func(c echo.Context) error {
		return s.handelDiv(c)
	})
	s.echoWebServer = e
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.echoWebServer.Start("0.0.0.0:8080")
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
		e.Logger.Infof("api server started at port 8080")
		return nil
	case <-time.After(1 * time.Second):
		return nil
	case <-ctx.Done():
		return fmt.Errorf("error strarting api server, %w", ctx.Err())
	}

}

func (s *Server) handelAdd(c echo.Context) error {
	ca, cb, err := s.parseParams(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	result := types.NewResultResponse().SetValues(ca, cb).SetResult(calc.Add(ca, cb))
	return c.JSON(http.StatusOK, result)
}

func (s *Server) handelSub(c echo.Context) error {
	a, b, err := s.parseParams(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	result := types.NewResultResponse().SetValues(a, b).SetResult(calc.Subtract(a, b))
	return c.JSON(http.StatusOK, result)
}

func (s *Server) handelMul(c echo.Context) error {
	a, b, err := s.parseParams(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	result := types.NewResultResponse().SetValues(a, b).SetResult(calc.Multiply(a, b))
	return c.JSON(http.StatusOK, result)
}

func (s *Server) handelDiv(c echo.Context) error {
	a, b, err := s.parseParams(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	result := types.NewResultResponse().SetValues(a, b).SetResult(calc.Divide(a, b))
	return c.JSON(http.StatusOK, result)
}

func (s *Server) parseParams(c echo.Context) (float64, float64, error) {
	var a, b float64
	valueA := c.QueryParam("a")
	valueB := c.QueryParam("b")
	if valueA == "" || valueB == "" {
		return 0, 0, fmt.Errorf("missing params")
	}
	if _, err := fmt.Sscanf(valueA, "%f", &a); err != nil {
		return 0, 0, fmt.Errorf("invalid param a")
	}
	if _, err := fmt.Sscanf(valueB, "%f", &b); err != nil {
		return 0, 0, fmt.Errorf("invalid param b")
	}
	return a, b, nil
}

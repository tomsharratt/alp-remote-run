package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tomsharratt/alp/evaluator"
	"github.com/tomsharratt/alp/lexer"
	"github.com/tomsharratt/alp/object"
	"github.com/tomsharratt/alp/parser"
)

type ExecuteRequest struct {
	Files []File
}

type File struct {
	Name    string
	Path    string
	Content string
}

type ExecuteRepsonse struct {
	Output string   `json:"output,omitempty"`
	Errors []string `json:"errors,omitempty"`
}

func handleExecute(c *gin.Context) {
	var req ExecuteRequest
	var res ExecuteRepsonse

	c.BindJSON(&req)

	if len(req.Files) == 0 {
		res.Errors = append(res.Errors, "no files provided.")
		c.JSON(http.StatusBadRequest, res)
		return
	}

	entry := req.Files[0]

	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	l := lexer.New(entry.Content)
	p := parser.New(l)

	program, err := p.ParseProgram(ctx)
	if err != nil {
		res.Errors = append(res.Errors, "program took too long to parse.")
		c.JSON(http.StatusBadRequest, res)
		return
	}
	if len(p.Errors()) != 0 {
		res.Errors = append(res.Errors, p.Errors()...)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	evaluated, err := evaluator.Eval(ctx, program, object.NewEnvironment())
	if err != nil {
		res.Errors = append(res.Errors, "program took too long to evaluate.")
		c.JSON(http.StatusBadRequest, res)
		return
	}
	res.Output = evaluated.Inspect()

	c.JSON(http.StatusOK, res)
}

func main() {
	r := gin.Default()

	r.POST("/execute", handleExecute)

	r.Run()
}

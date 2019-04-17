package covercheck

import (
	"context"
	"net/http"

	"golang.org/x/sync/errgroup"
)

// Controller provides handling checklist for each healthcheck path.
type Controller struct {
	path      string
	checklist []Checker
	fail      Response
	success   Response
}

// NewController returns Controller object.
func NewController(path string, checklist ...Checker) *Controller {
	if checklist == nil {
		checklist = []Checker{}
	}
	return &Controller{path, checklist, newFail(), newSuccess()}
}

// AddChecker provides embedding checker into controller.
func (c *Controller) AddChecker(checker Checker) {
	c.checklist = append(c.checklist, checker)
}

// HandlerFunc returns http.HandlerFunc implementation for running checklist
func (c *Controller) HandlerFunc() (string, http.HandlerFunc) {
	return c.path, c.runChecklist
}

func (c *Controller) runChecklist(w http.ResponseWriter, r *http.Request) {
	eg, ctx := errgroup.WithContext(context.Background())
	for _, checker := range c.checklist {
		eg.Go(func() error {
			return checker(ctx)
		})
	}
	if err := eg.Wait(); err != nil {
		c.fail.Render(w)
		return
	}
	c.success.Render(w)
}

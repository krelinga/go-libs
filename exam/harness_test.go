package exam_test

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/krelinga/go-deep/exam"
)

func TestHarness(t *testing.T) {
	t.Run("Run", func(t *testing.T) {
		h := exam.Harness{
			Name: "root",
		}
		log := h.Run(func(e exam.E) {
			e.Log("in root")
			e.Run("sub", func(e exam.E) {
				e.Log("in sub-test")
			})
		})
		if len(log.Log) != 1 || log.Log[0] != "in root" {
			t.Errorf("unexpected root log: %v", log.Log)
		}
		if log, ok := log.Find("sub"); !ok {
			t.Errorf("sub log not found")
		} else {
			if len(log.Log) != 1 || log.Log[0] != "in sub-test" {
				t.Errorf("unexpected sub log: %v", log.Log)
			}
		}
	})

	t.Run("Chdir", func(t *testing.T) {
		h := exam.Harness{}
		var origDir string
		h.Run(func(e exam.E) {
			var err error
			origDir, err = os.Getwd()
			if err != nil {
				t.Fatalf("Getwd failed: %v", err)
			}
			e.Chdir("/")
			newDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("Getwd failed: %v", err)
			}
			if newDir != "/" {
				t.Errorf("expected /, got %q", newDir)
			}
		})
		finalDir, err := os.Getwd()
		if err != nil {
			t.Fatalf("Getwd failed: %v", err)
		}
		if finalDir != origDir {
			t.Errorf("expected %q, got %q", origDir, finalDir)
		}
	})

	t.Run("Cleanup", func(t *testing.T) {
		h := exam.Harness{}
		called := false
		h.Run(func(e exam.E) {
			e.Cleanup(func() {
				called = true
			})
		})
		if !called {
			t.Errorf("cleanup not called")
		}
	})

	t.Run("Context", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			h := exam.Harness{}
			h.Run(func(e exam.E) {
				ctx := e.Context()
				if ctx == nil {
					t.Errorf("expected non-nil context")
				}
			})
		})

		t.Run("with value", func(t *testing.T) {
			type keyType struct{}
			parentCtx := context.WithValue(context.Background(), keyType{}, "test-value")
			h := exam.Harness{
				Ctx: parentCtx,
			}
			h.Run(func(e exam.E) {
				ctx := e.Context()
				v, ok := ctx.Value(keyType{}).(string)
				if !ok {
					t.Errorf("expected string value, got %T", ctx.Value(keyType{}))
				}
				if v != "test-value" {
					t.Errorf("expected test-value, got %q", v)
				}
			})
		})

		t.Run("is eventually canceled", func(t *testing.T) {
			deadlineCh := time.After(1 * time.Second)
			wait := sync.WaitGroup{}
			wait.Add(1)
			h := exam.Harness{}
			h.Run(func(e exam.E) {
				ctx := e.Context()
				go func() {
					select {
					case <-deadlineCh:
						// Timeout
						t.Errorf("context was not canceled in time")
					case <-ctx.Done():
						// Canceled as expected
					}
					wait.Done()
				}()
			})
			wait.Wait()
		})

		t.Run("sub-test is eventually canceled", func(t *testing.T) {
			deadlineCh := time.After(1 * time.Second)
			wait := sync.WaitGroup{}
			wait.Add(1)
			h := exam.Harness{}
			h.Run(func(e exam.E) {
				e.Run("sub", func(e exam.E) {
					ctx := e.Context()
					go func() {
						select {
						case <-deadlineCh:
							// Timeout
							t.Errorf("context was not canceled in time")
						case <-ctx.Done():
							// Canceled as expected
						}
						wait.Done()
					}()
				})
			})
			wait.Wait()
		})
	})
}

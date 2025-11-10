package exam_test

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/krelinga/go-deep"
	"github.com/krelinga/go-deep/exam"
	"github.com/krelinga/go-deep/match"
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

	t.Run("Deadline", func(t *testing.T) {
		tests := []struct {
			name     string
			deadline time.Time
			wantOk   bool
		}{
			{"no deadline", time.Time{}, false},
			{"with deadline", time.Now().Add(1 * time.Hour), true},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := exam.Harness{
					Deadline: tt.deadline,
				}
				h.Run(func(e exam.E) {
					dl, ok := e.Deadline()
					if ok != tt.wantOk {
						t.Errorf("expected ok=%v, got %v", tt.wantOk, ok)
					}
					if ok && !dl.Equal(tt.deadline) {
						t.Errorf("expected deadline %v, got %v", tt.deadline, dl)
					}
				})
			})
		}
	})

	t.Run("Recording Methods", func(t *testing.T) {
		tests := []struct {
			name    string
			f       func(exam.E)
			wantLog *exam.Log
		}{
			{
				name: "Error",
				f: func(e exam.E) {
					e.Error("test error")
				},
				wantLog: &exam.Log{
					Error: []string{"test error"},
					Fail:  true,
				},
			},
			{
				name: "Errorf",
				f: func(e exam.E) {
					e.Errorf("test %s", "errorf")
				},
				wantLog: &exam.Log{
					Error: []string{"test errorf"},
					Fail:  true,
				},
			},
			{
				name: "Fail",
				f: func(e exam.E) {
					e.Fail()
				},
				wantLog: &exam.Log{
					Fail: true,
				},
			},
			{
				name: "FailNow",
				f: func(e exam.E) {
					e.FailNow()
				},
				wantLog: &exam.Log{
					Fail:    true,
					FailNow: true,
				},
			},
			{
				name: "Fatal",
				f: func(e exam.E) {
					e.Fatal("fatal error")
				},
				wantLog: &exam.Log{
					Fatal:   []string{"fatal error"},
					Fail:    true,
					FailNow: true,
				},
			},
			{
				name: "Fatalf",
				f: func(e exam.E) {
					e.Fatalf("fatal %s", "errorf")
				},
				wantLog: &exam.Log{
					Fatal:   []string{"fatal errorf"},
					Fail:    true,
					FailNow: true,
				},
			},
			{
				name: "Helper",
				f: func(e exam.E) {
					e.Helper()
				},
				wantLog: &exam.Log{
					Helper: true,
				},
			},
			{
				name: "Log",
				f: func(e exam.E) {
					e.Log("a log message")
				},
				wantLog: &exam.Log{
					Log: []string{"a log message"},
				},
			},
			{
				name: "Logf",
				f: func(e exam.E) {
					e.Logf("log %s", "message")
				},
				wantLog: &exam.Log{
					Log: []string{"log message"},
				},
			},
			{
				name: "Parallel",
				f: func(e exam.E) {
					e.Parallel()
				},
				wantLog: &exam.Log{
					Parallel: true,
				},
			},
			{
				name: "Skip",
				f: func(e exam.E) {
					e.Skip("skip message")
				},
				wantLog: &exam.Log{
					Skip:    []string{"skip message"},
					SkipNow: true,
				},
			},
			{
				name: "Skipf",
				f: func(e exam.E) {
					e.Skipf("skip %s", "message")
				},
				wantLog: &exam.Log{
					Skip:    []string{"skip message"},
					SkipNow: true,
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := exam.Harness{}
				gotLog := h.Run(func(e exam.E) {
					tt.f(e)
				})
				matcher := match.Equal(tt.wantLog)
				result := matcher.Match(deep.NewEnv(), match.NewVals1(gotLog))
				if !result.Matched() {
					t.Errorf("log mismatch:\n%s", result)
				}
			})
		}
	})

	t.Run("Name", func(t *testing.T) {
		h := exam.Harness{
			Name: "root-name",
		}
		h.Run(func(e exam.E) {
			if e.Name() != "root-name" {
				t.Errorf("expected root-name, got %q", e.Name())
			}
			e.Run("sub-name", func(e exam.E) {
				if e.Name() != "root-name/sub-name" {
					t.Errorf("expected root-name/sub-name, got %q", e.Name())
				}
			})
		})
	})
}

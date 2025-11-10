package exam_test

import (
	"os"
	"testing"

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
}

package jitter_test

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/jitter"
)

func TestNew_StoresBaseAndFactor(t *testing.T) {
	j := jitter.New(5*time.Second, 0.1)
	if j.Base() != 5*time.Second {
		t.Fatalf("expected base 5s, got %v", j.Base())
	}
	if j.Factor() != 0.1 {
		t.Fatalf("expected factor 0.1, got %v", j.Factor())
	}
}

func TestNew_PanicsOnZeroBase(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero base")
		}
	}()
	jitter.New(0, 0.1)
}

func TestNew_PanicsOnZeroFactor(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero factor")
		}
	}()
	jitter.New(time.Second, 0)
}

func TestNew_PanicsOnFactorAboveOne(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for factor > 1")
		}
	}()
	jitter.New(time.Second, 1.1)
}

func TestNext_WithinBounds(t *testing.T) {
	base := 10 * time.Second
	factor := 0.2
	j := jitter.New(base, factor)

	lo := time.Duration(float64(base) * (1 - factor))
	hi := time.Duration(float64(base) * (1 + factor))

	for i := 0; i < 1000; i++ {
		d := j.Next()
		if d < lo || d > hi {
			t.Fatalf("Next() = %v out of expected range [%v, %v]", d, lo, hi)
		}
	}
}

func TestNext_AlwaysPositive(t *testing.T) {
	j := jitter.New(time.Millisecond, 1.0)
	for i := 0; i < 500; i++ {
		if d := j.Next(); d <= 0 {
			t.Fatalf("Next() returned non-positive duration: %v", d)
		}
	}
}

func TestNext_ConcurrentSafe(t *testing.T) {
	j := jitter.New(time.Second, 0.3)
	done := make(chan struct{})
	for i := 0; i < 8; i++ {
		go func() {
			for k := 0; k < 200; k++ {
				_ = j.Next()
			}
			done <- struct{}{}
		}()
	}
	for i := 0; i < 8; i++ {
		<-done
	}
}

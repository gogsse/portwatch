package sampler_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/sampler"
)

func TestNew_PanicsOnNegativeRate(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative rate")
		}
	}()
	sampler.New(-0.1)
}

func TestNew_ClampsRateAboveOne(t *testing.T) {
	s := sampler.New(2.5)
	if s.Rate() != 1.0 {
		t.Fatalf("expected rate clamped to 1.0, got %v", s.Rate())
	}
}

func TestAllow_RateOneAlwaysTrue(t *testing.T) {
	s := sampler.New(1.0)
	for i := 0; i < 100; i++ {
		if !s.Allow(8080) {
			t.Fatal("expected Allow to return true with rate=1.0")
		}
	}
}

func TestAllow_RateZeroAlwaysFalse(t *testing.T) {
	s := sampler.New(0.0)
	for i := 0; i < 100; i++ {
		if s.Allow(8080) {
			t.Fatal("expected Allow to return false with rate=0.0")
		}
	}
}

func TestCount_TracksCallsRegardlessOfDecision(t *testing.T) {
	s := sampler.New(0.0) // never emits, but must still count
	port := 443
	for i := 0; i < 10; i++ {
		s.Allow(port)
	}
	if got := s.Count(port); got != 10 {
		t.Fatalf("expected count 10, got %d", got)
	}
}

func TestCount_IndependentPerPort(t *testing.T) {
	s := sampler.New(1.0)
	s.Allow(80)
	s.Allow(80)
	s.Allow(443)

	if got := s.Count(80); got != 2 {
		t.Fatalf("expected count 2 for port 80, got %d", got)
	}
	if got := s.Count(443); got != 1 {
		t.Fatalf("expected count 1 for port 443, got %d", got)
	}
}

func TestReset_ClearsCounters(t *testing.T) {
	s := sampler.New(1.0)
	s.Allow(8080)
	s.Allow(8080)
	s.Reset()

	if got := s.Count(8080); got != 0 {
		t.Fatalf("expected count 0 after Reset, got %d", got)
	}
}

func TestAllow_PartialRateEmitsSomeFraction(t *testing.T) {
	s := sampler.New(0.5)
	allowed := 0
	const iterations = 10_000
	for i := 0; i < iterations; i++ {
		if s.Allow(9000) {
			allowed++
		}
	}
	// With rate=0.5 expect roughly 50% ± 5%
	lo, hi := iterations*45/100, iterations*55/100
	if allowed < lo || allowed > hi {
		t.Fatalf("expected allowed in [%d, %d], got %d", lo, hi, allowed)
	}
}

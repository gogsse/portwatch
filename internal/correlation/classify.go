package correlation

// BurstClass describes the inferred nature of a correlated burst.
type BurstClass int

const (
	// ClassNoise indicates too few events to draw a conclusion.
	ClassNoise BurstClass = iota
	// ClassRestart suggests a service restart pattern (mixed open/close).
	ClassRestart
	// ClassIntrusion suggests new listeners without corresponding closes.
	ClassIntrusion
	// ClassSweep suggests many ports closed simultaneously.
	ClassSweep
)

func (bc BurstClass) String() string {
	switch bc {
	case ClassRestart:
		return "restart"
	case ClassIntrusion:
		return "intrusion"
	case ClassSweep:
		return "sweep"
	default:
		return "noise"
	}
}

// Classify inspects a Burst and returns the most likely BurstClass.
// threshold is the minimum event count before a burst is classified.
func Classify(b *Burst, threshold int) BurstClass {
	if b == nil || len(b.Events) < threshold {
		return ClassNoise
	}

	var opened, closed int
	for _, e := range b.Events {
		switch e.Kind {
		case Opened:
			opened++
		case Closed:
			closed++
		}
	}

	switch {
	case opened > 0 && closed > 0:
		return ClassRestart
	case opened > closed:
		return ClassIntrusion
	default:
		return ClassSweep
	}
}

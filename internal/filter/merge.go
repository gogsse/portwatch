package filter

// Merge combines two Filters into a new Filter whose allowed set is the
// union of both. ExemptPrivileged is true if either source filter has it
// enabled.
func Merge(a, b *Filter) *Filter {
	union := make(map[int]struct{}, len(a.allowed)+len(b.allowed))

	for p := range a.allowed {
		union[p] = struct{}{}
	}
	for p := range b.allowed {
		union[p] = struct{}{}
	}

	return &Filter{
		allowed:         union,
		exemptPrivileged: a.exemptPrivileged || b.exemptPrivileged,
	}
}

// Ports returns the sorted list of explicitly allowed ports held by the filter.
func (f *Filter) Ports() []int {
	result := make([]int, 0, len(f.allowed))
	for p := range f.allowed {
		result = append(result, p)
	}
	sortInts(result)
	return result
}

// sortInts is a small helper to avoid importing sort in callers.
func sortInts(s []int) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}

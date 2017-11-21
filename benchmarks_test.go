package integration

import "testing"

// Run benchmarks in order to manage the dependency chain, and print a
// summary afterwards.
func BenchmarkMain(b *testing.B) {
	b.Run("Create", BenchmarkLoginCreateAccount)
	b.Run("Access", BenchmarkLoginRequestAccess)
	BenchmarkLoginSummary(b)
}

// Goal
// Build a 3-stage concurrent pipeline:
// generator → transformer → aggregator
// Stage 1: generator
// Generates integers 1..50
// Sends numbers to transformer
// Stage 2: transformer
// For each input n, compute f(n) = n * n
// Must run using K concurrent transformer workers
// Must obey backpressure (slow downstream should slow upstream)
// Stage 3: aggregator
// Receives squared values
// Computes:
// count
// sum
// average
// max
// min
// After transformer closes, aggregator produces a summary struct.
// Constraints
// No sync.WaitGroup
// You must use fan-out + fan-in correctly
// Pipeline must stop cleanly after all data is consumed
// aggregator must only close after all transformer workers finish
// Required Output
// Print final summary as:
// Count: ?
// Sum: ?
// Avg: ?
// Max: ?
// Min: ?

package main

func main() {

}

package main

func Process(b []float32) []float32 {
    for i := range b {
        b[i] *= 0.5
    }
    return b
}

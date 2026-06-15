package main

func Process(b []float32) []float32 {
	for i := range b {
		b[i] *= 0.5
	}
	return b
}

func Metadata() map[string]string {
	return map[string]string{
		"name":        "Example Mixer",
		"description": "Simple plugin that halves the output volume.",
		"version":     "0.1.0",
		"author":      "Iconictacoma",
	}
}

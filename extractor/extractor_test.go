package extractor_test

func mergeBuffers(buffers ...[]byte) []byte {
	combined := make([]byte, 0)

	for _, buf := range buffers {
		combined = append(combined, buf...)
	}

	return combined
}

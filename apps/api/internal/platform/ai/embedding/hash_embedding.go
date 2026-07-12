package embedding

import (
	"crypto/sha256"
	"encoding/binary"
	"math"
	"strconv"
)

const DefaultEmbeddingDim = 384

// HashEmbedding is a dev-only deterministic embedding used by both:
// - AI Runtime worker (Python)
// - API chat retrieval (Go)
//
// It must stay in sync with apps/ai-runtime/worker/embedding.py.
func HashEmbedding(text string, dim int) []float64 {
	if dim <= 0 {
		dim = DefaultEmbeddingDim
	}

	vec := make([]float64, dim)
	var sumSq float64

	for i := 0; i < dim; i++ {
		// Match Python: hashlib.sha256(f"{text}\0{i}".encode("utf-8")).digest()
		h := sha256.Sum256([]byte(text + "\x00" + strconv.Itoa(i)))
		u := binary.LittleEndian.Uint32(h[:4])

		// Match Python: (u % 10_000_000) / 10_000_000.0
		v := float64(u%10_000_000) / 10_000_000.0
		vec[i] = v
		sumSq += v * v
	}

	norm := math.Sqrt(sumSq) + 1e-12
	for i := range vec {
		vec[i] /= norm
	}

	return vec
}

func CosineSimilarity(a, b []float64) float64 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	var dot float64
	for i := 0; i < n; i++ {
		dot += a[i] * b[i]
	}
	return dot
}


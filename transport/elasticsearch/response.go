package transport

type EsHits[T any] struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source T `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

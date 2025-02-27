package metric

type Metric struct {
	Alloc         int64
	BuckHashSys   int64
	Frees         int64
	GCCPUFraction float64
	GCSys         int64
	HeapAlloc     int64
	HeapIdle      int64
	HeapInuse     int64
	HeapObjects   int64
	HeapReleased  int64
	HeapSys       int64
	LastGC        int64
	Lookups       int64
	MCacheInuse   int64
	MCacheSys     int64
	MSpanInuse    int64
	MSpanSys      int64
	Mallocs       int64
	NextGC        int64
	NumForcedGC   int64
	NumGC         int32
	OtherSys      int64
	PauseTotalNs  int64
	StackInuse    int64
	StackSys      int64
	Sys           int64
	TotalAlloc    int64
	PollCount     int64
	RandomValue   float64
}

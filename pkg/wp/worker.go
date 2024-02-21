package wp

type WorkerPool struct {
	JobQueue   chan []byte
	WorkerFunc func([]byte)
	MaxWorkers int
}

func NewWorkerPool(workerFunc func([]byte), maxWorkers int) *WorkerPool {
	return &WorkerPool{
		JobQueue:   make(chan []byte, maxWorkers),
		WorkerFunc: workerFunc,
		MaxWorkers: maxWorkers,
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.MaxWorkers; i++ {
		go func() {
			for job := range wp.JobQueue {
				wp.WorkerFunc(job)
			}
		}()
	}
}

func (wp *WorkerPool) AddJob(job []byte) {
	wp.JobQueue <- job
}

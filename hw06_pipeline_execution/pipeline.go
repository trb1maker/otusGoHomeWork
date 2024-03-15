package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func limitter(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for v := range in {
			select {
			case <-done:
				return
			default:
				out <- v
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = stage(limitter(in, done))
	}
	return in
}

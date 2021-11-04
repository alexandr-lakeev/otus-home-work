package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = proxyStage(stage, in, done)
	}

	return in
}

func proxyStage(stage Stage, in In, done In) Out {
	proxy := make(Bi)
	out := stage(proxy)

	go func() {
		defer close(proxy)
		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				proxy <- val
			}
		}
	}()

	return out
}

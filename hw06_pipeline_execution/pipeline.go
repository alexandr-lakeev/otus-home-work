package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var proxy Bi
	var out Out

	for _, stage := range stages {
		proxy = make(Bi)
		out = stage(proxy)

		go func(in In, proxy Bi) {
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
		}(in, proxy)

		in = out
	}

	return out
}

package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	result := make(Bi)
	go func() {
		defer close(result)
		pipelineChannel := in
		for _, stage := range stages {
			pipelineChannel = execStage(stage, pipelineChannel, done)
		}
		for numVal := range pipelineChannel {
			result <- numVal
		}
	}()
	return result
}

func execStage(stage Stage, in In, done In) Out {
	result := make(Bi)
	go func() {
		defer close(result)
		for resStage := range stage(in) {
			select {
			case <-done:
				return
			default:
				result <- resStage
			}
		}
	}()
	return result
}

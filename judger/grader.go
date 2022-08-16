package judger

type Grader interface {
	Grade()
}

type grader struct {
}

func NewGrader() *grader {
	return &grader{}
}

func (g *grader) Grade() {

}

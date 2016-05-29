package frame

type CommentSequencer interface {
	Sequencer

	Comment(language string, description string) CommentFramer
	AddComment(CommentFramer)
}

type CommentSequence struct {
	sequence map[string]CommentFramer
}

func NewCommentSequence() *CommentSequence {
	return &CommentSequence{
		sequence: make(map[string]CommentFramer),
	}
}

func (cs CommentSequence) Frames() []Framer {
	frames := []Framer{}
	for _, f := range cs.sequence {
		frames = append(frames, f)
	}
	return frames
}

func (cs CommentSequence) Comment(language string, description string) CommentFramer {
	return cs.sequence[language+description]
}

func (cs *CommentSequence) AddComment(cf CommentFramer) {
	id := cf.Language() + cf.Description()
	cs.sequence[id] = cf
}

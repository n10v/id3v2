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
	var (
		i      = 0
		frames = make([]Framer, len(cs.sequence))
	)

	for _, v := range cs.sequence {
		frames[i] = v
		i++
	}
	return frames
}

func (cs CommentSequence) Comment(language string, description string) CommentFramer {
	return cs.sequence[language+description]
}

func (cs *CommentSequence) AddFrame(f Framer) {
	cf := f.(CommentFramer)
	id := cf.Language() + cf.Description()
	cs.sequence[id] = cf
}

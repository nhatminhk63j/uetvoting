package event

import (
	"time"

	pb "github.com/nhatminhk63j/uetvoting/pb/event/v1"
)

type Event struct {
	ID           int `gorm:"primaryKey"`
	Name         string
	Participants string
	IsApproved   bool
	IsOpened     bool
	CreatedBy    int
	UpdatedBy    int
	UpdatedAt    time.Time
	Sections     []Section
}

type Section struct {
	ID          int `gorm:"primaryKey"`
	Name        string
	Description string
	Position    int
	EventID     int
	Questions   []Question
}

type Question struct {
	ID            int `gorm:"primaryKey"`
	Content       string
	Type          string
	IsRequired    bool
	LimitedChoice int
	Position      int
	SectionID     int
	Options       []Option
}

type Option struct {
	ID         int `gorm:"primaryKey"`
	Content    string
	QuestionID int
}

// ToEventUpsertInfo ...
func ToEventUpsertInfo(in *pb.UpsertEventRequest) *Event {
	sections := make([]Section, 0)
	for _, s := range in.Sections {
		questions := make([]Question, 0)
		for _, q := range s.Questions {
			options := make([]Option, 0)
			for _, o := range q.Options {
				options = append(options, Option{
					Content: o.Content,
				})
			}
			questions = append(questions, Question{
				Content:       q.Content,
				Position:      int(q.Position),
				Type:          q.Type,
				IsRequired:    q.IsRequired,
				LimitedChoice: int(q.LimitedChoice),
				Options:       options,
			})
		}
		sections = append(sections, Section{
			Name:        s.Name,
			Description: s.Description,
			Position:    int(s.Position),
			Questions:   questions,
		})
	}

	return &Event{
		ID:           int(in.Id),
		Name:         in.Name,
		Participants: in.Participants,
		Sections:     sections,
	}
}

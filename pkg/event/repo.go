package event

import (
	"context"

	"golang.org/x/xerrors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	UpsertEvent(ctx context.Context, event *Event) (eventID int, err error)
	GetEventByID(ctx context.Context, eventID int) (*Event, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{
		db: db,
	}
}

// UpsertEvent ...
func (r repo) UpsertEvent(ctx context.Context, newEvent *Event) (eventID int, err error) {
	if newEvent.ID == 0 {
		err = r.db.Create(newEvent).Error
		if err != nil {
			return 0, xerrors.Errorf("error creating newEvent: %w", err)
		}
		return newEvent.ID, nil
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		var oldEvent Event
		err := tx.Preload("Sections.Questions.Options").
			Preload("Sections.Questions").
			Preload(clause.Associations).
			First(&oldEvent, newEvent.ID).Error
		if err != nil {
			return xerrors.Errorf("error getting prevEvent: %w", err)
		}

		err = removeAllSections(ctx, tx, &oldEvent)
		if err != nil {
			return err
		}
		err = tx.Session(&gorm.Session{FullSaveAssociations: true}).Updates(newEvent).Error
		if err != nil {
			return xerrors.Errorf("error when full save associations: %w", err)
		}
		return nil
	})

	if err != nil {
		return 0, xerrors.Errorf("error updating newEvent: %w", err)
	}
	return newEvent.ID, nil
}

func removeAllSections(ctx context.Context, tx *gorm.DB, event *Event) error {
	sectionIDs := make([]interface{}, 0)
	questionIDs := make([]interface{}, 0)
	for _, section := range event.Sections {
		sectionIDs = append(sectionIDs, section.ID)
		for _, question := range section.Questions {
			questionIDs = append(questionIDs, question.ID)
		}
	}
	// delete all options
	if err := tx.Where("question_id in ?", questionIDs).Delete(Option{}).Error; err != nil {
		return xerrors.Errorf("error deleting options: %w", err)
	}

	// delete all questions
	if err := tx.Where("section_id in ?", sectionIDs).Delete(Question{}).Error; err != nil {
		return xerrors.Errorf("error deleting questions: %w", err)
	}

	// delete all sections
	if err := tx.Where("event_id = ?", event.ID).Delete(&Section{}).Error; err != nil {
		return xerrors.Errorf("error deleting sections: %w", err)
	}
	return nil
}

// GetEventByID ...
func (r repo) GetEventByID(ctx context.Context, eventID int) (*Event, error) {
	var event Event
	err := r.db.First(&event, eventID).Error
	if err != nil {
		return nil, xerrors.Errorf("error getting event by id: %w", err)
	}
	return &event, nil
}

package services

import (
	"context"
	"database/sql"
	"errors"
	"module_6/internal/database/repository"
	"module_6/internal/logger"
	"module_6/internal/models"
)

// NoteService handles note business logic
type NoteService struct {
	noteRepo     repository.NoteRepository
	categoryRepo repository.CategoryRepository
	logger       logger.Logger
}

func NewNoteService(
	noteRepo repository.NoteRepository,
	categoryRepo repository.CategoryRepository,
	log logger.Logger,
) *NoteService {
	return &NoteService{
		noteRepo:     noteRepo,
		categoryRepo: categoryRepo,
		logger:       log,
	}
}

func (s *NoteService) Create(ctx context.Context, userID int, req *models.NoteCreate) (*models.Note, error) {
	s.logger.Info("Creating note",
		logger.Int("user_id", userID),
		logger.String("title", req.Title))

	note := &models.Note{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	}

	if req.CategoryID != nil {
		category, err := s.categoryRepo.GetByID(ctx, *req.CategoryID)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				s.logger.Warn("Category not found",
					logger.Int("category_id", *req.CategoryID))
				return nil, models.ErrInvalidInput
			}
			s.logger.Error("Failed to get category", logger.Error(err))
			return nil, err
		}

		if category.UserID != userID {
			s.logger.Warn("Category access denied",
				logger.Int("category_id", *req.CategoryID),
				logger.Int("user_id", userID),
				logger.Int("owner_id", category.UserID))
			return nil, models.ErrForbidden
		}

		note.CategoryID = sql.NullInt64{
			Int64: int64(*req.CategoryID),
			Valid: true,
		}
		note.Category = category
	}

	if err := s.noteRepo.Create(ctx, note); err != nil {
		s.logger.Error("Failed to create note", logger.Error(err))
		return nil, err
	}

	s.logger.Info("Note created successfully",
		logger.Int("note_id", note.ID),
		logger.Int("user_id", userID))

	return note, nil
}

func (s *NoteService) GetAll(ctx context.Context, userID int, filter models.NoteFilter) ([]*models.Note, error) {
	s.logger.Debug("Getting notes for user",
		logger.Int("user_id", userID),
		logger.Any("filter", filter))

	notes, err := s.noteRepo.GetByUserID(ctx, userID, filter)
	if err != nil {
		s.logger.Error("Failed to get notes", logger.Error(err))
		return nil, err
	}

	s.logger.Debug("Notes retrieved",
		logger.Int("user_id", userID),
		logger.Int("count", len(notes)))

	return notes, nil
}

func (s *NoteService) GetByID(ctx context.Context, id, userID int) (*models.Note, error) {
	s.logger.Debug("Getting note",
		logger.Int("note_id", id),
		logger.Int("user_id", userID))

	note, err := s.noteRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			s.logger.Warn("Note not found", logger.Int("note_id", id))
		} else {
			s.logger.Error("Failed to get note", logger.Error(err))
		}
		return nil, err
	}

	if note.UserID != userID {
		s.logger.Warn("Access forbidden to note",
			logger.Int("note_id", id),
			logger.Int("user_id", userID),
			logger.Int("owner_id", note.UserID))
		return nil, models.ErrNotFound
	}

	return note, nil
}

func (s *NoteService) Update(ctx context.Context, id, userID int, req *models.NoteUpdate) (*models.Note, error) {
	s.logger.Info("Updating note",
		logger.Int("note_id", id),
		logger.Int("user_id", userID))

	note, err := s.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		note.Title = *req.Title
	}

	if req.Content != nil {
		note.Content = *req.Content
	}

	if req.CategoryID != nil {
		if *req.CategoryID == 0 {
			note.CategoryID = sql.NullInt64{Valid: false}
			note.Category = nil
		} else {
			category, err := s.categoryRepo.GetByID(ctx, *req.CategoryID)
			if err != nil {
				if errors.Is(err, models.ErrNotFound) {
					s.logger.Warn("Category not found",
						logger.Int("category_id", *req.CategoryID))
					return nil, models.ErrInvalidInput
				}
				s.logger.Error("Failed to get category", logger.Error(err))
				return nil, err
			}

			if category.UserID != userID {
				s.logger.Warn("Category access denied",
					logger.Int("category_id", *req.CategoryID),
					logger.Int("user_id", userID))
				return nil, models.ErrForbidden
			}

			note.CategoryID = sql.NullInt64{
				Int64: int64(*req.CategoryID),
				Valid: true,
			}
			note.Category = category
		}
	}

	if err := s.noteRepo.Update(ctx, note); err != nil {
		s.logger.Error("Failed to update note", logger.Error(err))
		return nil, err
	}

	s.logger.Info("Note updated successfully",
		logger.Int("note_id", id),
		logger.Int("user_id", userID))

	return note, nil
}

func (s *NoteService) Delete(ctx context.Context, id, userID int) error {
	s.logger.Info("Deleting note",
		logger.Int("note_id", id),
		logger.Int("user_id", userID))

	_, err := s.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	if err := s.noteRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete note", logger.Error(err))
		return err
	}

	s.logger.Info("Note deleted successfully",
		logger.Int("note_id", id),
		logger.Int("user_id", userID))

	return nil
}
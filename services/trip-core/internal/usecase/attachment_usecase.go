package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/makskozlov/atlasflow/services/trip-core/internal/domain"
)

type AttachmentUsecase struct {
	repo    TripRepository
	storage StoragePort
	clock   Clock
	ids     IDGen
}

func NewAttachmentUsecase(repo TripRepository, storage StoragePort, clock Clock, ids IDGen) *AttachmentUsecase {
	return &AttachmentUsecase{repo: repo, storage: storage, clock: clock, ids: ids}
}

type CreateUploadInput struct {
	TripID      string
	UserID      string
	PlaceID     *string
	DayID       *string
	Filename    string
	ContentType string
	FileSize    int64
	Visibility  domain.AttachmentVisibility
}

type UploadPresignedOutput struct {
	AttachmentID string
	UploadURL    string
	DownloadURL  string
	FileKey      string
}

func (u *AttachmentUsecase) CreateUpload(ctx context.Context, in CreateUploadInput) (UploadPresignedOutput, error) {
	t, err := u.repo.GetTrip(ctx, in.TripID)
	if err != nil {
		return UploadPresignedOutput{}, err
	}
	if err := t.AssertOwner(in.UserID); err != nil {
		return UploadPresignedOutput{}, err
	}

	if err := domain.ValidateAttachmentMeta(in.ContentType, in.FileSize); err != nil {
		return UploadPresignedOutput{}, err
	}

	attachmentID := u.ids.New()
	key := fmt.Sprintf("trips/%s/attachments/%s-%s", in.TripID, attachmentID, in.Filename)

	putURL, err := u.storage.PresignPut(ctx, key, in.ContentType, 15*time.Minute)
	if err != nil {
		return UploadPresignedOutput{}, fmt.Errorf("presign put url: %w", err)
	}

	getURL, err := u.storage.PresignGet(ctx, key, 24*time.Hour)
	if err != nil {
		return UploadPresignedOutput{}, fmt.Errorf("presign get url: %w", err)
	}

	vis := in.Visibility
	if vis == domain.AttachmentVisibilityUnspecified {
		vis = domain.AttachmentVisibilityPrivate // private by default
	}

	meta := domain.Attachment{
		ID:          attachmentID,
		TripID:      in.TripID,
		PlaceID:     in.PlaceID,
		DayID:       in.DayID,
		Filename:    in.Filename,
		ContentType: in.ContentType,
		FileSize:    in.FileSize,
		StorageKey:  key,
		DownloadURL: getURL,
		Visibility:  vis,
		IsReady:     false,
		CreatedAt:   u.clock.Now(),
	}

	if err := u.repo.CreateAttachment(ctx, meta); err != nil {
		return UploadPresignedOutput{}, fmt.Errorf("save attachment metadata: %w", err)
	}

	return UploadPresignedOutput{
		AttachmentID: attachmentID,
		UploadURL:    putURL,
		DownloadURL:  getURL,
		FileKey:      key,
	}, nil
}

func (u *AttachmentUsecase) ConfirmUpload(ctx context.Context, tripID, attachmentID, userID string) error {
	t, err := u.repo.GetTrip(ctx, tripID)
	if err != nil {
		return err
	}
	if err := t.AssertOwner(userID); err != nil {
		return err
	}

	att, err := u.repo.GetAttachment(ctx, attachmentID)
	if err != nil {
		return err
	}
	if att.TripID != tripID {
		return domain.ErrNotFound
	}

	return u.repo.ConfirmAttachmentReady(ctx, attachmentID)
}
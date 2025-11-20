package repository

import (
    "context"

    "github.com/google/uuid"
    "kerjakuy/internal/models"

    "gorm.io/gorm"
)

type ChatChannelRepository interface {
    Create(ctx context.Context, channel *models.ChatChannel) error
    FindByID(ctx context.Context, id uuid.UUID) (*models.ChatChannel, error)
    ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]models.ChatChannel, error)
}

type ChatChannelMemberRepository interface {
    AddMembers(ctx context.Context, members []models.ChatChannelMember) error
    RemoveMember(ctx context.Context, channelID, userID uuid.UUID) error
    ListMembers(ctx context.Context, channelID uuid.UUID) ([]models.ChatChannelMember, error)
}

type ChatMessageRepository interface {
    Create(ctx context.Context, message *models.ChatMessage) error
    ListByChannel(ctx context.Context, channelID uuid.UUID, limit int) ([]models.ChatMessage, error)
}

type ChatMessageReadRepository interface {
    MarkRead(ctx context.Context, read *models.ChatMessageRead) error
    ListByMessage(ctx context.Context, messageID uuid.UUID) ([]models.ChatMessageRead, error)
}

type chatChannelRepository struct {
    db *gorm.DB
}

type chatChannelMemberRepository struct {
    db *gorm.DB
}

type chatMessageRepository struct {
    db *gorm.DB
}

type chatMessageReadRepository struct {
    db *gorm.DB
}

func NewChatChannelRepository(db *gorm.DB) ChatChannelRepository {
    return &chatChannelRepository{db: db}
}

func NewChatChannelMemberRepository(db *gorm.DB) ChatChannelMemberRepository {
    return &chatChannelMemberRepository{db: db}
}

func NewChatMessageRepository(db *gorm.DB) ChatMessageRepository {
    return &chatMessageRepository{db: db}
}

func NewChatMessageReadRepository(db *gorm.DB) ChatMessageReadRepository {
    return &chatMessageReadRepository{db: db}
}

func (r *chatChannelRepository) Create(ctx context.Context, channel *models.ChatChannel) error {
    return r.db.WithContext(ctx).Create(channel).Error
}

func (r *chatChannelRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.ChatChannel, error) {
    var channel models.ChatChannel
    if err := r.db.WithContext(ctx).First(&channel, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &channel, nil
}

func (r *chatChannelRepository) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]models.ChatChannel, error) {
    var channels []models.ChatChannel
    if err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Find(&channels).Error; err != nil {
        return nil, err
    }
    return channels, nil
}

func (r *chatChannelMemberRepository) AddMembers(ctx context.Context, members []models.ChatChannelMember) error {
    if len(members) == 0 {
        return nil
    }
    return r.db.WithContext(ctx).Create(&members).Error
}

func (r *chatChannelMemberRepository) RemoveMember(ctx context.Context, channelID, userID uuid.UUID) error {
    return r.db.WithContext(ctx).Where("channel_id = ? AND user_id = ?", channelID, userID).Delete(&models.ChatChannelMember{}).Error
}

func (r *chatChannelMemberRepository) ListMembers(ctx context.Context, channelID uuid.UUID) ([]models.ChatChannelMember, error) {
    var members []models.ChatChannelMember
    if err := r.db.WithContext(ctx).Where("channel_id = ?", channelID).Find(&members).Error; err != nil {
        return nil, err
    }
    return members, nil
}

func (r *chatMessageRepository) Create(ctx context.Context, message *models.ChatMessage) error {
    return r.db.WithContext(ctx).Create(message).Error
}

func (r *chatMessageRepository) ListByChannel(ctx context.Context, channelID uuid.UUID, limit int) ([]models.ChatMessage, error) {
    var messages []models.ChatMessage
    query := r.db.WithContext(ctx).Where("channel_id = ?", channelID).Order("created_at desc")
    if limit > 0 {
        query = query.Limit(limit)
    }
    if err := query.Find(&messages).Error; err != nil {
        return nil, err
    }
    return messages, nil
}

func (r *chatMessageReadRepository) MarkRead(ctx context.Context, read *models.ChatMessageRead) error {
    return r.db.WithContext(ctx).Create(read).Error
}

func (r *chatMessageReadRepository) ListByMessage(ctx context.Context, messageID uuid.UUID) ([]models.ChatMessageRead, error) {
    var reads []models.ChatMessageRead
    if err := r.db.WithContext(ctx).Where("message_id = ?", messageID).Find(&reads).Error; err != nil {
        return nil, err
    }
    return reads, nil
}

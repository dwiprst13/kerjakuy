package service

import (
    "context"

    "github.com/google/uuid"
    "kerjakuy/internal/dto"
    "kerjakuy/internal/models"
    "kerjakuy/internal/repository"
)

type ChatService interface {
    CreateChannel(ctx context.Context, req dto.CreateChatChannelRequest, createdBy uuid.UUID) (*dto.ChatChannelDTO, error)
    ListWorkspaceChannels(ctx context.Context, workspaceID uuid.UUID) ([]dto.ChatChannelDTO, error)
    AddMembers(ctx context.Context, channelID uuid.UUID, userIDs []uuid.UUID) ([]dto.ChatChannelMemberDTO, error)
    RemoveMember(ctx context.Context, channelID, userID uuid.UUID) error
    SendMessage(ctx context.Context, req dto.CreateChatMessageRequest, senderID uuid.UUID) (*dto.ChatMessageDTO, error)
    ListMessages(ctx context.Context, channelID uuid.UUID, limit int) ([]dto.ChatMessageDTO, error)
    MarkMessageRead(ctx context.Context, messageID, userID uuid.UUID) error
}

type chatService struct {
    channelRepo repository.ChatChannelRepository
    memberRepo  repository.ChatChannelMemberRepository
    messageRepo repository.ChatMessageRepository
    readRepo    repository.ChatMessageReadRepository
}

func NewChatService(channelRepo repository.ChatChannelRepository, memberRepo repository.ChatChannelMemberRepository, messageRepo repository.ChatMessageRepository, readRepo repository.ChatMessageReadRepository) ChatService {
    return &chatService{channelRepo: channelRepo, memberRepo: memberRepo, messageRepo: messageRepo, readRepo: readRepo}
}

func (s *chatService) CreateChannel(ctx context.Context, req dto.CreateChatChannelRequest, createdBy uuid.UUID) (*dto.ChatChannelDTO, error) {
    channel := &models.ChatChannel{
        WorkspaceID: req.WorkspaceID,
        ProjectID:   req.ProjectID,
        Name:        req.Name,
        Type:        req.Type,
        CreatedBy:   createdBy,
    }
    if err := s.channelRepo.Create(ctx, channel); err != nil {
        return nil, err
    }
    members := make([]models.ChatChannelMember, 0, len(req.UserIDs)+1)
    members = append(members, models.ChatChannelMember{ChannelID: channel.ID, UserID: createdBy})
    for _, uid := range req.UserIDs {
        members = append(members, models.ChatChannelMember{ChannelID: channel.ID, UserID: uid})
    }
    if err := s.memberRepo.AddMembers(ctx, members); err != nil {
        return nil, err
    }
    return mapChannelToDTO(channel), nil
}

func (s *chatService) ListWorkspaceChannels(ctx context.Context, workspaceID uuid.UUID) ([]dto.ChatChannelDTO, error) {
    channels, err := s.channelRepo.ListByWorkspace(ctx, workspaceID)
    if err != nil {
        return nil, err
    }
    result := make([]dto.ChatChannelDTO, 0, len(channels))
    for i := range channels {
        result = append(result, *mapChannelToDTO(&channels[i]))
    }
    return result, nil
}

func (s *chatService) AddMembers(ctx context.Context, channelID uuid.UUID, userIDs []uuid.UUID) ([]dto.ChatChannelMemberDTO, error) {
    members := make([]models.ChatChannelMember, 0, len(userIDs))
    for _, uid := range userIDs {
        members = append(members, models.ChatChannelMember{ChannelID: channelID, UserID: uid})
    }
    if err := s.memberRepo.AddMembers(ctx, members); err != nil {
        return nil, err
    }
    stored, err := s.memberRepo.ListMembers(ctx, channelID)
    if err != nil {
        return nil, err
    }
    result := make([]dto.ChatChannelMemberDTO, 0, len(stored))
    for i := range stored {
        result = append(result, dto.ChatChannelMemberDTO{
            ID:        stored[i].ID,
            ChannelID: stored[i].ChannelID,
            UserID:    stored[i].UserID,
            JoinedAt:  stored[i].JoinedAt,
        })
    }
    return result, nil
}

func (s *chatService) RemoveMember(ctx context.Context, channelID, userID uuid.UUID) error {
    return s.memberRepo.RemoveMember(ctx, channelID, userID)
}

func (s *chatService) SendMessage(ctx context.Context, req dto.CreateChatMessageRequest, senderID uuid.UUID) (*dto.ChatMessageDTO, error) {
    message := &models.ChatMessage{
        ChannelID: req.ChannelID,
        SenderID:  senderID,
        Content:   req.Content,
        ReplyToID: req.ReplyToID,
    }
    if err := s.messageRepo.Create(ctx, message); err != nil {
        return nil, err
    }
    return mapMessageToDTO(message), nil
}

func (s *chatService) ListMessages(ctx context.Context, channelID uuid.UUID, limit int) ([]dto.ChatMessageDTO, error) {
    messages, err := s.messageRepo.ListByChannel(ctx, channelID, limit)
    if err != nil {
        return nil, err
    }
    result := make([]dto.ChatMessageDTO, 0, len(messages))
    for i := range messages {
        result = append(result, *mapMessageToDTO(&messages[i]))
    }
    return result, nil
}

func (s *chatService) MarkMessageRead(ctx context.Context, messageID, userID uuid.UUID) error {
    read := &models.ChatMessageRead{
        MessageID: messageID,
        UserID:    userID,
    }
    return s.readRepo.MarkRead(ctx, read)
}

func mapChannelToDTO(channel *models.ChatChannel) *dto.ChatChannelDTO {
    return &dto.ChatChannelDTO{
        ID:          channel.ID,
        WorkspaceID: channel.WorkspaceID,
        ProjectID:   channel.ProjectID,
        Name:        channel.Name,
        Type:        channel.Type,
        CreatedBy:   channel.CreatedBy,
        CreatedAt:   channel.CreatedAt,
    }
}

func mapMessageToDTO(message *models.ChatMessage) *dto.ChatMessageDTO {
    return &dto.ChatMessageDTO{
        ID:        message.ID,
        ChannelID: message.ChannelID,
        SenderID:  message.SenderID,
        Content:   message.Content,
        ReplyToID: message.ReplyToID,
        CreatedAt: message.CreatedAt,
    }
}

package handlers

import (
	"common/module/logger"
	pb "common/module/proto/message_service"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"message/module/application"
	"message/module/infrastructure/api"
)

type MessageHandler struct {
	messageService *application.MessageService
	userService    *application.UserService
	logInfo        *logger.Logger
	logError       *logger.Logger
}

func NewMessageHandler(messageService *application.MessageService, userService *application.UserService, logInfo *logger.Logger, logError *logger.Logger) *MessageHandler {
	return &MessageHandler{messageService: messageService, userService: userService, logInfo: logInfo, logError: logError}
}

func (m MessageHandler) MustEmbedUnimplementedMessageServiceServer() {
}

func (m MessageHandler) GetAllSent(_ context.Context, request *pb.GetRequest) (*pb.GetMultipleResponse, error) {
	sender, err := m.userService.GetByUsername(request.Username)
	if err != nil {
		m.logError.Logger.WithFields(logrus.Fields{
			"userId": request.Username,
		}).Errorf("No user in database")
		return nil, err
	}
	fmt.Println("sender[0].UserId")
	fmt.Println(sender[0].UserId)
	messages, err := m.messageService.GetAllSent(sender[0].UserId)

	fmt.Println(messages)
	response := &pb.GetMultipleResponse{Messages: []*pb.Message{}}
	for _, message := range messages {
		receiver, _ := m.userService.GetById(message.ReceiverId)
		current := api.MapMessageReply(message, receiver[0].Username, sender[0].Username)
		response.Messages = append(response.Messages, current)
	}

	return response, nil
}

func (m MessageHandler) GetAllReceived(_ context.Context, request *pb.GetRequest) (*pb.GetMultipleResponse, error) {
	receiver, err := m.userService.GetByUsername(request.Username)
	if err != nil {
		m.logError.Logger.WithFields(logrus.Fields{
			"userId": request.Username,
		}).Errorf("No user in database")
		return nil, err
	}
	messages, err := m.messageService.GetAllReceived(receiver[0].UserId)
	response := &pb.GetMultipleResponse{Messages: []*pb.Message{}}
	for _, message := range messages {
		sender, _ := m.userService.GetById(message.SenderId)
		current := api.MapMessageReply(message, receiver[0].Username, sender[0].Username)
		response.Messages = append(response.Messages, current)
	}

	return response, nil
}

func (m MessageHandler) SendMessage(_ context.Context, request *pb.SendMessageRequest) (*pb.MessageSentResponse, error) {

	userSender, _ := m.userService.GetByUsername(request.Message.SenderUsername)
	userReceiver, _ := m.userService.GetByUsername(request.Message.ReceiverUsername)

	model := api.MapNewMessage(request.Message.MessageText, userReceiver[0].UserId, userSender[0].UserId)
	message, err := m.messageService.SendMessage(model)

	if err != nil {
		m.logError.Logger.WithFields(logrus.Fields{
			"userId": request.Message.SenderUsername,
		}).Errorf("Can not send message")
		return nil, err
	}

	response := api.MapMessageReply(message, request.Message.ReceiverUsername, request.Message.SenderUsername)
	return &pb.MessageSentResponse{Message: response}, nil
}

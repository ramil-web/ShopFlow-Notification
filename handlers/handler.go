package handlers

import (
	"encoding/json"
	"log"
	"shopflow/notification/services"
)

type HandlerContext struct {
	emailService services.EmailSender
}

func NewHandlerContext(emailService services.EmailSender) *HandlerContext {
	return &HandlerContext{emailService: emailService}
}

type RegistrationEvent struct {
	Email    string `json:"email"`
	Login    string `json:"login"`
	Password string `json:"password"`
	UserID   int    `json:"user_id"`
}

type ApplicationEvent struct {
	Email  string `json:"email"`
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	Text   string `json:"text"`
	File   string `json:"file_url"`
}

// ProcessRegistration принимает только []byte и логирует ошибки
func (h *HandlerContext) ProcessRegistration(msg []byte) {
	var event RegistrationEvent
	err := json.Unmarshal(msg, &event)
	if err != nil {
		log.Println("[ERROR] Failed to unmarshal registration event:", err)
		return
	}

	log.Printf("[INFO] Processing registration for %s\n", event.Email)

	err = h.emailService.SendEmail(event.Email, event.Login, event.Password)
	if err != nil {
		log.Printf("[ERROR] Failed to send email to %s: %v\n", event.Email, err)
		return
	}

	log.Printf("[INFO] Email sent successfully to %s\n", event.Email)
}

func (h *HandlerContext) ProcessApplication(msg []byte) {
	var event ApplicationEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		log.Println("[ERROR] Failed to unmarshal application event:", err)
		return
	}

	log.Printf("[INFO] ProcessApplication: %+v\n", event)

	if event.Email == "" {
		log.Println("[WARN] Email is empty, skipping sending")
		return
	}

	log.Printf("[INFO] Sending application email to %s\n", event.Email)
	if err := h.emailService.SendApplicationEmail(event.Email, event.ID, event.Text, event.File); err != nil {
		log.Printf("[ERROR] Failed to send email to %s: %v\n", event.Email, err)
		return
	}

	log.Printf("[INFO] Application email sent successfully to %s\n", event.Email)
}

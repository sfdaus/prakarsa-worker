package handlers

import (
	"context"
	"prakarsa-app/repository"
	"prakarsa-app/transport/clients"
	"prakarsa-app/worker"
)

type NotifApproved struct{ Client *clients.NotifClient }

func NewNotifApproved(c *clients.NotifClient) *NotifApproved { return &NotifApproved{Client: c} }

func (h *NotifApproved) Name() string     { return "notif.thread_application_approved" } // hanya nama; tidak dipakai filter query karena outbox khusus notif
func (h *NotifApproved) Concurrency() int { return 8 }

func (h *NotifApproved) Handle(ctx context.Context, r repository.Row) error {
	hdr := worker.ParseHeaders(r.HeadersJSON)
	return h.Client.Send(ctx, clients.CreateNotification{
		UserID:        r.UserID,
		Type:          r.Type,
		ReferenceType: r.RefType,
		ReferenceID:   r.RefID,
		Title:         r.Title,
		Message:       r.Message,
		Headers:       hdr,
	})
}

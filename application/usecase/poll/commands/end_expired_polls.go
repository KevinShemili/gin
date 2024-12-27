package commands

import (
	"encoding/json"
	"fmt"
	"gin/application/repository/contracts"
	"gin/application/usecase/poll/results"
	"gin/infrastructure/mail"
	"gin/infrastructure/websocket"
	"log"
	"time"
)

func EndExpiredPolls(UnitOfWork contracts.IUnitOfWork) error {

	polls, err := UnitOfWork.IPollRepository().GetExpiredPolls(time.Now())
	if err != nil {
		return err
	}

	uof, err := UnitOfWork.Begin()
	if err != nil {
		return err
	}
	defer uof.Rollback()

	for _, poll := range polls {
		poll.IsEnded = true
		if err := uof.IPollRepository().Update(poll); err != nil {
			return err
		}
		fmt.Printf("Ended Poll: %d", poll.ID)

	}

	if err := uof.Commit(); err != nil {
		return err
	}

	for _, endedPoll := range polls {
		var broadcastData results.BroadcastExpiry
		broadcastData.BroadcastType = "poll-ended"
		broadcastData.Data.PollID = endedPoll.ID

		message, _ := json.Marshal(broadcastData)
		websocket.BroadcastMessage(string(message))
	}

	go func() {

		for _, poll := range polls {

			results := []map[string]interface{}{}

			for _, category := range poll.Categories {
				results = append(results, map[string]interface{}{
					"CategoryName": category.Name,
					"Votes":        len(category.Votes),
				})
			}

			if err := mail.SendEmail(
				poll.Creator.Email,
				"Poll Has Ended",
				"../../../infrastructure/mail/templates/expired_poll_template.html",
				map[string]interface{}{
					"PollTitle": poll.Title,
					"Results":   results,
				},
			); err != nil {
				log.Printf("Failed to send email. %v", err)
			}
		}
	}()

	return nil
}

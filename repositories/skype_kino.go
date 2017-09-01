package repositories

import (
	"fmt"
	"github.com/andboson/skypeapi"
	"log"
)

func sendFilmsReplyMessage(activity *skypeapi.Activity, location, platform string) {
	name, ok := KinoNamesRu[location]
	if !ok {
		skypeapi.SendReplyMessage(activity, "Не знаю такое место", SkypeToken.AccessToken)
		log.Printf("unknown location: %s", location)
		return
	}
	films := GetMovies(location)

	if platform == "web" {
		skypeapi.SendReplyMessage(activity, "Фильмы в "+name, SkypeToken.AccessToken)
		for _, film := range films {
			var filmText = " \n "
			filmText += fmt.Sprintf("\r\n **%s** ", film.Title)
			filmText += fmt.Sprintf("\r\n %s", film.TimeBlock)
			filmText += fmt.Sprintf("[:](%s)", URL_PREFIX+"/"+film.Img)
			skypeapi.SendReplyMessage(activity, filmText, SkypeToken.AccessToken)
		}
	} else {
		sendReplyMessageRich(activity, "Фильмы в "+name, SkypeToken.AccessToken, films)
	}

}

func sendChoicePlaceReplyMessage(activity *skypeapi.Activity, message, authorizationToken string) error {
	responseActivity := &skypeapi.Activity{
		Type:         activity.Type,
		From:         activity.Recipient,
		Conversation: activity.Conversation,
		Recipient:    activity.From,
		Text:         message,
		InputHint:    "место (lyubava\\plaza)",
		SuggestedActions: skypeapi.SuggestedActions{
			Actions: []skypeapi.CardAction{
				{
					Title: "Любава",
					Type:  "imBack",
					Value: "lyubava",
				},
				{
					Title: "Днепроплаза",
					Type:  "imBack",
					Value: "plaza",
				},
			},
		},
		ReplyToID: activity.ID,
	}
	replyUrl := fmt.Sprintf("%vv3/conversations/%v/activities", activity.ServiceURL, activity.Conversation.ID)

	return skypeapi.SendActivityRequest(responseActivity, replyUrl, authorizationToken)
}

func sendReplyMessageRich(activity *skypeapi.Activity, message, authorizationToken string, films []Film) error {
	var attchmts []skypeapi.Attachment

	for _, film := range films {
		var att = skypeapi.Attachment{
			ContentType: "application/vnd.microsoft.card.hero",
			Content: skypeapi.AttachmentContent{
				Title: film.Title,
				Text:  film.TimeBlock,
				Images: []skypeapi.CardImage{
					{
						URL: URL_PREFIX + "/" + film.Img,
					},
				},
			},
		}

		attchmts = append(attchmts, att)
	}

	responseActivity := &skypeapi.Activity{
		Type:             activity.Type,
		AttachmentLayout: "carousel",
		From:             activity.Recipient,
		Conversation:     activity.Conversation,
		Recipient:        activity.From,
		Text:             message,
		Attachments:      attchmts,
		ReplyToID:        activity.ID,
	}
	replyUrl := fmt.Sprintf("%vv3/conversations/%v/activities", activity.ServiceURL, activity.Conversation.ID)
	return skypeapi.SendActivityRequest(responseActivity, replyUrl, authorizationToken)
}

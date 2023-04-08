package telebot

import (
	"encoding/json"
	"errors"
	"strconv"
)

type MessageThreadID int

// ForumTopic represents a forum topic.
type ForumTopic struct {
	Name            string `json:"name"`
	MessageThreadID int    `json:"message_thread_id"`

	// Color of the topic icon in RGB format.
	//
	// Currently, must be one of 7322096 (0x6FB9F0), 16766590 (0xFFD67E),
	// 13338331 (0xCB86DB), 9367192 (0x8EEE98), 16749490 (0xFF93B2),
	// or 16478047 (0xFB6F5F).
	IconColor int64 `json:"icon_color"`

	// (Optional) Unique identifier of the custom emoji shown as the topic icon.
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
}

// ForumTopicCreated represents a service message about a new
// forum topic created in the chat.
type ForumTopicCreated struct {
	Name string `json:"name"`

	// Color of the topic icon in RGB format.
	//
	// Currently, must be one of 7322096 (0x6FB9F0), 16766590 (0xFFD67E),
	// 13338331 (0xCB86DB), 9367192 (0x8EEE98), 16749490 (0xFF93B2),
	// or 16478047 (0xFB6F5F).
	IconColor int64 `json:"icon_color"`

	// (Optional) Unique identifier of the custom emoji shown as the topic icon.
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
}

// ForumTopicEdited represents a service message about a forum topic
// closed in the chat.
type ForumTopicEdited struct {
	// (Optional) Name of the topic.
	Name string `json:"name,omitempty"`

	// (Optional) Unique identifier of the custom emoji shown as the topic icon.
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
}

// ForumTopicClosed represents a service message about an edited forum topic.
type ForumTopicClosed struct{}

// ForumTopicReopened represents a service message about a forum topic
// reopened in the chat.
type ForumTopicReopened struct{}

// GeneralForumTopicHidden represents a service message about General forum topic
// hidden in the chat. Currently holds no information.
type GeneralForumTopicHidden struct{}

// GeneralForumTopicUnhidden represents a service message about General forum topic
// unhidden in the chat. Currently holds no information.
type GeneralForumTopicUnhidden struct{}

// CreateForumTopic creates a topic in a forum supergroup chat.
// The bot must be an administrator in the chat for this to work and must have
// the CanManageTopics administrator rights. Returns information about the created
// topic as a *ForumTopic object.
func (b *Bot) CreateForumTopic(chat *Chat, ft *ForumTopic) (*ForumTopic, error) {
	if ft == nil {
		return nil, errors.New("telebot: forum topic is nil")
	}
	params := map[string]string{
		"chat_id": chat.Recipient(),
		"name":    ft.Name,
	}

	if ft.IconColor != 0 {
		params["icon_color"] = strconv.FormatInt(ft.IconColor, 10)
	}

	if ft.IconCustomEmojiID != "" {
		params["icon_custom_emoji_id"] = ft.IconCustomEmojiID
	}

	data, err := b.Raw("createForumTopic", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result *ForumTopic
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result, nil
}

// EditForumTopic edits name and icon of a topic in a forum supergroup chat.
// The bot must be an administrator in the chat for this to work and must have
// CanManageTopics administrator rights, unless it is the creator of the topic.
//
// The parameters name and icon are optional. If they are omitted, the existing
// values are kept.
func (b *Bot) EditForumTopic(chat *Chat, msgThreadID int, name string, icon string) error {
	params := map[string]string{
		"chat_id":           chat.Recipient(),
		"message_thread_id": strconv.Itoa(msgThreadID),
	}

	if name != "" {
		params["name"] = name
	}

	if icon != "" {
		params["icon_custom_emoji_id"] = icon
	}

	_, err := b.Raw("editForumTopic", params)
	return err
}

// CloseForumTopic closes an open topic in a forum supergroup chat.
// The bot must be an administrator in the chat for this to work and must have
// CanManageTopics administrator rights, unless it is the creator of the topic.
func (b *Bot) CloseForumTopic(chat *Chat, msgThreadID int) error {
	params := map[string]string{
		"chat_id":           chat.Recipient(),
		"message_thread_id": strconv.Itoa(msgThreadID),
	}

	_, err := b.Raw("closeForumTopic", params)
	return err
}

// ReopenForumTopic reopens a closed topic in a forum supergroup chat.
// The bot must be an administrator in the chat for this to work and must have
// CanManageTopics administrator rights, unless it is the creator of the topic.
func (b *Bot) ReopenForumTopic(chat *Chat, msgThreadID int) error {
	params := map[string]string{
		"chat_id":           chat.Recipient(),
		"message_thread_id": strconv.Itoa(msgThreadID),
	}

	_, err := b.Raw("reopenForumTopic", params)
	return err
}

// DeleteForumTopic deletes a forum topic along with all its messages in
// a forum supergroup chat. The bot must be an administrator in the chat for
// this to work and must have CanManageTopics administrator rights, unless
// it is the creator of the topic.
func (b *Bot) DeleteForumTopic(chat *Chat, msgThreadID int) error {
	params := map[string]string{
		"chat_id":           chat.Recipient(),
		"message_thread_id": strconv.Itoa(msgThreadID),
	}

	_, err := b.Raw("deleteForumTopic", params)
	return err
}

// UnpinAllForumTopicMessages clears the list of pinned messages in a forum topic.
// The bot must be an administrator in the chat for this to work and must have the
// CanPinMessages administrator right in the supergroup.
func (b *Bot) UnpinAllForumTopicMessages(chat *Chat, msgThreadID int) error {
	params := map[string]string{
		"chat_id":           chat.Recipient(),
		"message_thread_id": strconv.Itoa(msgThreadID),
	}

	_, err := b.Raw("unpinAllForumTopicMessages", params)
	return err
}

// GetForumTopicIconStickers returns custom emoji stickers, which can be
// used as a forum topic icon by any user.
func (b *Bot) GetForumTopicIconStickers() ([]Sticker, error) {
	data, err := b.Raw("getForumTopicIconStickers", nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result []Sticker
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result, nil
}

// EditGeneralForumTopic edits the name of the 'General' topic in a forum
// supergroup chat. The bot must be an administrator in the chat for this
// to work and must have CanManageTopics administrator rights.
func (b *Bot) EditGeneralForumTopic(chat *Chat, name string) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
		"name":    name,
	}

	_, err := b.Raw("editGeneralForumTopic", params)
	return err
}

// CloseGeneralForumTopic closes an open 'General' topic in a forum supergroup chat.
// The bot must be an administrator in the chat for this to work and must have the
// CanManageTopics administrator rights.
func (b *Bot) CloseGeneralForumTopic(chat *Chat) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	_, err := b.Raw("closeGeneralForumTopic", params)
	return err
}

// ReopenGeneralForumTopic reopens a closed 'General' topic in a forum supergroup chat.
// The bot must be an administrator in the chat for this to work and must have the
// CanManageTopics administrator rights. The topic will be automatically unhidden
// if it was hidden.
func (b *Bot) ReopenGeneralForumTopic(chat *Chat) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	_, err := b.Raw("reopenGeneralForumTopic", params)
	return err
}

// HideGeneralForumTopic hides the 'General' topic in a forum supergroup chat.
// The bot must be an administrator in the chat for this to work and must have the
// CanManageTopics administrator rights. The topic will be automatically closed
// if it was open.
func (b *Bot) HideGeneralForumTopic(chat *Chat) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	_, err := b.Raw("hideGeneralForumTopic", params)
	return err
}

// UnhideGeneralForumTopic unhides the 'General' topic in a forum supergroup chat.
// The bot must be an administrator in the chat for this to work and must have the
// CanManageTopics administrator rights.
func (b *Bot) UnhideGeneralForumTopic(chat *Chat) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	_, err := b.Raw("unhideGeneralForumTopic", params)
	return err
}

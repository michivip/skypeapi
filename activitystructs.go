/*
MIT License

Copyright (c) 2017 MichiVIP

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package skypeapi

import "time"

type Activity struct {
	// Type of activity. One of these values: contactRelationUpdate, conversationUpdate, deleteUserData,
	// message, ping, typing, endOfConversation. For details about activity types, see Activities overview:
	// 	https://docs.microsoft.com/en-us/bot-framework/rest-api/bot-framework-rest-connector-activities
	Type string `json:"type,omitempty"`
	// The action to apply or that was applied. Use the type property to determine context for the action.
	// For example, if type is contactRelationUpdate, the value of the action property would be add if the
	// user added your bot to their contacts list, or remove if they removed your bot from their contacts list.
	Action string `json:"action,omitempty"`
	// ID that uniquely identifies the activity on the channel.
	ID string `json:"id,omitempty"`
	// Date and time that the message was sent in the UTC time zone, expressed in ISO-8601 format.
	Timestamp time.Time `json:"timestamp,omitempty"`
	// An ID that uniquely identifies the channel. Set by the channel.
	ChannelID string `json:"channelId,omitempty"`
	// URL that specifies the channel's service endpoint. Set by the channel.
	ServiceURL string `json:"serviceUrl,omitempty"`
	// A ChannelAccount object that specifies the sender of the message.
	From ChannelAccount `json:"from,omitempty"`
	// A ConversationAccount object that defines the conversation to which the activity belongs.
	Conversation ConversationAccount `json:"conversation,omitempty"`
	// A ChannelAccount object that specifies the recipient of the message.
	Recipient ChannelAccount `json:"recipient,omitempty"`
	// Array of Attachment objects that defines additional information to include in the message. Each
	// attachment may be either a media file (e.g., audio, video, image, file) or a rich card.
	Attachments []Attachment `json:"attachments,omitempty"`
	// Layout of the rich card attachments that the message includes. One of these values: carousel, list.
	// For more information about rich card attachments, see Add rich card attachments to messages:
	// 	https://docs.microsoft.com/en-us/bot-framework/rest-api/bot-framework-rest-connector-add-rich-cards
	AttachmentLayout string `json:"attachmentLayout,omitempty"`
	// An object that contains channel-specific content. Some channels provide features that require additional
	// information that cannot be represented using the attachment schema. For those cases, set this property to
	// the channel-specific content as defined in the channel's documentation. For more information, see Implement
	// channel-specific functionality:
	// 	https://docs.microsoft.com/en-us/bot-framework/rest-api/bot-framework-rest-connector-channeldata
	ChannelData interface{} `json:"channelData,omitempty"`
	// An ID that uniquely identifies the channel. Set by the channel.
	ChannelId string `json:"channelId,omitempty"`
	// Array of objects that represents the entities that were mentioned in the message. Objects in this array
	// may be any Schema.org object. For example, the array may include Mention objects that identify someone
	// who was mentioned in the conversation and Place objects that identify a place that was mentioned in the
	// conversation.
	Entities []interface{} `json:"entities,omitempty"`
	// Flag that indicates whether or not history is disclosed. Default value is false.
	HistoryDisclosed bool `json:"historyDisclosed,omitempty"`
	// Value that indicates whether your bot is accepting, expecting, or ignoring user input after the message
	// is delivered to the client. One of these values: acceptingInput, expectingInput, ignoringInput.
	InputHint string `json:"inputHint,omitempty"`
	// Locale of the language that should be used to display text within the message, in the format
	// <language>-<country>. The channel uses this property to indicate the user's language, so that your bot
	// may specify display strings in that language. Default value is en-US.
	Locale string `json:"locale,omitempty"`
	// Date and time that the message was sent in the local time zone, expressed in ISO-8601 format.
	LocalTimestamp string `json:"localTimestamp,omitempty"`
	// Array of ChannelAccount objects that represents the list of users that joined the conversation. Present
	// only if activity type is "conversationUpdate" and users joined the conversation.
	MembersAdded []ChannelAccount `json:"membersAdded,omitempty"`
	// Array of ChannelAccount objects that represents the list of users that left the conversation. Present
	// only if activity type is "conversationUpdate" and users left the conversation.
	MembersRemoved []ChannelAccount `json:"membersRemoved,omitempty"`
	// A ConversationReference object that defines a particular point in a conversation.
	RelatesTo ConversationReference `json:"relatesTo,omitempty"`
	// The ID of the message to which this message replies. To reply to a message that the user sent, set this
	// property to the ID of the user's message. Not all channels support threaded replies. In these cases, the
	// channel will ignore this property and use time ordered semantics (timestamp) to append the message to the
	// conversation.
	ReplyToID string `json:"replyToId,omitempty"`
	// Text to be spoken by your bot on a speech-enabled channel. To control various characteristics of your
	// bot's speech such as voice, rate, volume, pronunciation, and pitch, specify this property in
	// Speech Synthesis Markup Language (https://msdn.microsoft.com/en-us/library/hh378377(v=office.14).aspx) format.
	Speak string `json:"speak,omitempty"`
	// A SuggestedActions object that defines the options from which the user can choose.
	SuggestedActions SuggestedActions `json:"suggestedActions,omitempty"`
	// Summary of the information that the message contains. For example, for a message that is sent on an email
	// channel, this property may specify the first 50 characters of the email message.
	Summary string `json:"summary,omitempty"`
	// Text of the message that is sent from user to bot or bot to user. See the channel's documentation for
	// limits imposed upon the contents of this property.
	Text string `json:"text,omitempty"`
	// Format of the message's text. One of these values: markdown, plain, xml. For details about text format,
	// see Create messages: https://docs.microsoft.com/en-us/bot-framework/rest-api/bot-framework-rest-connector-create-messages
	TextFormat string `json:"textFormat,omitempty"`
	// Topic of the conversation to which the activity belongs.
	TopicName string `json:"topicName,omitempty"`
}

type SuggestedActions struct {
	// Array of strings that contains the IDs of the recipients to whom the suggested actions should be displayed.
	To []string `json:"to,omitempty"`
	// Array of CardAction objects that define the suggested actions.
	Actions []CardAction `json:"actions,omitempty"`
}

type CardAction struct {
	// Type of action to perform. For a list of valid values, see Add rich card attachments to messages:
	// 	https://docs.microsoft.com/en-us/bot-framework/rest-api/bot-framework-rest-connector-add-rich-cards
	Type string `json:"type,omitempty"`
	// Text of the button. Only applicable for a button's action.
	Title string `json:"title,omitempty"`
	// URL of an image to display on the button. Only applicable for a button's action.
	Image string `json:"image,omitempty"`
	// Contents of the action. The value of this property will vary according to the action type. For more
	// information, see Add rich card attachments to messages.
	Value string `json:"value,omitempty"`
}

type ConversationReference struct {
	// ID that uniquely identifies the activity that this object references.
	ActivityID string `json:"activityId,omitempty"`
	// A ChannelAccount object that identifies the bot in the conversation that this object references.
	Bot ChannelAccount `json:"bot,omitempty"`
	// An ID that uniquely identifies the channel in the conversation that this object references.
	ChannelID string `json:"channelId,omitempty"`
	// A ConversationAccount object that defines the conversation that this object references.
	Conversation ConversationAccount `json:"conversation,omitempty"`
	// URL that specifies the channel's service endpoint in the conversation that this object references.
	ServiceUrl string `json:"serviceUrl,omitempty"`
	// A ChannelAccount object that identifies the user in the conversation that this object references.
	User ChannelAccount `json:"user,omitempty"`
}

type Attachment struct {
	// The media type of the content in the attachment. For media files, set this property to known media types
	// such as image/png, audio/wav, and video/mp4. For rich cards, set this property to one of these vendor-specific types:
	// 	- application/vnd.microsoft.card.adaptive: A rich card that can contain any combination of text, speech, images, buttons, and input fields. Set the content property to an AdaptiveCard object.
	// 	- application/vnd.microsoft.card.animation: A rich card that plays animation. Set the content property to an AnimationCard object.
	// 	- application/vnd.microsoft.card.audio: A rich card that plays audio files. Set the content property to an AudioCard object.
	// 	- application/vnd.microsoft.card.video: A rich card that plays videos. Set the content property to a VideoCard object.
	// 	- application/vnd.microsoft.card.hero: A Hero card. Set the content property to a HeroCard object.
	// 	- application/vnd.microsoft.card.thumbnail: A Thumbnail card. Set the content property to a ThumbnailCard object.
	// 	- application/vnd.microsoft.com.card.receipt: A Receipt card. Set the content property to a ReceiptCard object.
	// 	- application/vnd.microsoft.com.card.signin: A user Sign In card. Set the content property to a SignInCard object.
	ContentType string `json:"contentType,omitempty"`
	// URL for the content of the attachment. For example, if the attachment is an image, set contentUrl to
	// the URL that represents the location of the image. Supported protocols are: HTTP, HTTPS, File, and Data.
	ContentUrl string `json:"contentUrl,omitempty"`
	// The content of the attachment. If the attachment is a rich card, set this property to the rich card
	// object. This property and the contentUrl property are mutually exclusive.
	Content AttachmentContent `json:"content,omitempty"`
	// Name of the attachment.
	Name string `json:"name,omitempty"`
	// URL to a thumbnail image that the channel can use if it supports using an alternative, smaller form of
	// content or contentUrl. For example, if you set contentType to application/word and set contentUrl to the
	// location of the Word document, you might include a thumbnail image that represents the document. The
	// channel could display the thumbnail image instead of the document. When the user clicks the image, the
	// channel would open the document.
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`
}

type AttachmentContent struct {
	Title string `json:"title,omitempty"`
	// AdaptiveCard
	Type     string       `json:"type,omitempty"`
	Subtitle string       `json:"subtitle,omitempty"`
	Text     string       `json:"text,omitempty"`
	Images   []CardImage  `json:"images,omitempty"`
	Buttons  []CardAction `json:"buttons,omitempty"`
	// slice with TextBlock, Select
	Body    []interface{} `json:"body,omitempty"`
	Actions []Action      `json:"actions,omitempty"`
	Tap     *CardAction   `json:"tap,omitempty"`
}

type CardImage struct {
	URL string      `json:"url"`
	Alt string      `json:"alt,omitempty"`
	Tap *CardAction `json:"tap,omitempty"`
}

type Select struct {
	// Input.ChoiceSet
	Type string `json:"type"`
	ID   string `json:"id,omitempty"`
	// compact
	Style   string         `json:"style,omitempty"`
	Choices []SelectChoice `json:"choices"`
}

type SelectChoice struct {
	Title      string `json:"title"`
	Value      string `json:"value"`
	IsSelected bool   `json:"isSelected,omitempty"`
}

type TextBlock struct {
	// TextBlock
	Type   string `json:"type"`
	Text   string `json:"text"`
	Size   string `json:"size,omitempty"`
	Weight string `json:"weight,omitempty"`
}

type Action struct {
	// Action.Http
	Type string `json:"type,omitempty"`
	// POST
	Method string `json:"method,omitempty"`
	URL    string `json:"url,omitempty"`
	Title  string `json:"title,omitempty"`
}

type ChannelAccount struct {
	// ID that uniquely identifies the bot or user on the channel.
	ID string `json:"id,omitempty"`
	// Name of the bot or user.
	Name string `json:"name,omitempty"`
}

type ConversationAccount struct {
	// The ID that identifies the conversation. The ID is unique per channel. If the channel starts the
	// conversion, it sets this ID; otherwise, the bot sets this property to the ID that it gets back in
	// the response when it starts the conversation (see Starting a conversation).
	ID string `json:"id,omitempty"`
	// Flag to indicate whether or not this is a group conversation. Set to true if this is a group
	// conversation; otherwise, false. The default is false.
	IsGroup bool `json:"isGroup,omitempty"`
	// A display name that can be used to identify the conversation.
	Name string `json:"name,omitempty"`
}

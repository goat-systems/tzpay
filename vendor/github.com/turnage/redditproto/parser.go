// redditproto provides protobuffer definitions and JSON parsing utilities for
// Reddit data types.
//
// A note about JSON parsing utilities: These expect to receive JSON unmodified
// from when it is received as a response from a Reddit endpoint which claims to
// provide the named type. E.g. provide to "ParseComment" exactly the body of
// the response from a call to /by_id/{fullname_of_comment}.json.
package redditproto

import (
	"encoding/json"
	"fmt"
)

// Reddit's responses include a "kind" field that contains a string representing
// the type of the "data" field. These constants are derived from those values.
const (
	listingKind = "Listing"
	commentKind = "t1"
	linkKind    = "t3"
	messageKind = "t4"
)

// ParseComment parses a JSON message expected to be a Comment.
func ParseComment(raw json.RawMessage) (*Comment, error) {
	thing, err := handleThing(raw)
	if err != nil {
		return nil, err
	}

	comment, ok := thing.(*Comment)
	if !ok {
		return nil, fmt.Errorf("JSON message was not a comment")
	}

	return comment, nil
}

// ParseThread parses a JSON message expected to represent a Link comment page.
func ParseThread(raw json.RawMessage) (*Link, error) {
	// The JSON message should be a top level array, holding the link in the
	// first index of a Listing at the first index of the top level array,
	// and holding all the comments in a Listing at the second index of the
	// top level array. I don't know why it's done this way...
	listings := []interface{}{
		&redditThing{},
		&redditThing{},
	}

	if err := json.Unmarshal(raw, &listings); err != nil {
		return nil, err
	}

	if len(listings) != 2 {
		return nil, fmt.Errorf("the top-level JSON message was corrupt")
	}

	rawLink := listings[0].(*redditThing)
	rawComments := listings[1].(*redditThing)

	linkThing, err := parseThing(rawLink)
	if err != nil {
		return nil, err
	}

	linkBuffer, ok := linkThing.(*listingBuffer)
	if !ok {
		return nil, fmt.Errorf("link JSON message was nonlisting")
	}

	if len(linkBuffer.links) != 1 {
		return nil, fmt.Errorf("found an unexpected number of links")
	}

	link := linkBuffer.links[0]

	commentsThing, err := parseThing(rawComments)
	if err != nil {
		return nil, err
	}

	commentsBuffer, ok := commentsThing.(*listingBuffer)
	if !ok {
		return nil, fmt.Errorf("comments JSON message was nonlisting")
	}

	link.Comments = commentsBuffer.comments

	return link, nil
}

// ParseListing parses a JSON message expected to be a listing with any mix of
// Links, Comments, or Messages.
func ParseListing(raw json.RawMessage) (
	[]*Link,
	[]*Comment,
	[]*Message,
	error,
) {
	thing, err := handleThing(raw)
	if err != nil {
		return nil, nil, nil, err
	}

	buffer, ok := thing.(*listingBuffer)
	if !ok {
		return nil, nil, nil, fmt.Errorf("JSON message was nonlisting")
	}

	return buffer.links, buffer.comments, buffer.messages, nil
}

// handleThing unmarshals Things and parses them.
func handleThing(raw json.RawMessage) (interface{}, error) {
	thing, err := unmarshalThing(raw)
	if err != nil {
		return nil, err
	}

	return parseThing(thing)
}

// parseThing parses reddit Things and returns the protobuffer that represents
// the Thing according to its Kind.
func parseThing(thing *redditThing) (interface{}, error) {
	switch thing.Kind {
	case listingKind:
		return unmarshalListing(thing.Data)
	case commentKind:
		return unmarshalComment(thing.Data)
	case linkKind:
		return unmarshalLink(thing.Data)
	case messageKind:
		return unmarshalMessage(thing.Data)
	}
	return nil, fmt.Errorf("Unrecognized message kind")
}

// unmarshalThing unmarshals a JSON message into a redditThing, but leaves the
// Data field as raw JSON so it can be unmarshalled according to the Kind field.
func unmarshalThing(raw json.RawMessage) (*redditThing, error) {
	thing := &redditThing{}
	return thing, json.Unmarshal(raw, thing)
}

// unmarshalListing unmarshals a JSON message into a listing, whose children are
// unmarshaled into redditThings and then parsed. A buffer is returned with a
// slice of comments, links, messages (etc; see the struct definition) contained
// in the listing; generally the caller should know what to expect.
func unmarshalListing(raw json.RawMessage) (*listingBuffer, error) {
	listing := &redditListing{}
	if err := json.Unmarshal(raw, listing); err != nil {
		return nil, err
	}

	buffer := &listingBuffer{}
	for _, childThing := range listing.Children {
		childInterface, err := parseThing(childThing)
		if err != nil {
			return nil, err
		}

		if link, ok := childInterface.(*Link); ok {
			buffer.links = append(buffer.links, link)
		} else if comment, ok := childInterface.(*Comment); ok {
			buffer.comments = append(buffer.comments, comment)
		} else if message, ok := childInterface.(*Message); ok {
			buffer.messages = append(buffer.messages, message)
		} else {
			return nil, fmt.Errorf("corrupted listing child")
		}
	}

	return buffer, nil
}

// unmarshalComment unmarshals a JSON message into a Comment protobuffer, and
// recursively unmarshals the reply tree.
func unmarshalComment(raw json.RawMessage) (*Comment, error) {
	buffer := &commentResponse{}
	if err := json.Unmarshal(raw, buffer); err != nil {
		return nil, err
	}

	replies, _, err := unmarshalReplyTree(buffer.Replies)

	return &Comment{
		ApprovedBy:          buffer.ApprovedBy,
		Author:              buffer.Author,
		AuthorFlairCssClass: buffer.AuthorFlairCssClass,
		AuthorFlairText:     buffer.AuthorFlairText,
		BannedBy:            buffer.BannedBy,
		Body:                buffer.Body,
		BodyHtml:            buffer.BodyHtml,
		Gilded:              buffer.Gilded,
		LinkAuthor:          buffer.LinkAuthor,
		LinkUrl:             buffer.LinkUrl,
		NumReports:          buffer.NumReports,
		ParentId:            buffer.ParentId,
		Replies:             replies,
		Subreddit:           buffer.Subreddit,
		SubredditId:         buffer.SubredditId,
		Distinguished:       buffer.Distinguished,
		Created:             buffer.Created,
		CreatedUtc:          buffer.CreatedUtc,
		Ups:                 buffer.Ups,
		Downs:               buffer.Downs,
		Likes:               buffer.Likes,
		Id:                  buffer.Id,
		Name:                buffer.Name,
	}, err
}

// unmarshalReplyTree unmarshals the reply field of comments. Sometimes this
// field is a listing, sometimes it is a string. This function handles whatever
// it happens to be and returns a slice of the replies.
func unmarshalReplyTree(raw json.RawMessage) ([]*Comment, []*Message, error) {
	repliesThing, err := unmarshalThing(raw)
	if err != nil {
		// When a supply tree is not included, the field is a string,
		// which the JSON unmarshaller chokes on.
		return nil, nil, nil
	}

	bufferInterface, err := parseThing(repliesThing)
	if err != nil {
		return nil, nil, err
	}

	buffer, ok := bufferInterface.(*listingBuffer)
	if !ok {
		return nil,
			nil,
			fmt.Errorf("listing buffer corrupted or mislabeled")
	}

	return buffer.comments, buffer.messages, nil
}

// unmarshalLink unmarshals a JSON message into a Link protobuffer.
func unmarshalLink(raw json.RawMessage) (*Link, error) {
	link := &Link{}
	if err := json.Unmarshal(raw, link); err != nil {
		return nil, err
	}

	return link, nil
}

// unmarshalMessage unmarshals a JSON message into a Message protobuffer.
func unmarshalMessage(raw json.RawMessage) (*Message, error) {
	buffer := &messageResponse{}
	if err := json.Unmarshal(raw, buffer); err != nil {
		return nil, err
	}

	_, replies, err := unmarshalReplyTree(buffer.Replies)
	if err != nil {
		return nil, err
	}

	return &Message{
		Author:           buffer.Author,
		BodyHtml:         buffer.BodyHtml,
		Body:             buffer.Body,
		Context:          buffer.Context,
		FirstMessageName: buffer.FirstMessageName,
		Likes:            buffer.Likes,
		LinkTitle:        buffer.LinkTitle,
		New:              buffer.New,
		ParentId:         buffer.ParentId,
		Subject:          buffer.Subject,
		Subreddit:        buffer.Subreddit,
		WasComment:       buffer.WasComment,
		Created:          buffer.Created,
		CreatedUtc:       buffer.CreatedUtc,
		Id:               buffer.Id,
		Name:             buffer.Name,
		Messages:         replies,
	}, nil
}

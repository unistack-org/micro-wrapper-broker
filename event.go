package broker

import (
	"strings"

	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/metadata"
)

type event struct {
	event broker.Event
}

func (evt *event) Topic() string {
	return evt.event.Topic()
}

func (evt *event) Payload() interface{} {
	return nil
}

func (evt *event) ContentType() string {
	return "raw"
}

func (evt *event) AddHeader(md metadata.Metadata) {
	if evt.event.Message().Header == nil {
		evt.event.Message().Header = make(metadata.Metadata)
	}
	for k, v := range md {
		evt.event.Message().Header[strings.Title(k)] = v
	}
}

func (evt *event) Header() map[string]string {
	return evt.event.Message().Header
}

func (evt *event) Body() []byte {
	return evt.event.Message().Body
}

func (evt *event) Codec() codec.Reader {
	return nil
}

func (evt *event) Event() broker.Event {
	return evt.event
}

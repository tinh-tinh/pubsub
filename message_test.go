package pubsub_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/pubsub/v2"
)

func Test_NewMessage(t *testing.T) {
	m := pubsub.NewMessage("topic", "content")
	require.Equal(t, "topic", m.GetTopic())
	require.Equal(t, "content", m.GetContent())
}

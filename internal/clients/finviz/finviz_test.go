package finviz_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jnericks/obibot/internal/clients/finviz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHeatMap(t *testing.T) {
	resp, err := finviz.GetHeatMap(context.Background())
	require.NoError(t, err)

	assert.NotEmpty(t, resp.URL)
	fmt.Println(resp.URL)
}

package flamboyant_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
)

func TestOn(t *testing.T) {
	assert.False(t, flamboyant.On(context.Background()), "sin banderas, el diseño de siempre")
	assert.True(t, flamboyant.On(feature.Enable(context.Background(), flamboyant.Flag)))
	assert.True(t, flamboyant.On(feature.With(context.Background(), feature.Flags{Compare: true})), "Comparar pinta el shell flamboyant")
	assert.False(t, flamboyant.On(feature.Enable(context.Background(), "otra")))
}

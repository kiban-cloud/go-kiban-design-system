package flamboyant_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
)

func TestConNombreNuevo(t *testing.T) {
	casos := []struct {
		path, want string
		ok         bool
	}{
		{"/workfloo", "/solicitudes", true},
		{"/workfloo/workfloos/123?tab=2", "/solicitudes/workfloos/123?tab=2", true},
		{"/workfloo?x=1", "/solicitudes?x=1", true},
		{"/rekon/payments", "/pagos/payments", true},
		{"/link/vault", "/fuentes/vault", true},
		{"/klin/deliveries", "/cumplimiento/deliveries", true},
		{"/klin-internal/tools/clean-data", "/cumplimiento-interno/tools/clean-data", true},
		{"/klinx", "/klinx", false},
		{"/workfloo-templates", "/workfloo-templates", false},
		{"/crm/customers", "/crm/customers", false},
		{"/kiban-cloud/profile", "/kiban-cloud/profile", false},
		{"/portal/workfloo/1", "/portal/workfloo/1", false},
		{"", "", false},
	}
	for _, c := range casos {
		got, ok := flamboyant.ConNombreNuevo(c.path)
		assert.Equal(t, c.want, got, c.path)
		assert.Equal(t, c.ok, ok, c.path)
	}
}

func TestConNombreViejo(t *testing.T) {
	casos := []struct {
		path, want string
		ok         bool
	}{
		{"/solicitudes/workfloos/123", "/workfloo/workfloos/123", true},
		{"/pagos", "/rekon", true},
		{"/fuentes/test-cases?page=2", "/link/test-cases?page=2", true},
		{"/cumplimiento/deliveries", "/klin/deliveries", true},
		{"/cumplimiento-interno/deliveries/9", "/klin-internal/deliveries/9", true},
		{"/pagosx", "/pagosx", false},
		{"/workfloo/workfloos", "/workfloo/workfloos", false},
	}
	for _, c := range casos {
		got, ok := flamboyant.ConNombreViejo(c.path)
		assert.Equal(t, c.want, got, c.path)
		assert.Equal(t, c.ok, ok, c.path)
	}
}

func TestRuta(t *testing.T) {
	apagado := context.Background()
	prendido := feature.Enable(context.Background(), flamboyant.Flag)

	assert.Equal(t, "/workfloo/workfloos", flamboyant.Ruta(apagado, "/workfloo/workfloos"), "sin flamboyant, el nombre de siempre")
	assert.Equal(t, "/solicitudes/workfloos", flamboyant.Ruta(prendido, "/workfloo/workfloos"))
	assert.Equal(t, "/crm/customers", flamboyant.Ruta(prendido, "/crm/customers"), "crm no cambia")
	assert.Equal(t, "", flamboyant.Ruta(prendido, ""))
}

func TestConRutas(t *testing.T) {
	ctx := context.Background()
	assert.False(t, flamboyant.AceptaNombresNuevos(ctx))
	assert.False(t, flamboyant.PorNombreNuevo(ctx))

	viejo := flamboyant.ConRutas(ctx, false)
	assert.True(t, flamboyant.AceptaNombresNuevos(viejo))
	assert.False(t, flamboyant.PorNombreNuevo(viejo))

	nuevo := flamboyant.ConRutas(ctx, true)
	assert.True(t, flamboyant.AceptaNombresNuevos(nuevo))
	assert.True(t, flamboyant.PorNombreNuevo(nuevo))
}

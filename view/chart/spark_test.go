package chart_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/chart"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, c.Render(context.Background(), &buf))
	return buf.String()
}

// La serie va de izquierda a derecha con 4 px de margen, el mínimo abajo y
// el máximo arriba; el último valor (hoy) lleva el punto.
func TestSpark_TrazaLaSerie(t *testing.T) {
	body := render(t, chart.Spark(chart.SparkOptions{Values: []float64{1, 2, 3}, Color: "#8000FF", Width: 100, Height: 40, Class: "h-14"}))

	assert.Contains(t, body, "data-kiban-spark")
	assert.Contains(t, body, `class="block w-full h-14"`)
	assert.Contains(t, body, `viewBox="0 0 100 40"`)
	assert.Contains(t, body, `d="M4 34 L50 22 L96 10"`, "la línea")
	assert.Contains(t, body, `d="M4 34 L50 22 L96 10 L96 40 L4 40 Z"`, "el área cierra contra el piso")
	assert.Contains(t, body, `cx="96" cy="10"`, "el punto en el último valor")
	assert.Contains(t, body, `fill="#8000FF"`)
	assert.Contains(t, body, `stroke="#8000FF"`)
}

// Una serie plana queda abajo; un solo valor se traza como línea plana.
func TestSpark_PlanaYUnSoloValor(t *testing.T) {
	plana := render(t, chart.Spark(chart.SparkOptions{Values: []float64{5, 5}, Width: 100, Height: 40}))
	assert.Contains(t, plana, `d="M4 34 L96 34"`)

	uno := render(t, chart.Spark(chart.SparkOptions{Values: []float64{7}, Width: 100, Height: 40}))
	assert.Contains(t, uno, `d="M4 34 L96 34"`)
	assert.Contains(t, uno, `cx="96" cy="34"`)
}

// Sin valores, o con un lienzo demasiado chico, no se dibuja nada.
func TestSpark_SinSerieOLienzoChicoNoPinta(t *testing.T) {
	assert.Empty(t, render(t, chart.Spark(chart.SparkOptions{Width: 100, Height: 40})))
	assert.Empty(t, render(t, chart.Spark(chart.SparkOptions{Values: []float64{1, 2}, Width: 8, Height: 40})))
	assert.Empty(t, render(t, chart.Spark(chart.SparkOptions{Values: []float64{1, 2}, Width: 100, Height: 16})))
}

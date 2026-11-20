package chart_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/chart"
)

func TestStack_TramosProporcionalesYSinCeros(t *testing.T) {
	body := render(t, chart.Stack(chart.StackOptions{Parts: []chart.StackPart{
		{Label: "Exitosas", Value: 3, Display: "3", Color: "#00A884"},
		{Label: "En curso", Value: 0, Color: "#999"},
		{Label: "Con error", Value: 1, Display: "1", Color: "#E5484D"},
	}}))
	assert.Contains(t, body, "width:75.00%;background:#00A884")
	assert.Contains(t, body, "width:25.00%;background:#E5484D")
	assert.NotContains(t, body, "width:0;background:#999", "un tramo en cero no se pinta")
	assert.Contains(t, body, "En curso", "pero sigue en la leyenda")

	assert.NotContains(t, render(t, chart.Stack(chart.StackOptions{HideLegend: true, Parts: []chart.StackPart{{Label: "A", Value: 1}}})), ">A<")
}

func TestRank_ProporcionalAlMayor(t *testing.T) {
	body := render(t, chart.Rank(chart.RankOptions{Color: "#0047FF", Rows: []chart.RankRow{
		{Name: "Alta", Value: 10, Display: "10", Share: "67 %"},
		{Name: "Baja", Value: 5, Display: "5"},
	}}))
	assert.Contains(t, body, "width:100.00%")
	assert.Contains(t, body, "width:50.00%")
	assert.Contains(t, body, "67 %")
}

func TestCompare_Diferencia(t *testing.T) {
	o := chart.CompareOptions{Delta: true,
		Series: [2]chart.StackPart{{Label: "Original", Color: "#999"}, {Label: "Versión", Color: "#0047FF"}},
		Rows: []chart.CompareRow{
			{Label: "Aprobado", A: 100, B: 166, ADisplay: "100", BDisplay: "166"},
			{Label: "Rechazado", A: 50, B: 8, ADisplay: "50", BDisplay: "8"},
			{Label: "Igual", A: 3, B: 3, ADisplay: "3", BDisplay: "3"},
		}}
	body := render(t, chart.Compare(o))
	assert.Contains(t, body, ">+66<")
	assert.Contains(t, body, ">−42<")
	assert.Contains(t, body, "bg-[#FDECEC]", "baja en rojo")
	assert.Contains(t, body, "width:100.00%", "la barra mayor llena el ancho")
	o.Delta = false
	assert.False(t, strings.Contains(render(t, chart.Compare(o)), "+66"))
}

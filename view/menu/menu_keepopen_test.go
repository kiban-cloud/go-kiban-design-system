package menu_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/kiban-cloud/go-kiban-design-system/view/menu"
)

func renderMenu(t *testing.T, cfg menu.Config) string {
	t.Helper()
	var buf bytes.Buffer
	if err := menu.Menu(cfg).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// Por default el item se cierra solo tras el click — comportamiento previo.
func TestMenuItem_ClosesOnClickByDefault(t *testing.T) {
	body := renderMenu(t, menu.Config{
		ID:    "row-menu",
		Items: []menu.MenuItem{{Label: "Editar", OnClick: "irAEditar()"}},
	})

	if !strings.Contains(body, "irAEditar();window.kibanCloseMenu(&#39;row-menu&#39;)") {
		t.Errorf("expected the close call appended after the user expression: %s", body)
	}
}

// Un item sin OnClick igual se cierra: el click sólo sirve para dismissar.
func TestMenuItem_ClosesEvenWithoutOnClick(t *testing.T) {
	body := renderMenu(t, menu.Config{
		ID:    "row-menu",
		Items: []menu.MenuItem{{Label: "Cerrar"}},
	})

	if !strings.Contains(body, "window.kibanCloseMenu(&#39;row-menu&#39;)") {
		t.Errorf("expected the close call: %s", body)
	}
}

// KeepOpen deja el panel abierto para que el item pueda mostrar su
// feedback dentro (un check tras copiar, por ejemplo). El cierre queda a
// cargo del item.
func TestMenuItem_KeepOpenSuppressesAutoClose(t *testing.T) {
	body := renderMenu(t, menu.Config{
		ID: "row-menu",
		Items: []menu.MenuItem{
			{Label: "Copiar ID", OnClick: "copiar(this)", KeepOpen: true},
		},
	})

	if !strings.Contains(body, `onclick="copiar(this)"`) {
		t.Errorf("expected the user expression untouched: %s", body)
	}
	if strings.Contains(body, "kibanCloseMenu") {
		t.Errorf("KeepOpen must not append the close call: %s", body)
	}
}

// KeepOpen sin OnClick no emite onclick alguno: sin él, el click del item
// no haría nada y el atributo vacío sólo sería ruido.
func TestMenuItem_KeepOpenWithoutOnClickEmitsNoHandler(t *testing.T) {
	body := renderMenu(t, menu.Config{
		ID:    "row-menu",
		Items: []menu.MenuItem{{Label: "Inerte", KeepOpen: true}},
	})

	// El trigger del kebab SIEMPRE lleva su onclick (kibanToggleMenu);
	// lo que no debe llevar handler es el item.
	item := body[strings.Index(body, `role="menuitem"`):]
	if strings.Contains(item, "onclick") {
		t.Errorf("expected no onclick on the item: %s", item)
	}
}

// KeepOpen no afecta a los demás items del mismo menú.
func TestMenuItem_KeepOpenIsPerItem(t *testing.T) {
	body := renderMenu(t, menu.Config{
		ID: "row-menu",
		Items: []menu.MenuItem{
			{Label: "Copiar ID", OnClick: "copiar(this)", KeepOpen: true},
			{Label: "Eliminar", OnClick: "borrar()"},
		},
	})

	if !strings.Contains(body, "borrar();window.kibanCloseMenu(&#39;row-menu&#39;)") {
		t.Errorf("the sibling item must still auto-close: %s", body)
	}
}

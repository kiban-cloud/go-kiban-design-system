package layout

import (
	"context"
	"testing"

	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
)

func tieneModulo(ms []fsModulo, key string) (fsModulo, bool) {
	for _, m := range ms {
		if m.Key == key {
			return m, true
		}
	}
	return fsModulo{}, false
}

func TestFsTools_CarteraOcultaSinBandera(t *testing.T) {
	ms, _ := fsTools(context.Background(), Config{})
	if _, ok := tieneModulo(ms, "core"); ok {
		t.Fatal("sin FlagCartera, Cartera no sale en el menú")
	}
}

func TestFsTools_CarteraConBandera(t *testing.T) {
	ctx := feature.Enable(context.Background(), FlagCartera)
	ms, _ := fsTools(ctx, Config{})
	m, ok := tieneModulo(ms, "core")
	if !ok {
		t.Fatal("con FlagCartera, Cartera sale en el menú")
	}
	if m.Href != CarteraHref {
		t.Fatalf("entrada de Cartera: %q, esperaba %q", m.Href, CarteraHref)
	}
	if _, ok := tieneModulo(ms, "folo"); ok {
		t.Fatal("FlagCartera no prende Monitoreo")
	}
}

// Package flamboyant dice si un request se pinta con el diseño «kiban cloud
// flamboyant». Lo consultan el shell (view/layout) y los componentes que
// tienen variante flamboyant (botón, checkbox, badge…): con la bandera
// apagada cada uno se ve como siempre.
//
// Vive aparte de view/layout porque layout importa componentes (avatar, por
// ejemplo) y esos componentes necesitan esta pregunta: ponerla en layout
// haría un ciclo.
package flamboyant

import (
	"context"

	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
)

// Flag es la bandera del diseño. La prende middleware/authcookie (feature
// flag de la compañía y elección del usuario) o, en kiban-proto, el selector
// de fase.
const Flag = "shell-flamboyant"

// On dice si el request se pinta flamboyant. En modo Comparar también: el
// shell envuelve la pantalla entera y no se pueden mostrar los dos diseños
// lado a lado.
func On(ctx context.Context) bool {
	return feature.Comparing(ctx) || feature.Enabled(ctx, Flag)
}

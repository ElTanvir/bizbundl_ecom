package hero

import (
	"bizbundl/pkg/components/registry"

	"github.com/a-h/templ"
)

func init() {
	registry.Register(&registry.Component{
		Type: "hero",
		Renderer: func(props map[string]interface{}) templ.Component {
			return View(mapProps(props))
		},
	})
}

func mapProps(props map[string]interface{}) Props {
	return Props{
		Title:           getString(props, "Title"),
		Subtitle:        getString(props, "Subtitle"),
		ButtonText:      getString(props, "ButtonText"),
		ButtonLink:      getString(props, "ButtonLink"),
		BackgroundImage: getString(props, "BackgroundImage"),
		Align:           getString(props, "Align"),
		Variant:         getString(props, "Variant"),
	}
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

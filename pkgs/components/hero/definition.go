package hero

import (
	"bizbundl/pkgs/components/registry"
	"bizbundl/pkgs/components/utils"

	"github.com/a-h/templ"
)

func init() {
	registry.Register(&registry.Component{
		Type:        "hero",
		Title:       "Hero Section",
		Description: "A large banner area usually placed at the top of the page.",
		Category:    "Content",
		Variants: map[string]registry.VariantDefinition{
			"standard": {
				Name:        "Standard Image",
				Description: "Title, Subtitle, and Background Image",
				Props: map[string]registry.PropDefinition{
					"Title":           {Type: registry.TypeString, Required: true, Default: "Welcome"},
					"Subtitle":        {Type: registry.TypeString},
					"BackgroundImage": {Type: registry.TypeImage, Description: "URL for background"},
					"Align":           {Type: registry.TypeString, Default: "center"},
				},
			},
			"video": {
				Name:        "Background Video",
				Description: "Hero with looped background video",
				Props: map[string]registry.PropDefinition{
					"Title":    {Type: registry.TypeString, Required: true},
					"VideoURL": {Type: registry.TypeString, Required: true, Description: "MP4 URL"},
					"Overlay":  {Type: registry.TypeBoolean, Default: true},
				},
			},
		},
		Renderer: func(props map[string]interface{}) templ.Component {
			return View(mapProps(props))
		},
	})
}

func mapProps(props map[string]interface{}) Props {
	return Props{
		Title:           utils.GetString(props, "Title"),
		Subtitle:        utils.GetString(props, "Subtitle"),
		ButtonText:      utils.GetString(props, "ButtonText"),
		ButtonLink:      utils.GetString(props, "ButtonLink"),
		BackgroundImage: utils.GetString(props, "BackgroundImage"),
		Align:           utils.GetString(props, "Align"),
		Variant:         utils.GetString(props, "Variant"),
	}
}
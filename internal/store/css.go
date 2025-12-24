package store

import "fmt"

const cssVariableKey = "css:variables"

func SetCssVariable(data string) {
	s := Get()
	formattedData := fmt.Sprintf("<style>%s</style>", data)
	s.Set(cssVariableKey, formattedData)
}

func GetCssVariable() string {
	s := Get()
	v, ok := s.Get(cssVariableKey)
	if !ok {
		return ""
	}
	u, _ := v.(string)
	return u
}

var CriticalCSS = `
<style>
      :root {
        /* Fonts */
        /* Mapped to Design System: 'Inter' (sans) and 'Nunito' (display) */
        --font-body: "Inter", sans-serif;
        --font-title: "Nunito", sans-serif;

        /* Light Theme Colors (Mapped to "Playful Doodle" Palette) */
        --color-surface: #ffffff; /* brand-bg */
        --color-surface-alt: #fbf5ff; /* brand-bg-start (light purple) */
        --color-on-surface: #555555; /* brand-text-secondary */
        --color-on-surface-strong: #212121; /* brand-text-primary */

        --color-primary: #7f00ff; /* brand-primary (Purple) */
        --color-on-primary: #ffffff; /* Text on primary */

        --color-secondary: #ccff00; /* brand-accent-lime */
        --color-on-secondary: #212121; /* Text on lime */

        --color-outline: #d4d4d4; /* Default light border (e.g., header) */
        --color-outline-strong: #212121; /* Strong border (e.g., service cards) */

        --color-info: #0284c7; /* Unused in this design, kept as-is */
        --color-on-info: #000000;

        --color-success: #059669; /* Unused in this design, kept as-is */
        --color-on-success: #000000;

        --color-warning: #ff4500; /* brand-accent-orange */
        --color-on-warning: #ffffff; /* Text on orange */

        --color-danger: #fff5f5; /* Mapped to brand-bg-end (light red) */
        --color-on-danger: #ef4444; /* Original danger red, for text use */

        /* Border Radius */
        --radius-none: 0;
        --radius-radius: 1rem; /* Mapped to borderRadius-brand */

        /* Shadows */
        --shadow-sm: 0 1px 2px 0 rgb(0 0 0 / 0.05);
        --shadow-md: 0 4px 10px rgba(0, 0, 0, 0.05); /* Mapped to boxShadow-brand */
        --shadow-lg: 0 10px 15px -3px rgb(0 0 0 / 0.1),
          0 4px 6px -4px rgb(0 0 0 / 0.1);
        --shadow-xl: 0 20px 25px -5px rgb(0 0 0 / 0.1),
          0 8px 10px -6px rgb(0 0 0 / 0.1);
      }
    </style>
  `

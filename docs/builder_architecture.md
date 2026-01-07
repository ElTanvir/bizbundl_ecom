# Architecture: Page Builder & Resolver V2

**Goal:** Allow admins to build pages from pre-compiled Atomic Components (`pkg/components`) with high performance and strict type safety.

## 1. Core Concepts

### A. Atomic Component System
We use a **Modular Monolith** approach where each UI block is a self-contained package in `pkg/components/`.

**Structure:**
```text
pkg/components/
├── registry/               # The "Brain". Maps keys ("product_grid") to implementations.
├── product_grid/           # The "Component".
│   ├── definition.go       # Registration & Main Dispatcher Logic.
│   ├── resolver.go         # The "Smart Dispatcher" (Resolver V2).
│   └── variants/           # Isolated sub-types.
│       ├── grid/           # Variant A
│       │   ├── definition.go # Variant Definition
│       │   └── view.templ    # Unique UI
│       └── carousel/       # Variant B
```

### B. The Registry (The Truth Source)
The `pkgs/components/registry` package holds the map of all available components.
*   **Safety:** It panics on duplicate registration (Safe Start).
*   **Type:** `map[string]*Component`.

---

## 2. The "Resolver V2" Flow (Dispatcher Pattern)

Old systems often used a giant `switch` statement or a single resolver for all variants. **Resolver V2** uses a **Dispatcher Pattern**:

### The Flow:
1.  **Request:** User requests a page with a `product_grid` section, prop `Variant: "carousel"`.
2.  **Global Resolver (`product_grid/resolver.go`):**
    *   The Engine calls the Component's Main Resolver.
    *   It looks at `props["Variant"]`.
    *   It checks the Registry: `Does component.Variants["carousel"] have a custom Resolver?`
3.  **Variant Dispatch:**
    *   **YES:** It calls `carousel.Resolver.Resolve(ctx, props)`.
        *   *Result:* Fetches "New Arrivals" (specific to carousel logic).
    *   **NO:** It falls back to the Default Resolver (`grid`).
4.  **Result:** The specific data needed for that variant is returned.

**Benefit:**
*   **Isolation:** A "Masonry" variant can have totally different data logic (infinite scroll) than a "Slider" variant (static list), without touching the core code.

---

## 3. How to Add New Things

### A. How to Add a New **Variant** (e.g., "Masonry" for Product Grid)
1.  **Create Folder:** `pkgs/components/product_grid/variants/masonry/`
2.  **Create View:** `view.templ`
    ```go
    package masonry
    templ View(props map[string]any) { ... }
    ```
3.  **Create Definition:** `definition.go`
    ```go
    package masonry
    func Definition(svc *service.Catalog) registry.VariantDefinition {
        return registry.VariantDefinition{
            Name: "masonry",
            Resolver: &resolver{...}, // Optional: Custom Data Logic
            Renderer: func(p) { return View(p) },
        }
    }
    ```
4.  **Register It:** Go to `pkgs/components/product_grid/definition.go` and add:
    ```go
    c.Variants["masonry"] = masonry.Definition(catalogSvc)
    ```
    *Done! It is now available in the Page Builder.*

### B. How to Add a New **Component** (e.g., "Testimonial")
1.  **Create Folder:** `pkgs/components/testimonial/`
2.  **Define Structure:**
    *   `definition.go`: Define the `Component` struct (Type, Title, Props).
    *   `view.templ`: The UI.
3.  **Register It:**
    *   Go to `cmd/server/main.go` (or your component init module).
    *   Call `testimonial.Register()`.

---

## 4. Usage in Database
The DB stores simple JSON. The Resolvers hydrate it.
```json
{
  "type": "product_grid",
  "props": {
    "Variant": "carousel",
    "Limit": 8,
    "Title": "Hot Deals"
  }
}
```

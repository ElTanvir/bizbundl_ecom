# Architecture: High-Performance "Strict" Component Builder

**Goal:** Allow admins to build pages from pre-compiled components without sacrificing performance or safety. "Little Freedom, High Reliability."

## 1. Core Philosophy
*   **Compile Over Interpret:** We do not parse HTML/Templates from the database at runtime. All HTML structures are compiled binaries (`.templ`).
*   **Structure Over Style:** Admins configure *Data Sources* (e.g., "Category: Electronics"), not *CSS Styles* (e.g., "Padding: 10px").
*   **Fail Early:** Configuration is validated at "Save Time", ensuring broken pages never reach the database.

---

## 2. Technical Components

### A. The Schema (Database)
We store the "Layout Configuration" as a JSONB array in the `pages` table.

**Table:** `pages`
| Column | Type | Description |
| :--- | :--- | :--- |
| `id` | `UUID` | PK |
| `route` | `string` | e.g. `/`, `/landing/black-friday` |
| `sections` | `JSONB` | Array of Section Configs |

**JSON Structure:**
```json
[
  {
    "id": "hero_v1",
    "props": {
       "title": "Summer Sale",
       "bg_image": "/uploads/summer.jpg",
       "cta_link": "/collections/summer"
    }
  },
  {
    "id": "product_carousel",
    "variant": "grid_compact",
    "data_source": {
       "type": "category",
       "value": "electronics",
       "limit": 8
    }
  }
]
```

### B. The Registry (Go Interface)
The core of the system is the `Component` interface. Every "Block" (Hero, Carousel, Text) must implement this.

```go
type Component interface {
    // ID returns the unique string identifier (e.g., "hero_v1")
    ID() string

    // Validate checks if the raw JSON config is valid.
    // e.g., Checks if "category_id" actually exists in DB.
    // Runs ONLY when Admin clicks "Save".
    Validate(ctx context.Context, config json.RawMessage) error

    // FetchData retrieves necessary data for the component.
    // Runs concurrently at Request time.
    FetchData(ctx context.Context, config json.RawMessage) (any, error)

    // Render produces the HTML.
    // Uses the data returned by FetchData.
    Render(data any) templ.Component
}
```

### C. The Request Lifecycle (Runtime Flow)
1.  **Request:** User hits `/`.
2.  **Lookup:** Handler fetches `sections` JSON from DB.
3.  **Hydration (The "Waterfall" Killer):**
    *   Handler iterates over sections.
    *   Launches a **Goroutine** for each component's `FetchData()`.
    *   Uses `errgroup` to wait for all data concurrently.
4.  **Rendering:**
    *   Once data is ready, `base_layout.templ` iterates through sections.
    *   Calls `registry.Get(id).Render(data)`.
5.  **Response:** HTML stream sent to user.

---

## 3. Safety Mechanisms
1.  **Validation Hook:** The API endpoint `POST /admin/pages`:
    *   Decodes JSON.
    *   Loops through sections -> calls `component.Validate()`.
    *   **If any fail:** Returns 400 Error ("Category 'X' does not exist").
    *   **Result:** Impossible to save a broken page configuration.
2.  **Graceful Degrade:** If `FetchData` fails at runtime (e.g., DB glitch), the component can return a "Empty" state or be skipped entirely, keeping the rest of the page alive.

## 4. Admin UI (The Experience)
*   **Implementation:** HTMX + JSON Form.
*   **Sidebar:** List of available components (fetched from Registry).
*   **Main Area:** Live Preview (iframe) or Block List.
*   **Editor:** When a block is clicked, a specific Form (defined by the Component) appears to edit `props`.

# Middleware & Caching Architecture Report
**Context:** Go (Fiber) App + Cloudflare Edge + Limited VPS Budget.

---

## 1. The Strategy: "Edge First, Origin Safe"
Since you are using Cloudflare, your Go server (the "Origin") should almost **never** serve full HTML content. Ideally, 95%+ of traffic stops at Cloudflare.

### The Flow
1.  **User** requests `your-shop.com`.
2.  **Cloudflare (Edge)** checks if it has the page.
    *   *Hit:* Serves instantly. Your VPS sleeps.
    *   *Miss:* Asks your Go VPS.
3.  **Go VPS (Origin)**:
    *   Checks **Fiber RAM Cache** (Level 2).
    *   *Hit:* Returns 200 OK (from RAM). DB sleeps.
    *   *Miss:* Renders template, Queries DB, Returns 200 OK.

---

## 2. Final Recommendations for Your Code

### A. Use the Built-in `etag` Middleware (Mandatory)
*   **Why:** Cloudflare relies on ETags to ask "Did this change?".
*   **How it works:** Cloudflare sends `If-None-Match: "xyz"`. If your Go app sees nothing changed, it sends `304 Not Modified` (0KB body).
*   **Benefit:** Saves your bandwidth and speeds up Cloudflare re-validation.
*   **Action:**
    ```go
    import "github.com/gofiber/fiber/v2/middleware/etag"
    app.Use(etag.New()) // Put this near the top
    ```

### B. Use the Built-in `cache` Middleware (Recommended)
*   **Why:** Even with Cloudflare, you need a "shield" for your database. If Cloudflare clears its cache (e.g., new deployment), 1000 users might hit your server at once. Memory cache answers them without waking up the Database.
*   **Config:** Set a short expiration (e.g., 1 minute). Cloudflare handles long-term storage; this is just for burst protection.
*   **Action:**
    ```go
    import "github.com/gofiber/fiber/v2/middleware/cache"
    
    app.Use(cache.New(cache.Config{
        Expiration:   1 * time.Minute,
        CacheControl: true, // Tells Cloudflare "Yes, you can cache this"
    }))
    ```

### C. Remove Custom `RouteCacheMiddleware`
*   **Why:** As analyzed previously, your custom implementation isn't thread-safe for memory limits and creates a Denial of Service risk.
*   **Action:** Delete `internal/middleware/route_cache.go` and `internal/store/route_cache.go`.

---

## 3. Controlling Cloudflare (The Secret Sauce)
You control Cloudflare using the `Cache-Control` header sent from Go.

| Header Value | Meaning | Recommended For |
| :--- | :--- | :--- |
| `public, max-age=60` | "Cache this for 60 seconds anywhere" | **Dynamic Pages** (Home, Product) |
| `public, max-age=31536000, immutable` | "Cache forever, it never changes" | **Statics** (Images, JS, CSS) |
| `private, no-cache` | "Never cache this shared" | **Cart, Checkout, Admin** |

### Implementation in Handlers
```go
// Homepage (Cached)
app.Get("/", func(c *fiber.Ctx) error {
    c.Set("Cache-Control", "public, max-age=60") // Cloudflare holds for 1 min
    return c.Render("home", fiber.Map{})
})

// Checkout (Never Cached)
app.Get("/checkout", func(c *fiber.Ctx) error {
    c.Set("Cache-Control", "private, no-cache")
    return c.Render("checkout", fiber.Map{})
})
```

---

## 4. Summary Checklist
1.  [ ] **Delete** your custom middleware files.
2.  [ ] **Add** `app.Use(etag.New())`.
3.  [ ] **Add** `app.Use(cache.New())` (for DB protection).
4.  [ ] **Audit** your handlers to send correct `Cache-Control` headers.

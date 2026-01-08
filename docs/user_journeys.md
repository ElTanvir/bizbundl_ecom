# User Journeys

## Journey 1: The Aspiring Mogul (Platform User)
**Persona**: "Serious Sam" wants to start 3 different niche brands.

1.  **Acquisition**: Sam lands on `bizbundl.com` (Custom UI). Sees "Start your Empire".
2.  **Sign Up**: Signs up with email/password. (Global User).
3.  **Onboarding**: Redirected to `/dashboard`. It's empty.
4.  **Creation**: Clicks "Create New Shop".
    *   Enters: "Neon Vibes" (Name), `neon-vibes` (Subdomain).
    *   Selects Plan: "Starter - $29/mo".
    *   Enters Payment (or Trial).
5.  **Provisioning**: System creates `shop_neon` schema.
6.  **Success**: Sam sees "Neon Vibes" in his dashboard.
7.  **Switching**: Clicks "Manage Shop". Redirected to `neon-vibes.bizbundl.com/admin` (Auto-logged in via token exchange).

## Journey 2: The Shopper (Tenant Customer)
**Persona**: "Casual Cathy" looking for neon icons.

1.  **Discovery**: Cathy lands on `neon-vibes.bizbundl.com` (Page Builder UI).
2.  **Browsing**: Views Product Grid. Adds "Icon Pack" to Cart.
3.  **Checkout**: Enters email. (Local User checking: Does `shop_neon.users` have this email? No).
4.  **Purchase**: Pays.
5.  **Account**: Account auto-created in `shop_neon.users`.
6.  **Isolation**: Cathy goes to `other-shop.bizbundl.com`. She is **NOT** logged in. (Correct).

## Journey 3: The Scaling Owner (Expansion)
**Persona**: "Serious Sam" is successful.

1.  **Dashboard**: Sam logs into `bizbundl.com`.
2.  **Expansion**: Clicks "Add Shop 2".
3.  **Creation**: Creates "Cozy Merino" (`cozy.bizbundl.com`).
4.  **Management**: Sam now sees 2 shops. He can toggle between them.
5.  **Billing**: He receives 2 separate invoices (or 1 consolidated, depending on implementation).

# UI Gotchas

Known pitfalls and non-obvious requirements for UI components in SportLink.

---

## Google Profile Pictures — `referrerPolicy: 'no-referrer'`

### Problem

Google user profile picture URLs (`lh3.googleusercontent.com`) fail to load when
the browser sends a `Referer` header. The image request returns a 403 and MUI's
`Avatar` component falls back to showing the first letter of the `alt` prop instead
of the photo.

This is **not** a data problem — the URL is stored correctly in DynamoDB.
The picture will show fine in one place and silently break in another depending on
whether `referrerPolicy` was included.

### Rule

**Every `<Avatar>` that renders a Google profile picture must include
`imgProps={{ referrerPolicy: 'no-referrer' }}`.**

```tsx
// ✅ Correct
<Avatar
  src={account.Picture}
  alt={displayName}
  imgProps={{ referrerPolicy: 'no-referrer' }}
/>

// ❌ Will silently show a letter fallback instead of the photo
<Avatar
  src={account.Picture}
  alt={displayName}
/>
```

### Where this applies

Any component that renders a user's `Picture` field fetched from `/api/account/:id`.
Current usages:

| File | Context |
|------|---------|
| `components/Layout.tsx` | Logged-in user avatar in the navbar |
| `features/matchoffer/ui/pages/MyOffersPage.tsx` | Requester avatar in the Solicitudes dialog |
| `features/matchrequest/ui/pages/MyReceivedRequestsPage.tsx` | Requester avatar in received requests list |

### Why the navbar worked but other places didn't

`Layout.tsx` was written with `referrerPolicy: 'no-referrer'` from the start.
Later components that introduced new `Avatar` usages copied the pattern without
this prop, causing the bug to reappear silently.

### Checklist

When adding a new `<Avatar>` that renders a user profile picture:

- [ ] `imgProps={{ referrerPolicy: 'no-referrer' }}` is present
- [ ] Tested with an account that has a Google profile picture set

---

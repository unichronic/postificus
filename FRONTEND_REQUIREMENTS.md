# Frontend Requirements - Postificus

## Core User Flows

### 1. Authentication & User Management
- **Login/Signup** (if user accounts are needed)
- **Session management** with JWT tokens
- **Profile management**

### 2. Platform Connection Management
**Purpose:** Connect social media accounts to publish to

**UI Requirements:**
- Platform cards (Medium, LinkedIn, Dev.to)
- "Connect Account" buttons per platform
- Status indicators (Connected/Not Connected)
- Platform-specific credential inputs:
  - **Medium:** UID, SID, XSRF (from browser cookies)
  - **LinkedIn:** li_at cookie
  - **Dev.to:** Session token
- Instructions on how to extract cookies
- Test connection functionality
- Disconnect/Re-authenticate options

**API Endpoints Used:**
- `POST /api/credentials` - Save platform credentials
- `GET /api/credentials` - List connected platforms
- `DELETE /api/credentials/{platform}` - Remove connection

### 3. Content Editor/Composer
**Purpose:** Write blog posts with markdown support

**UI Requirements:**
- Rich markdown editor with preview pane
- Split view: Editor | Preview
- Title input field
- Content textarea with syntax highlighting
- Tag input (multi-select or comma-separated)
- Featured image upload (future)
- Save as draft functionality
- Auto-save (optional)

**Components:**
- Monaco Editor or CodeMirror for markdown
- Markdown renderer (marked.js or similar)
- Tag selector component

### 4. Publishing Interface
**Purpose:** Select platforms and publish content

**UI Requirements:**
- Platform selector (checkboxes for Medium, LinkedIn, Dev.to)
- Per-platform toggle switches
- Platform-specific options:
  - **Medium:** Canonical URL, tags
  - **LinkedIn:** Share URL (for article links)
  - **Dev.to:** Tags, series
- "Publish" button (primary CTA)
- "Schedule" option (future)
- Publishing status modal/toast notifications

**API Endpoints Used:**
- `POST /publish/medium` - Publish to Medium
- `POST /publish/linkedin` - Publish to LinkedIn
- `POST /publish/devto` - Publish to Dev.to

### 5. Queue/Job Status Monitor
**Purpose:** Track publishing jobs in real-time

**UI Requirements:**
- Job list/table showing:
  - Post title
  - Target platform(s)
  - Status (Queued, Processing, Success, Failed)
  - Timestamp
  - Published URL (on success)
  - Error message (on failure)
- Real-time updates (WebSocket or polling)
- Filter by status/platform
- Retry failed jobs button
- Clear completed jobs

**API Endpoints Used:**
- `GET /api/jobs` - List all jobs
- `GET /api/jobs/{id}` - Get job details
- `POST /api/jobs/{id}/retry` - Retry failed job

### 6. Dashboard/Overview
**Purpose:** High-level view of publishing activity

**UI Requirements:**
- Summary cards:
  - Total posts published
  - Posts published this week/month
  - Connected platforms count
  - Success/failure rate
- Recent activity feed
- Quick publish button

---

## Page Structure

### 1. `/login` - Authentication
- Login form
- Sign up link (if applicable)

### 2. `/dashboard` - Home
- Overview cards
- Recent posts
- Quick actions

### 3. `/editor` or `/compose` - Content Editor
- Full-screen markdown editor
- Preview pane
- Metadata sidebar (tags, platforms)

### 4. `/platforms` or `/connections` - Platform Management
- Grid of platform cards
- Connection status
- Add/remove credentials

### 5. `/jobs` or `/queue` - Job Monitor
- Job queue table
- Filters and search
- Job details modal

### 6. `/settings` - User Settings
- Profile settings
- Default publishing options
- Notification preferences
- API keys (if exposing API to users)

---

## Key Components Needed

### 1. **Navigation Bar**
- Logo/Brand
- Links: Dashboard, Editor, Platforms, Jobs
- User menu (profile, settings, logout)

### 2. **Platform Card Component**
```
┌─────────────────────────┐
│  [Medium Logo]          │
│  Medium                 │
│  ● Connected            │
│  [Disconnect] [Test]    │
└─────────────────────────┘
```

### 3. **Markdown Editor Component**
- Editor pane (left)
- Live preview (right)
- Toolbar (bold, italic, link, etc.)

### 4. **Publishing Modal**
```
┌──────────────────────────────┐
│  Select Platforms to Publish │
│  ☑ Medium                    │
│  ☑ LinkedIn                  │
│  ☐ Dev.to                    │
│                              │
│  [Cancel]  [Publish Now]     │
└──────────────────────────────┘
```

### 5. **Job Status Component**
- Status badge (color-coded)
- Progress spinner (for in-progress jobs)
- Success/error icons
- Expandable error details

### 6. **Toast/Notification System**
- Success: "Published to Medium successfully!"
- Error: "Failed to publish to LinkedIn: {error}"
- Info: "Job queued, processing..."

---

## Technical Requirements

### State Management
- **Option A:** React Context + useReducer
- **Option B:** Redux/Zustand/Jotai
- **Recommended:** Zustand (simpler than Redux, better than Context for complex state)

**State to Manage:**
- User session/auth
- Connected platforms
- Current draft
- Job queue
- UI state (modals, toasts)

### Data Fetching
- **Option A:** fetch API with custom hooks
- **Option B:** React Query (recommended)
- **Option C:** SWR

**React Query Benefits:**
- Automatic caching
- Background refetching
- Optimistic updates
- Retry logic

### Real-time Updates (Job Queue)
- **Option A:** Polling (simple, works everywhere)
  ```js
  useQuery(['jobs'], fetchJobs, { refetchInterval: 5000 })
  ```
- **Option B:** WebSocket (more efficient)
  ```js
  const ws = new WebSocket('ws://localhost:8080/ws/jobs')
  ```
- **Recommended:** Polling for MVP, WebSocket for production

### Form Handling
- **Option A:** Controlled components (vanilla React)
- **Option B:** React Hook Form (recommended)
- **Benefits:** Validation, error handling, performance

### Styling
- **Option A:** Tailwind CSS (utility-first, fast development)
- **Option B:** CSS Modules + vanilla CSS
- **Option C:** Styled Components (CSS-in-JS)
- **Recommended:** Tailwind for rapid prototyping

### Markdown Rendering
- **Library:** `marked` or `react-markdown`
- **Syntax Highlighting:** `prismjs` or `highlight.js`

---

## API Integration

### Authentication
```js
// Login
POST /api/auth/login
{ email, password }
→ { token, user }

// Store token in localStorage or cookie
localStorage.setItem('token', token)

// Attach to all requests
headers: { Authorization: `Bearer ${token}` }
```

### Publishing Flow
```js
// 1. User writes content in editor
const [title, setTitle] = useState('')
const [content, setContent] = useState('')
const [tags, setTags] = useState([])

// 2. User selects platforms
const [platforms, setPlatforms] = useState({
  medium: true,
  linkedin: false,
  devto: true
})

// 3. User clicks "Publish"
const publish = async () => {
  const promises = []
  
  if (platforms.medium) {
    promises.push(
      fetch('/publish/medium', {
        method: 'POST',
        body: JSON.stringify({ title, content, tags })
      })
    )
  }
  
  if (platforms.devto) {
    promises.push(
      fetch('/publish/devto', {
        method: 'POST',
        body: JSON.stringify({ title, content, tags })
      })
    )
  }
  
  const results = await Promise.allSettled(promises)
  // Show success/error for each platform
}
```

### Job Monitoring
```js
// Polling approach
const { data: jobs } = useQuery(
  ['jobs'],
  () => fetch('/api/jobs').then(r => r.json()),
  { refetchInterval: 5000 }
)

// Display jobs in table
jobs.map(job => (
  <tr key={job.id}>
    <td>{job.title}</td>
    <td>{job.platform}</td>
    <td><StatusBadge status={job.status} /></td>
  </tr>
))
```

---

## Responsive Design Considerations

### Mobile-First Breakpoints
- **Mobile:** < 640px (1 column)
- **Tablet:** 640px - 1024px (2 columns)
- **Desktop:** > 1024px (3 columns, sidebars)

### Mobile Optimizations
- Stack editor/preview vertically
- Hamburger menu for navigation
- Bottom sheet for publishing modal
- Swipeable platform cards

---

## Accessibility (a11y) Requirements

- Semantic HTML (nav, main, aside, etc.)
- ARIA labels for icons/buttons
- Keyboard navigation (Tab, Enter, Escape)
- Focus management (modals trap focus)
- Color contrast (WCAG AA minimum)
- Screen reader support

---

## Performance Considerations

### Code Splitting
```js
// Lazy load routes
const Editor = lazy(() => import('./pages/Editor'))
const Jobs = lazy(() => import('./pages/Jobs'))
```

### Image Optimization
- Lazy load images
- Use WebP format
- Responsive images (srcset)

### Bundle Size
- Tree-shake unused code
- Minify for production
- Use CDN for static assets

---

## Security Considerations

### Frontend Security
- **XSS Protection:** Sanitize HTML/markdown
- **CSRF Tokens:** Include in all POST requests
- **Secure Cookies:** httpOnly, secure, sameSite
- **Input Validation:** Client-side validation (server validates too)
- **Secrets:** Never expose API keys in frontend code

### Credential Storage
- Store platform credentials on **backend only**
- Frontend only shows connection status
- Use secure API endpoints to save/update

---

## Minimal Viable Product (MVP) Scope

### Must-Have Features (Phase 1)
1. ✅ Authentication (login/logout)
2. ✅ Markdown editor
3. ✅ Platform connection management (Medium, LinkedIn, Dev.to)
4. ✅ Single-platform publishing
5. ✅ Basic job status display

### Nice-to-Have (Phase 2)
6. Multi-platform publishing (one click)
7. Real-time job updates (WebSocket)
8. Draft saving
9. Scheduling posts
10. Analytics dashboard

### Future Enhancements (Phase 3)
11. Collaborative editing
12. Template library
13. SEO optimization tools
14. Image management/upload
15. Advanced analytics
16. Publishing to more platforms

---

## Tech Stack Recommendation

### Frontend Framework
**React** with TypeScript (type safety, better DX)

### Build Tool
**Vite** (faster than Create React App)

### Routing
**React Router v6**

### State Management
**Zustand** (simple, performant)

### Data Fetching
**React Query** (caching, background refetch)

### Forms
**React Hook Form** (validation, performance)

### Styling
**Tailwind CSS** (rapid development)

### Markdown
**react-markdown** + **remark-gfm** (GitHub Flavored Markdown)

### Icons
**Lucide React** or **React Icons**

### UI Components (Optional)
- **Headless:** Radix UI or Headless UI
- **Full Suite:** Shadcn/ui (recommended - customizable, Tailwind-based)

---

## File Structure

```
frontend/
├── public/
│   └── favicon.ico
├── src/
│   ├── components/
│   │   ├── editor/
│   │   │   ├── MarkdownEditor.tsx
│   │   │   └── PreviewPane.tsx
│   │   ├── platforms/
│   │   │   ├── PlatformCard.tsx
│   │   │   └── ConnectionModal.tsx
│   │   ├── jobs/
│   │   │   ├── JobTable.tsx
│   │   │   └── StatusBadge.tsx
│   │   ├── common/
│   │   │   ├── Button.tsx
│   │   │   ├── Modal.tsx
│   │   │   └── Toast.tsx
│   │   └── layout/
│   │       ├── Navbar.tsx
│   │       └── Sidebar.tsx
│   ├── pages/
│   │   ├── Dashboard.tsx
│   │   ├── Editor.tsx
│   │   ├── Platforms.tsx
│   │   ├── Jobs.tsx
│   │   └── Settings.tsx
│   ├── hooks/
│   │   ├── useAuth.ts
│   │   ├── usePlatforms.ts
│   │   └── useJobs.ts
│   ├── api/
│   │   ├── auth.ts
│   │   ├── platforms.ts
│   │   ├── publish.ts
│   │   └── jobs.ts
│   ├── store/
│   │   ├── authStore.ts
│   │   ├── editorStore.ts
│   │   └── jobStore.ts
│   ├── types/
│   │   ├── platform.ts
│   │   ├── job.ts
│   │   └── post.ts
│   ├── utils/
│   │   ├── markdown.ts
│   │   └── validation.ts
│   ├── App.tsx
│   └── main.tsx
├── package.json
├── tsconfig.json
├── vite.config.ts
└── tailwind.config.js
```

---

## Summary

**Core Requirements:**
1. **Editor:** Markdown writing with live preview
2. **Platforms:** Connect/manage Medium, LinkedIn, Dev.to accounts
3. **Publishing:** One-click multi-platform publishing
4. **Monitoring:** Track job status in real-time
5. **Dashboard:** Overview of activity and quick actions

**Recommended Stack:**
- React + TypeScript + Vite
- Tailwind CSS for styling
- React Query for data fetching
- Zustand for state management
- React Hook Form for forms

**MVP Timeline:** 2-3 weeks for a single developer

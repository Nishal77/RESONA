# RESONA — CLAUDE.md
> Complete project context file. Paste this at the start of ANY Claude conversation.
> This file covers every single detail about the Resona project.
---
## 1. WHO IS BUILDING THIS
- **Name:** Nishal Poojary
- **Handle:** @nishalbuilds
- **College:** MCA (Master of Computer Applications) — MITE, Mangalore Institute of Technology and Engineering, Karnataka, India
- **Work:** Backend Engineer at Devify Labs (digital studio), Mangalore
- **Core stack expertise:** TypeScript, NestJS, Node.js, PostgreSQL, Redis, Docker, React
- **GitHub:** github.com/Nishal77
- **Email:** nishalnpoojary@gmail.com
- **Building solo — no teammates**
- **Purpose:** MCA Major Project + Published Research Paper
---
## 2. WHAT IS RESONA
Resona is a **vernacular-first social media web application** for Bharat's regional language creators.

**One-liner:**
> "A community + content discovery web app where Bharat's regional language creators get discovered — without speaking English."

**Think:** Reddit + Instagram, but vernacular-first. Built for Bharat.

**Why the name:** "Resona" = Resonance. When a voice hits the right frequency, it amplifies. Every regional creator deserves that amplification. The core algorithm (VRS) literally measures resonance — the name and the product are one.
---
## 3. THE PROBLEM WE ARE SOLVING
Existing platforms (Instagram, Twitter, Reddit, YouTube) optimize discovery entirely for English content.
- A Kannada creator with 500 followers gets buried under an English creator with 5000
- Regional language content gets lower algorithmic reach despite stronger community engagement
- There is no platform that surfaces content BY language natively
- There is no community structure built AROUND regional languages
- Existing platforms treat English as the default — regional languages are second-class

**What failed before:**
- ShareChat — tried vernacular but had poor community structure and content quality issues
- Koo — failed because it tried to clone Twitter instead of building vernacular-native UX
- Neither introduced a vernacular-specific scoring or discovery mechanism

**The gap Resona fills:**
No platform has combined (1) language-native discovery algorithm + (2) regional community structure + (3) transparent engagement scoring into one product.
---
## 4. TYPE OF APPLICATION
- **Responsive Web App** — NOT a mobile app
- Works perfectly on desktop browsers AND mobile browsers
- Built mobile-first — every screen designed for 375px first, then scaled up
- No App Store or Play Store required
- Single Page Application (React SPA)
- RESTful API backend
---
## 5. SUPPORTED LANGUAGES
| Language | VRS Locality Score | Primary Region |
|---|---|---|
| Kannada | 1.0 | Karnataka — HERO LANGUAGE for all demos |
| Tamil | 1.0 | Tamil Nadu |
| Telugu | 1.0 | Andhra Pradesh + Telangana |
| Malayalam | 1.0 | Kerala |
| Hindi | 1.0 | North India |
| Code-mixed (regional + English) | 0.6 | Any region |
| Pure English | 0.3 | — (lowest weight by design) |

**Important:** Kannada is the hero language. All demo data, screenshots, example posts, and seed data use Kannada content.
---
## 6. THE CORE INNOVATION — VERNACULAR RESONANCE SCORE (VRS)
This is the **original research contribution** of this project. No existing platform or academic paper defines vernacular resonance algorithmically. This is what makes the project unique and the research paper publishable.

### 6.1 The Full Formula
```
VRS = (Engagement Rate × Language Locality Score × Share Velocity)
       ÷ (Hours Since Posted + 1)
```

### 6.2 Each Component Explained

**Engagement Rate**
```
EngagementRate = (likes + comments + shares) ÷ views
```
Measures how actively people engage per view. A post with 100 views and 20 engagements scores higher than one with 10000 views and 20 engagements.

**Language Locality Score**
Automatically calculated from Google Cloud Translate API detection:
- Pure regional language (Kannada, Tamil, Telugu, Malayalam, Hindi) → `1.0`
- Code-mixed content (regional + English mixed) → `0.6`
- Pure English → `0.3`

**Share Velocity**
```
ShareVelocity = total shares within first 2 hours after posting
```
Early shares signal strong community resonance. A post shared 10 times in 2 hours beats one shared 50 times over 3 days.

**Time Decay**
```
(Hours Since Posted + 1)
```
Denominator causes natural decay — older posts lose VRS unless they keep getting engagement. Keeps feed fresh.

### 6.3 How VRS Drives the App
- Every post gets a `vrs_score` stored in the database
- Home feed is sorted by `vrs_score` DESC — NOT chronological, NOT follower count
- VRS score badge is displayed publicly on every post card — full transparency, no black box
- Background cron job recalculates `vrs_score` for all posts every 30 minutes
- Redis caches the sorted feed between recalculations
- Explore page surfaces posts that cross a global VRS threshold (cross-language viral)

### 6.4 VRS on User Profiles
- `vrs_total` on the users table = sum of all VRS scores across all their posts
- Displayed on profile as "Total Resonance Score"
- Acts as a creator reputation metric — pure vernacular creators rank higher

### 6.5 Language Detection Flow
```
User creates post
→ Content sent to Google Cloud Translate API
→ API returns: detected_language + language_confidence (0.0 to 1.0)
→ If confidence >= 0.70 → auto-set detected_language, calculate locality_score
→ If confidence < 0.70 → show language selector to user, ask them to confirm
→ Store detected_language, language_confidence, language_locality_score in posts table
→ VRS calculation uses language_locality_score
```
---
## 7. COMPLETE FEATURE LIST

### Feature 1 — Authentication System
- Email + password registration with validation
- Email + password login
- Google OAuth 2.0 login (one-click sign in)
- JWT access token (expires in 15 minutes)
- JWT refresh token (expires in 7 days, stored in httpOnly cookie)
- Refresh token rotation on every use
- Logout invalidates refresh token
- All private routes protected via NestJS JwtAuthGuard
- `@CurrentUser()` decorator on all protected controllers

### Feature 2 — Onboarding Flow
- Shown only to new users after first login/register
- **Step 1:** Pick your primary language — mandatory, cannot skip — 6 options shown as large clickable cards
- **Step 2:** Pick your state and city — optional but recommended
- **Step 3:** Follow 3 starter communities — curated list shown based on chosen language
- After completing onboarding → redirect to /home
- Onboarding state tracked in `users.onboarding_completed` boolean
- If user refreshes during onboarding — resume from last completed step

### Feature 3 — Home Feed
- Posts sorted by `vrs_score` DESC (highest VRS at top)
- Language filter bar at top of feed:
  - "All" tab — shows posts from all languages
  - Individual language tabs (Kannada / Tamil / Telugu / Malayalam / Hindi)
  - Filter is instant — no page reload
- Shows only posts from communities user has joined + posts from users they follow
- New users (no follows/communities) → see top VRS posts from their primary language
- Infinite scroll pagination (20 posts per page, load more on scroll)
- Floating "+" Create Post button — fixed bottom right on mobile
- Pull to refresh on mobile
- Each post card in feed shows:
  - User avatar + username + language badge
  - Post content (truncated at 280 chars with "Read more")
  - Media thumbnail if image/video
  - VRS badge (score + color coded: green >0.8, amber 0.4-0.8, grey <0.4)
  - Like count, comment count, share count
  - Community name it was posted in
  - Time since posted ("2h ago", "3d ago")
  - Like / Comment / Share / Save action buttons

### Feature 4 — Post Creation
- Accessible via floating "+" button or /create route
- Post types supported:
  - **Text post** — plain text, up to 2000 characters
  - **Image post** — single image upload, max 10MB, formats: jpg/png/webp
  - **Video post** — short video, max 50MB, formats: mp4/webm
- Language detection:
  - Runs automatically as user types (debounced 1 second)
  - Shows detected language + confidence score live in UI
  - If confidence < 70% → shows language dropdown for manual selection
- Tag input:
  - User types tags with # prefix
  - Auto-suggests existing popular tags in their language
  - Max 5 tags per post
- Community selector:
  - Dropdown shows communities user has joined
  - Can post to general feed (no community selected)
- Submit → POST /posts API call
- On success → redirect to new post's detail page
- Media upload flow:
  - File selected → POST /upload/image or /upload/video
  - Supabase Storage returns public URL
  - URL stored in post creation payload

### Feature 5 — Post Detail Page
- Route: /posts/:id
- Shows full post content (no truncation)
- VRS score badge — large, prominently displayed at top
- Language badge showing detected language
- Author info — avatar, username, language badge, follow button
- Media — full image or video player if media attached
- Engagement bar — like count, comment count, share count, save count
- Action buttons — Like / Comment / Share / Save
- Share → copies post URL to clipboard + shows "Copied!" toast
- Comments section below post:
  - Top-level comments sorted by like count
  - Each comment shows avatar, username, content, like button, reply button
  - Nested replies (1 level deep — replies to replies not supported)
  - Load more comments button (10 per page)
  - Comment input box at bottom (logged-in users only)
  - Guest users see "Login to comment" prompt

### Feature 6 — Explore Page
- Route: /explore
- **Trending Tags section:**
  - Shows top 10 trending tags for user's primary language
  - Language selector to switch (Kannada / Tamil / Telugu / Malayalam / Hindi / All)
  - Each tag shows name + usage count + 24h growth percentage
  - Click tag → opens tag feed page /explore/tags/:tagname
- **Top Communities section:**
  - Shows 6 communities sorted by (member_count + weekly_post_count)
  - Each community card shows: name, language, member count, description, Join button
- **Cross-language Viral section:**
  - Posts with VRS > global threshold (recalculated daily)
  - Labelled "Trending Across Bharat"
  - Shows posts from all languages — promotes cross-community discovery
- **Search bar:**
  - Full-text search across post content
  - Language filter dropdown alongside search
  - Results sorted by VRS score
  - Instant results as user types (debounced 300ms)

### Feature 7 — Communities
- Route: /c/:slug
- **Community Detail Page:**
  - Banner image + avatar + community name
  - Language badge + member count + post count
  - Description
  - Join / Leave button
  - "Snap of the Week" banner — shows the highest VRS post of the current week, auto-selected every Monday at midnight via cron job
  - Community feed (posts sorted by VRS, filterable by Latest/Top/Hot)
  - Member list (visible to joined members)
- **Community Creation:**
  - Any logged-in user can create a community
  - Fields: name, slug (auto-generated from name), description, primary language, secondary languages, avatar upload
  - Creator automatically becomes admin
- **Community Roles:**
  - member — can post, comment, like
  - moderator — can delete posts/comments, pin posts
  - admin — full control including role management
- **Join/Leave logic:**
  - Join → insert into community_members, increment community.member_count
  - Leave → delete from community_members, decrement community.member_count
- **Post in Community:**
  - When creating post, user can select a community from dropdown
  - Post appears in both general feed and community feed
- **Weekly Snap of the Week:**
  - Cron job runs every Monday at 00:00
  - Selects post with highest vrs_score from past 7 days in each community
  - Stored in community as `snap_of_week_post_id`
  - Displayed as highlighted banner at top of community page

### Feature 8 — User Profiles
- Route: /:username
- **Own Profile:**
  - Editable via Settings
- **Public Profile:**
  - Avatar (Supabase Storage URL)
  - Full name + username
  - Bio (max 160 chars)
  - Language badge (primary language)
  - State + city
  - Total VRS score ("Resonance Score")
  - Follower count + Following count (both clickable → opens list modal)
  - "Follow" button (if viewing another user's profile)
  - Post grid — user's posts sorted by vrs_score (3 columns on desktop, 2 on mobile)
  - Tabs on profile: Posts / Communities / Saved (own profile only)
- **Follow / Unfollow:**
  - Follow → insert into follows, increment follower_count and following_count
  - Unfollow → delete from follows, decrement counts
  - Follow creates a notification for the followed user

### Feature 9 — Notifications
- Route: /notifications
- **Notification types:**
  - `like` — "username liked your post"
  - `comment` — "username commented on your post"
  - `follow` — "username started following you"
  - `trending` — "Your post is trending in the Kannada community" (triggered when post VRS crosses community threshold)
- **Notification bell in navbar:**
  - Shows unread count badge (red dot with number)
  - Clicking bell → opens notifications dropdown (top 5 recent)
  - "See all" link → goes to /notifications page
- **Full notifications page:**
  - Paginated list of all notifications (20 per page)
  - Each notification shows: actor avatar, message, time, read/unread state
  - Click notification → navigates to relevant post/profile
  - "Mark all as read" button at top
  - Individual mark as read on click
- **Notifications created by:**
  - Like action → create notification for post owner
  - Comment action → create notification for post owner
  - Follow action → create notification for followed user
  - VRS scheduler → creates trending notification when post VRS crosses threshold

### Feature 10 — Settings
- Route: /settings
- **Profile section:**
  - Change avatar (upload to Supabase Storage)
  - Change full name
  - Change bio
  - Change state/city
- **Language section:**
  - Change primary language
  - This changes the default feed language filter
- **Account section:**
  - Change email address (requires password confirmation)
  - Change password (requires current password)
- **Danger zone:**
  - Delete account (requires typing "DELETE" to confirm)
  - Deletes user + all their posts + engagements (cascading deletes in DB)
---
## 8. ALL 10 PAGES WITH ROUTES
| Page | Route | Auth Required | Description |
|---|---|---|---|
| Landing | / | No | Public homepage with hero, language preview, sign up CTA |
| Onboarding | /onboarding | Yes | 3-step language/community setup for new users |
| Home Feed | /home | Yes | Main feed sorted by VRS |
| Explore | /explore | No | Trending tags, communities, cross-language viral |
| Post Detail | /posts/:id | No (view) / Yes (engage) | Full post with comments |
| Create Post | /create | Yes | Text/image/video post creation |
| Community | /c/:slug | No (view) / Yes (join/post) | Community feed and detail |
| Profile | /:username | No (view) / Yes (follow) | User profile page |
| Notifications | /notifications | Yes | All notifications |
| Settings | /settings | Yes | Profile and account settings |
---
## 9. TECH STACK — COMPLETE

### Frontend
| Layer | Technology | Version | Notes |
|---|---|---|---|
| Framework | React | 18.x | Functional components + hooks only |
| Build tool | Vite | 5.x | Fast HMR, fast builds |
| Language | TypeScript | 5.x | Strictly typed, no `any` |
| Styling | TailwindCSS | 3.x | Utility-first, mobile-first |
| Routing | React Router | v6 | Declarative routes, nested routes |
| Data fetching | TanStack Query | v5 | Caching, background refetch, infinite scroll |
| HTTP client | Axios | 1.x | JWT interceptor, refresh token logic |
| State management | Zustand | 4.x | Auth state only — server state via React Query |
| Form handling | React Hook Form | 7.x | Validation, controlled inputs |
| Icons | Lucide React | latest | Consistent icon set |
| Toasts | React Hot Toast | latest | Success/error notifications |
| Date formatting | date-fns | latest | "2h ago" style formatting |

### Backend
| Layer | Technology | Version | Notes |
|---|---|---|---|
| Language | Go | 1.22+ | Idiomatic Go, no `interface{}` abuse |
| Framework | Gin | 1.9.x | HTTP router + middleware |
| Database | PostgreSQL | 15.x | Via Supabase |
| ORM | GORM | 2.x | Models, migrations, associations |
| Cache | Redis | 7.x | go-redis/v9 client |
| Auth JWT | golang-jwt/jwt | v5 | Access + refresh tokens |
| Google OAuth | golang.org/x/oauth2 | latest | Google strategy |
| Validation | go-playground/validator | v10 | Struct tag validation |
| File upload | multipart/form-data | stdlib | In-memory before Supabase upload |
| Scheduler | robfig/cron | v3 | VRS cron job every 30 min |
| Config | godotenv + viper | latest | .env management |
| HTTP client | net/http | stdlib | For Google Translate API calls |
| Password hash | bcrypt | golang.org/x/crypto | 10 cost rounds |
| UUID | google/uuid | latest | UUID v4 generation |
| Logger | zerolog | latest | Structured JSON logging |

### Database & Storage
| Service | Purpose | Why |
|---|---|---|
| Supabase PostgreSQL | Primary database | Managed PostgreSQL, generous free tier, excellent dashboard |
| Supabase Storage | Image + video uploads | Replaces AWS S3 entirely — no AWS account needed |
| Redis | Feed cache + rate limiting | Upstash Redis (serverless Redis, free tier) |

### External APIs
| API | Purpose | Notes |
|---|---|---|
| Google Cloud Translate API | Language detection on every post | Free tier: 500k chars/month |
| Google OAuth 2.0 | Social login | Via Google Cloud Console |

### DevOps
| Tool | Purpose |
|---|---|
| Docker | Backend container |
| docker-compose | Local dev environment (NestJS + PostgreSQL + Redis) |
| Railway | Backend production deployment |
| Vercel | Frontend production deployment |
| GitHub | Version control |
---
## 10. DATABASE SCHEMA — ALL TABLES

### users
```sql
CREATE TABLE users (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username         VARCHAR(50) UNIQUE NOT NULL,
  email            VARCHAR(255) UNIQUE NOT NULL,
  password_hash    VARCHAR,                                    -- null for Google OAuth users
  google_id        VARCHAR UNIQUE,                             -- for Google OAuth users
  full_name        VARCHAR(100),
  avatar_url       VARCHAR,                                    -- Supabase Storage public URL
  bio              TEXT,                                       -- max 160 chars enforced in DTO
  primary_language VARCHAR(20) NOT NULL DEFAULT 'kannada',    -- kannada | tamil | telugu | malayalam | hindi | english
  state            VARCHAR(100),
  city             VARCHAR(100),
  vrs_total        FLOAT DEFAULT 0,                           -- sum of all post vrs_scores, updated on each recalculation
  follower_count   INT DEFAULT 0,
  following_count  INT DEFAULT 0,
  onboarding_completed BOOLEAN DEFAULT FALSE,                 -- false = show onboarding flow after login
  created_at       TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at       TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### posts
```sql
CREATE TABLE posts (
  id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id                  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  content_text             TEXT,                              -- null if media-only post
  media_url                VARCHAR,                           -- Supabase Storage public URL, null if text-only
  media_type               VARCHAR(10),                       -- image | video | null
  community_id             UUID REFERENCES communities(id) ON DELETE SET NULL,  -- null = posted to general feed
  detected_language        VARCHAR(20),                       -- result from Google Translate API
  language_confidence      FLOAT,                             -- 0.0 to 1.0 confidence score from API
  language_locality_score  FLOAT NOT NULL DEFAULT 0.3,        -- 0.3 | 0.6 | 1.0 — calculated from detected_language
  vrs_score                FLOAT NOT NULL DEFAULT 0,          -- recalculated every 30 minutes by VRS scheduler
  like_count               INT NOT NULL DEFAULT 0,
  comment_count            INT NOT NULL DEFAULT 0,
  share_count              INT NOT NULL DEFAULT 0,
  view_count               INT NOT NULL DEFAULT 0,
  save_count               INT NOT NULL DEFAULT 0,
  created_at               TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at               TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_posts_vrs_score ON posts(vrs_score DESC);
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_community_id ON posts(community_id);
CREATE INDEX idx_posts_detected_language ON posts(detected_language);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
```

### communities
```sql
CREATE TABLE communities (
  id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name                  VARCHAR(100) NOT NULL,
  slug                  VARCHAR(100) UNIQUE NOT NULL,         -- url-safe version of name, e.g. "kannada-memes"
  description           TEXT,
  primary_language      VARCHAR(20) NOT NULL,
  secondary_languages   VARCHAR(20)[] DEFAULT '{}',
  avatar_url            VARCHAR,                              -- Supabase Storage public URL
  banner_url            VARCHAR,                              -- Supabase Storage public URL
  member_count          INT NOT NULL DEFAULT 0,
  post_count            INT NOT NULL DEFAULT 0,
  snap_of_week_post_id  UUID REFERENCES posts(id) ON DELETE SET NULL,  -- updated every Monday midnight by cron
  created_by            UUID REFERENCES users(id) ON DELETE SET NULL,
  created_at            TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_communities_primary_language ON communities(primary_language);
CREATE INDEX idx_communities_slug ON communities(slug);
```

### community_members
```sql
CREATE TABLE community_members (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  community_id  UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
  user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role          VARCHAR(20) NOT NULL DEFAULT 'member',        -- member | moderator | admin
  joined_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  UNIQUE(community_id, user_id)
);

CREATE INDEX idx_community_members_user_id ON community_members(user_id);
CREATE INDEX idx_community_members_community_id ON community_members(community_id);
```

### comments
```sql
CREATE TABLE comments (
  id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  post_id            UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  user_id            UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  content            TEXT NOT NULL,
  parent_comment_id  UUID REFERENCES comments(id) ON DELETE CASCADE,
  -- null = top-level comment, set = reply (1 level deep only)
  like_count         INT NOT NULL DEFAULT 0,
  created_at         TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_parent_comment_id ON comments(parent_comment_id);
```

### engagements
```sql
CREATE TABLE engagements (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  post_id     UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  type        VARCHAR(10) NOT NULL,                           -- like | share | save | view
  created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  UNIQUE(post_id, user_id, type)                             -- prevents duplicate likes/saves/views per user per post
);

CREATE INDEX idx_engagements_post_id ON engagements(post_id);
CREATE INDEX idx_engagements_created_at ON engagements(created_at);
-- Critical for VRS share velocity calculation (first 2 hours)
CREATE INDEX idx_engagements_post_type_created ON engagements(post_id, type, created_at);
```

### tags
```sql
CREATE TABLE tags (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name            VARCHAR(100) NOT NULL,
  language        VARCHAR(20),                                -- the language this tag primarily belongs to
  usage_count     INT NOT NULL DEFAULT 0,
  trending_score  FLOAT NOT NULL DEFAULT 0,                  -- recalculated daily based on recent usage
  created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  UNIQUE(name, language)
);

CREATE INDEX idx_tags_language_trending ON tags(language, trending_score DESC);
```

### post_tags
```sql
CREATE TABLE post_tags (
  post_id  UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  tag_id   UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
  PRIMARY KEY (post_id, tag_id)
);
```

### follows
```sql
CREATE TABLE follows (
  follower_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  following_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  PRIMARY KEY (follower_id, following_id),
  CHECK (follower_id != following_id)                        -- users cannot follow themselves
);

CREATE INDEX idx_follows_follower_id ON follows(follower_id);
CREATE INDEX idx_follows_following_id ON follows(following_id);
```

### notifications
```sql
CREATE TABLE notifications (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,  -- the user who receives this notification
  type        VARCHAR(20) NOT NULL,                          -- like | comment | follow | trending
  actor_id    UUID REFERENCES users(id) ON DELETE CASCADE,  -- the user who triggered this (null for system notifications)
  post_id     UUID REFERENCES posts(id) ON DELETE CASCADE,  -- the related post (null for follow notifications)
  message     TEXT,                                          -- pre-computed e.g. "nishal liked your post"
  read        BOOLEAN NOT NULL DEFAULT FALSE,
  created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_id_read ON notifications(user_id, read);
CREATE INDEX idx_notifications_user_id_created ON notifications(user_id, created_at DESC);
```

### refresh_tokens
```sql
CREATE TABLE refresh_tokens (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash  VARCHAR NOT NULL,                              -- bcrypt hash of the refresh token
  expires_at  TIMESTAMP WITH TIME ZONE NOT NULL,
  created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
```
---
## 11. ALL API ENDPOINTS

### Response format (ALL endpoints)
```json
Success:   { "success": true, "data": {}, "message": "Operation successful" }
Error:     { "success": false, "error": "Error message", "statusCode": 400 }
Paginated: { "success": true, "data": [], "meta": { "page": 1, "limit": 20, "total": 100, "hasMore": true } }
```

### Auth — /auth
```
POST   /auth/register
       Body: { username, email, password, fullName }
       Returns: { accessToken, user }

POST   /auth/login
       Body: { email, password }
       Returns: { accessToken, user }
       Sets: httpOnly cookie with refreshToken

POST   /auth/google
       Body: { googleToken }
       Returns: { accessToken, user }

POST   /auth/refresh
       Cookie: refreshToken
       Returns: { accessToken }

POST   /auth/logout
       Header: Authorization Bearer
       Clears: refreshToken cookie, invalidates token in DB
```

### Users — /users
```
GET    /users/me
GET    /users/me/feed          Query: { page?, limit?, language? }
GET    /users/me/communities
PUT    /users/me               Body: { fullName?, bio?, state?, city?, primaryLanguage?, avatarUrl? }
GET    /users/:username
GET    /users/:id/posts        Query: { page?, limit? }
GET    /users/:id/followers    Query: { page?, limit? }
GET    /users/:id/following    Query: { page?, limit? }
POST   /users/:id/follow
DELETE /users/:id/follow
```

### Posts — /posts
```
GET    /posts                  Query: { page?, limit?, language? }
POST   /posts                  Body: { contentText?, mediaUrl?, mediaType?, communityId?, tagIds? }
GET    /posts/:id
DELETE /posts/:id
POST   /posts/:id/like
DELETE /posts/:id/like
POST   /posts/:id/share
POST   /posts/:id/view
POST   /posts/:id/save
DELETE /posts/:id/save
GET    /posts/:id/comments     Query: { page?, limit? }
POST   /posts/:id/comments     Body: { content, parentCommentId? }
DELETE /posts/:id/comments/:commentId
POST   /posts/:id/comments/:commentId/like
```

### Communities — /communities
```
GET    /communities            Query: { page?, limit?, language? }
POST   /communities            Body: { name, description, primaryLanguage, secondaryLanguages?, avatarUrl?, bannerUrl? }
GET    /communities/:slug
PUT    /communities/:id
GET    /communities/:id/posts  Query: { page?, limit?, sort? (vrs|latest|top) }
GET    /communities/:id/members
GET    /communities/:id/snap-of-week
POST   /communities/:id/join
DELETE /communities/:id/join
PUT    /communities/:id/members/:userId/role   Body: { role }
```

### Explore — /explore
```
GET    /explore/trending                Query: { language?, limit? }
GET    /explore/tags                    Query: { language?, limit? }
GET    /explore/tags/:tagName/posts     Query: { page?, limit?, language? }
GET    /explore/communities             Query: { language?, limit? }
GET    /explore/search                  Query: { q, language?, page?, limit? }
```

### Upload — /upload
```
POST   /upload/image     Body: multipart/form-data { file }   Validates: jpg|png|webp, max 10MB
POST   /upload/video     Body: multipart/form-data { file }   Validates: mp4|webm, max 50MB
```

### Notifications — /notifications
```
GET    /notifications              Query: { page?, limit? }
GET    /notifications/unread-count
PUT    /notifications/:id/read
PUT    /notifications/read-all
```

### Language — /language (internal)
```
POST   /language/detect    Body: { text }    Returns: { language, confidence }
Note: called internally by posts service — not directly by frontend
```
---
## 12. VRS SERVICE IMPLEMENTATION

### vrs.service.ts — Core calculation
```typescript
// Calculates VRS for a single post
calculateVRS(postId: string): Promise<number>
// Steps:
// 1. Fetch post (like_count, comment_count, share_count, view_count, created_at, language_locality_score)
// 2. EngagementRate = (likes + comments + shares) / views (0 if views = 0)
// 3. Share velocity = count of shares WHERE type='share' AND created_at <= post.created_at + 2 hours
// 4. hoursSincePosted = (NOW - created_at) in hours
// 5. VRS = (engagementRate * localityScore * shareVelocity) / (hoursSincePosted + 1)
// 6. UPDATE posts SET vrs_score = VRS WHERE id = postId
// 7. Return VRS value
```

### vrs.scheduler.ts — Background job
```typescript
// Runs every 30 minutes: '0,30 * * * *'
// 1. Fetch ALL post IDs (paginated batches of 100)
// 2. For each batch: calculateVRS() in parallel (Promise.all)
// 3. UPDATE users SET vrs_total = SUM(post.vrs_score) for each user
// 4. Update trending_score on tags based on recent usage
// 5. Trigger "Snap of the Week" check (runs only on Mondays at 00:00)
// 6. Check for posts crossing VRS threshold → create trending notifications
// 7. Invalidate Redis feed cache keys
```

### Redis caching strategy
```
Key: feed:language:{language}:page:{page}
Key: feed:user:{userId}:page:{page}
Key: explore:trending:{language}
TTL: 30 minutes (matches VRS recalculation interval)
On VRS recalculation complete → DEL all feed:* keys → next request rebuilds cache
```
---
## 13. BACKEND FOLDER STRUCTURE
```
backend/               ← lives inside the root (frontend) folder
├── cmd/
│   └── server/
│       └── main.go                        ← Entry point, wires everything
├── internal/
│   ├── auth/
│   │   ├── handler.go                     ← Gin route handlers
│   │   ├── service.go                     ← Business logic
│   │   ├── repository.go                  ← DB queries
│   │   └── dto.go                         ← Request/response structs
│   ├── users/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── dto.go
│   ├── posts/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── dto.go
│   ├── communities/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── dto.go
│   ├── comments/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── dto.go
│   ├── engagements/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   ├── explore/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   ├── notifications/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   ├── upload/
│   │   ├── handler.go
│   │   └── service.go                     ← Supabase Storage integration
│   ├── tags/
│   │   ├── service.go
│   │   └── repository.go
│   ├── language/
│   │   └── service.go                     ← Google Translate API wrapper
│   ├── vrs/
│   │   ├── service.go                     ← VRS calculation engine ← THESIS CORE
│   │   └── scheduler.go                   ← robfig/cron job every 30 minutes
│   └── middleware/
│       ├── auth.go                        ← JWT validation middleware
│       ├── optional_auth.go               ← Optional JWT (guest-viewable routes)
│       └── response.go                    ← Standard response helpers
├── pkg/
│   ├── database/
│   │   └── postgres.go                    ← GORM connection + auto-migrate
│   ├── redis/
│   │   └── client.go                      ← go-redis client wrapper
│   ├── supabase/
│   │   └── storage.go                     ← Supabase Storage REST client
│   ├── models/
│   │   ├── user.go                        ← GORM model structs
│   │   ├── post.go
│   │   ├── community.go
│   │   ├── community_member.go
│   │   ├── comment.go
│   │   ├── engagement.go
│   │   ├── tag.go
│   │   ├── post_tag.go
│   │   ├── follow.go
│   │   ├── notification.go
│   │   └── refresh_token.go
│   └── config/
│       └── config.go                      ← Viper config loader
├── docker-compose.yml
├── Dockerfile
├── .env
├── .env.example
└── go.mod
```

### Go patterns used throughout:
- Handler → Service → Repository layering (no circular deps)
- Dependency injection via constructor functions (no global state)
- `context.Context` threaded through every DB/Redis call
- Errors wrapped with `fmt.Errorf("operation: %w", err)`
- All HTTP responses via `pkg/models/response.go` helpers
- GORM models in `pkg/models/` — shared across all internal packages
---
## 14. PROJECT FOLDER STRUCTURE (ROOT = FRONTEND)
```
major/                                 ← ROOT = Frontend (Vite/React)
├── src/
│   ├── pages/
│   │   ├── Landing.tsx          → /
│   │   ├── Onboarding.tsx       → /onboarding
│   │   ├── Home.tsx             → /home
│   │   ├── Explore.tsx          → /explore
│   │   ├── PostDetail.tsx       → /posts/:id
│   │   ├── CreatePost.tsx       → /create
│   │   ├── Community.tsx        → /c/:slug
│   │   ├── Profile.tsx          → /:username
│   │   ├── Notifications.tsx    → /notifications
│   │   └── Settings.tsx         → /settings
│   ├── components/
│   │   ├── post/
│   │   │   ├── PostCard.tsx         ← VRS badge always shown
│   │   │   ├── PostActions.tsx      ← Like/Comment/Share/Save
│   │   │   ├── PostMedia.tsx        ← Image/Video renderer
│   │   │   └── CommentItem.tsx      ← Single comment with reply
│   │   ├── community/
│   │   │   ├── CommunityCard.tsx
│   │   │   └── SnapOfWeek.tsx       ← Weekly top post banner
│   │   ├── user/
│   │   │   ├── UserAvatar.tsx
│   │   │   └── FollowButton.tsx
│   │   ├── vrs/
│   │   │   └── VRSBadge.tsx         ← THE VRS score visual — green/amber/grey
│   │   ├── common/
│   │   │   ├── LanguageBadge.tsx
│   │   │   ├── LanguageFilter.tsx
│   │   │   ├── InfiniteScroll.tsx
│   │   │   ├── Modal.tsx
│   │   │   ├── Spinner.tsx
│   │   │   └── EmptyState.tsx
│   │   ├── layout/
│   │   │   ├── Navbar.tsx           ← Top bar: logo, search, notifications bell
│   │   │   ├── Sidebar.tsx          ← Desktop only
│   │   │   ├── BottomNav.tsx        ← Mobile only
│   │   │   └── AppLayout.tsx
│   │   └── forms/
│   │       └── CreatePostModal.tsx
│   ├── hooks/
│   │   ├── useAuth.ts
│   │   ├── useFeed.ts
│   │   ├── useCommunityFeed.ts
│   │   ├── usePost.ts
│   │   ├── useInfiniteScroll.ts
│   │   └── useNotifications.ts
│   ├── api/
│   │   ├── client.ts                ← Axios instance, JWT interceptor, refresh logic
│   │   ├── auth.api.ts
│   │   ├── posts.api.ts
│   │   ├── communities.api.ts
│   │   ├── users.api.ts
│   │   ├── explore.api.ts
│   │   ├── upload.api.ts
│   │   └── notifications.api.ts
│   ├── store/
│   │   └── auth.store.ts            ← Zustand: { user, accessToken, setUser, logout }
│   ├── types/
│   │   ├── user.types.ts
│   │   ├── post.types.ts
│   │   ├── community.types.ts
│   │   ├── notification.types.ts
│   │   └── api.types.ts             ← ApiResponse<T>, PaginatedResponse<T>
│   ├── utils/
│   │   ├── formatDate.ts
│   │   ├── formatVRS.ts
│   │   └── languageUtils.ts
│   ├── constants/
│   │   └── languages.ts
│   ├── App.tsx
│   └── main.tsx
├── package.json
├── vite.config.ts
├── tailwind.config.ts
├── tsconfig.json
├── index.html
├── .env.example
└── backend/                           ← Go backend lives here
    ├── cmd/server/main.go
    ├── internal/
    ├── pkg/
    ├── go.mod
    ├── docker-compose.yml
    ├── Dockerfile
    └── .env.example
```
---
## 15. ENVIRONMENT VARIABLES

### Backend (.env)
```env
# Supabase
DATABASE_URL=postgresql://postgres:[password]@db.[project-ref].supabase.co:5432/postgres
SUPABASE_URL=https://[project-ref].supabase.co
SUPABASE_SERVICE_KEY=your_supabase_service_role_key
SUPABASE_STORAGE_BUCKET=resona-media

# Redis (Upstash)
REDIS_URL=rediss://:[password]@[host].upstash.io:6379

# JWT
JWT_SECRET=your_long_random_jwt_secret_here
JWT_REFRESH_SECRET=your_long_random_refresh_secret_here
JWT_EXPIRES_IN=15m
JWT_REFRESH_EXPIRES_IN=7d

# Google OAuth
GOOGLE_CLIENT_ID=your_google_client_id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_CALLBACK_URL=http://localhost:8080/auth/google/callback

# Google Translate API
GOOGLE_TRANSLATE_API_KEY=your_google_translate_api_key

# App Config
PORT=8080
FRONTEND_URL=http://localhost:5173
APP_ENV=development

# VRS Config
VRS_TRENDING_THRESHOLD=0.75
VRS_SHARE_VELOCITY_HOURS=2
```

> Note: Port is 8080 (Go convention). Frontend VITE_API_URL must match.

### Frontend (.env)
```env
VITE_API_URL=http://localhost:3000
VITE_GOOGLE_CLIENT_ID=your_google_client_id.apps.googleusercontent.com
```
---
## 16. DOCKER SETUP (LOCAL DEV)
```yaml
# docker-compose.yml
version: '3.8'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: resona
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
volumes:
  postgres_data:
  redis_data:
```
> Note: In production, Supabase provides PostgreSQL. Docker postgres is for local dev only.
> Redis in production = Upstash (serverless Redis, free tier, zero config).
---
## 17. KEY CODING RULES — ALWAYS FOLLOW
1. **Always Go on backend** — idiomatic Go, no `interface{}` abuse, no global state, explicit error handling
2. **Go patterns strictly** — Handler → Service → Repository per feature, constructor DI, context threading
3. **VRS is sacred** — never simplify, shortcut, or remove the VRS algorithm or VRS badge
4. **VRSBadge on every PostCard** — always visible, always showing the score, no exceptions
5. **Mobile-first responsive** — design for 375px first, then md: and lg: breakpoints
6. **Kannada is the hero language** — all seed data, demos, example content uses Kannada
7. **No over-engineering** — solo project, keep it shippable not perfect
8. **Language detection on every post** — Google Translate API called server-side on post create
9. **Redis for feed caching** — never query PostgreSQL on every feed request
10. **UUID everywhere** — all primary keys are UUID, no integer IDs
11. **go-playground/validator on all request structs** — validate binding tags on every handler input
12. **Supabase Storage for all media** — never local file storage, no AWS
13. **Standard API response always:**
    ```json
    { "success": true, "data": {}, "message": "..." }
    { "success": false, "error": "...", "statusCode": 400 }
    ```
14. **httpOnly cookie for refresh token** — never store refresh token in localStorage
15. **Indexes on all foreign keys and sort columns** — especially posts.vrs_score
---
## 18. SEED DATA

### Starter Communities
```
slug: "kannada-memes"     | language: kannada   | name: "Kannada Memes"
slug: "mangalore-vibes"   | language: kannada   | name: "Mangalore Vibes"
slug: "tamil-poetry"      | language: tamil     | name: "Tamil Poetry"
slug: "telugu-trends"     | language: telugu    | name: "Telugu Trends"
slug: "malayalam-humour"  | language: malayalam | name: "Malayalam Humour"
slug: "hindi-shayari"     | language: hindi     | name: "Hindi Shayari"
slug: "bharat-builders"   | language: kannada   | name: "Bharat Builders"
slug: "coastal-karnataka" | language: kannada   | name: "Coastal Karnataka"
```

### Demo user for testing
```
username: nishal_demo
email: demo@resona.in
password: Demo@1234
primaryLanguage: kannada
city: Mangalore
state: Karnataka
```
---
## 19. CURRENT PROJECT STATUS
> **Backend and frontend fully scaffolded. All layers written.**

### Backend (Go/Gin) — complete:
- All 10+ modules: auth, users, posts, communities, comments, engagements, explore, notifications, upload, tags, language, VRS
- GORM models for all 10 DB tables
- VRS service + 30-min cron scheduler
- JWT auth with httpOnly refresh token cookie
- Redis feed cache invalidation on VRS recalculation
- Supabase Storage upload client
- Google Translate language detection service
- docker-compose.yml for local PostgreSQL + Redis

### Frontend (React/Vite/TypeScript) — complete:
- All 10 pages with full routing
- PostCard always shows VRSBadge (sacred rule upheld)
- Language-native LanguageFilter and LanguageBadge
- TanStack Query infinite scroll feed
- Zustand auth store with access token refresh interceptor
- Mobile-first design with BottomNav and Navbar

### To run locally:
```bash
# From project root (major/)

# Backend deps + DB
cd backend && go mod tidy
docker-compose up -d
cp .env.example .env   # fill in keys
go run ./cmd/server    # starts on :8080

# Frontend (new terminal, from major/)
cd ..
npm install
cp .env.example .env   # VITE_API_URL=http://localhost:8080
npm run dev            # starts on :5173
```
---
*CLAUDE.md — Resona — Nishal Poojary — MITE MCA Major Project — 2026*

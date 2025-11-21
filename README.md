# Project Management SaaS Backend (Golang + GORM + PostgreSQL)

Backend untuk aplikasi **Project Management + Real-Time Chat** yang dirancang sebagai SaaS multi-tenant.  
User dapat membuat workspace, mengelola project, task board (kanban), hingga berkolaborasi lewat chat seperti WhatsApp.

---

## Fitur Utama (MVP)

### Authentication & Users
- Register & Login (JWT / Session)
- Manajemen profil user
- Avatar upload
- Keamanan password (bcrypt)

### Multi-Tenant Workspace
- Create workspace
- Invite member (role: owner, admin, member)
- Workspace billing-ready (free/standard/pro)
- Unique workspace slug

### Project Management (Trello-like)
- Project per workspace
- Board → Column → Task (Kanban)
- Drag & drop task
- Task assignees
- Task comments
- Attachments
- Activity log

### Real-Time Chat (WhatsApp-style)
- Channel per workspace
- Channel per project
- Group chat & direct message
- Message reply & read status (future)

---

## Arsitektur & Teknologi

### Backend:
- **Golang 1.22+**
- **Gin** 
- **GORM** ORM
- **PostgreSQL**
- **UUID primary key**
- **Clean Folder Structure** (modular service layer)

### Maybe next step:
- Redis (session cache / rate limit)
- WebSocket (real-time chat)
- S3-compatible storage (file upload)

---

## Struktur Folder 

/cmd
/server
main.go

/internal
/config
/database
/models
/repositories
/services
/handlers (controller)
/middlewares
/utils

/pkg (helper libs)


## Roadmap
- v0.1: Auth + User
- v0.2: Workspace + Members
- v0.3: Project + Board + Task
- v0.4: Comments + Activity
- v0.5: Chat realtime
- v1.0: Billing + Production

---

Dikembangkan oleh **dwiprst13** — dibuat khusus agar lebih simpel, cepat, dan terjangkau dibanding SaaS project management lain.

# PotTrack API

Authentication

- `POST /api/auth/register` - body `{email,nama,password}`
- `POST /api/auth/login` - body `{email,password}` returns `{token}`

Pots

- `GET /api/pot` - list user's pots
- `POST /api/pot` - create pot
- `GET /api/pot/:id` - get pot detail
- `PATCH /api/pot/:id` - update pot
- `POST /api/pot/:id/siklus` - start new cycle for pot

Cycles

- `POST /api/siklus/:id/arsip` - archive cycle (body `{alasan: selesai|gagal}`)
- `GET /api/siklus/:id/tahap` - get stages
- `PUT /api/siklus/:id/tahap` - replace stages
- `POST /api/siklus/:id/tugas` - create task
- `GET /api/siklus/:id/catatan` - list notes
- `POST /api/siklus/:id/catatan` - add note
- `GET /api/siklus/:id/insiden` - list incidents
- `POST /api/siklus/:id/insiden` - create incident
- `GET /api/siklus/:id/panen` - list harvest records
- `POST /api/siklus/:id/panen` - create harvest record

Tasks

- `GET /api/tugas?filter=today|overdue|all` - list tasks
- `PATCH /api/tugas/:id` - update task
- `POST /api/tugas/:id/selesai` - complete task (body `{catatan_teks?, foto_urls?}`)

Incidents

- `PATCH /api/insiden/:id/selesai` - close incident


All protected routes require header `Authorization: Bearer <token>` obtained from login.

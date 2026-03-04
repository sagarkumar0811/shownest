# ShowNest

ShowNest is an event discovery and ticket booking platform that allows users to explore shows, view venues, and reserve seats.  
The project focuses on building a scalable system for managing events, show schedules, and real-time ticket availability.

## PostgreSQL Installation

### Installation Details

| Component | Path/Value |
|-----------|------------|
| **Installation Directory** | `/Library/PostgreSQL/18` |
| **Server Installation Directory** | `/Library/PostgreSQL/18` |
| **Data Directory** | `/Library/PostgreSQL/18/data` |
| **Database Port** | `5433` |
| **Database Superuser** | `postgres` |
| **Operating System Account** | `postgres` |
| **Database Service** | `postgresql-18` |
| **Command Line Tools** | `/Library/PostgreSQL/18` |
| **pgAdmin4** | `/Library/PostgreSQL/18/pgAdmin 4` |
| **Stack Builder** | `/Library/PostgreSQL/18` |
| **Installation Log** | `/tmp/install-postgresql.log` |

## Database Architecture

```
PostgreSQL Server
│
├── postgres (admin)
│
├── user (database)
│    └── owned by user_service
│
├── booking (database)
│    └── owned by booking_service
│
├── catalog (database)
│    └── owned by catalog_service
│
├── merchant (database)
│    └── owned by merchant_service
│
├── payment (database)
│    └── owned by payment_service
│
├── seat (database)
│    └── owned by seat_service
│
├── search (database)
│    └── owned by search_service
│
└── notification (database)
     └── owned by notification_service
```

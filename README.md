# Shownest

A Go microservices backend for a full-stack event and venue booking platform. The system allows users to discover shows, browse venues, lock seats in real time, and complete bookings with QR-based tickets. Merchants can onboard, manage their venues and halls, publish events with showtimes, configure seating layouts, and apply dynamic pricing and promotional offers.

## Implemented

- **user** — handles user registration, login, OTP-based phone verification, JWT session management, and profile updates
- **merchant** — manages merchant onboarding, document submission, venue and hall creation, and hall configuration (capacity, seating type)
- **catalog** — stores and serves events (shows, concerts, sports), showtimes, and associated media; the source of truth for what is playing, where, and when
- **inventory** — manages seat definitions per hall, seat categories with pricing, real-time seat locking with Redis-backed TTLs, and seat availability queries
- **booking** — orchestrates the full booking lifecycle: seat reservation, payment handoff, confirmation, cancellation, and QR ticket generation with expiry
- **pricing** — computes final ticket prices by applying dynamic pricing rules and time-based surge logic, validates and redeems coupons, and manages convenience fee configuration

## Databases

Each service owns an isolated PostgreSQL database. Migrations live under `migrations/<service>/schema.sql`.

| Service | Database | Owner | Superuser |
|---|---|---|---|
| user | `user` | `user_service` | `postgres` |
| merchant | `merchant` | `merchant_service` | `postgres` |
| catalog | `catalog` | `catalog_service` | `postgres` |
| inventory | `inventory` | `inventory_service` | `postgres` |
| booking | `booking` | `booking_service` | `postgres` |
| pricing | `pricing` | `pricing_service` | `postgres` |

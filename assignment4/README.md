# Golang REST Backend (sql + postgres)

## Run

1. Create PostgreSQL database and user.
2. Set env vars:
   - DB_HOST
   - DB_PORT
   - DB_USER
   - DB_PASSWORD
   - DB_NAME
3. Run schema and sample data:
   - `psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f setup.sql`
4. Start server:
   - `go run .`
5. API:
   - `GET /users?page=1&page_size=10&name=alice&order_by=name`
   - `GET /users/common-friends?user1=<id>&user2=<id>`

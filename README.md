# MIBiG submission platform

This is the repository that contains the backend and frontend for the MIBiG submission platform

This is a platform that enables users to submit new entries to the MIBiG database.
Submissions can then be reviewed by users and curated

The software stack is as follows:
- Backend: golang
  - Gin-gonic
  - Gorm
- Frontend: Flask (currently, replacement TBD)
    - jinja2
    - bootstrap
- Database: TBD
- Testing: TBD 

## Installation

### Backend

1. Install golang 1.18+
2. Install gin-gonic and gorm:
```bash
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
```
3. Clone the repository:
```bash
git clone git@github.com:adraismawur/mibig-entry.git
```
4. Run the backend:
```bash
cd api
go run mibig-entry.go
``` 

### Frontend

TBD

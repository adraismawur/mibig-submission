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
- Database: sqlite/postgres

## Installation

Clone the repository
`git clone git@github.com:adraismawur/mibig-entry.gi`

### Backend

1. Install golang 1.18+
2. Clone the repository:
```bash
git clone git@github.com:adraismawur/mibig-entry.git
```
3. Change your working directory
`cd api`
3. Install go dependencies:
```bash
go get .
```
4. Run the backend:
```bash
go run mibig-submission.go
``` 

### Frontend

Using a different terminal:

1. Change your working directory
`cd web`
2. Create a conda/mamba environment
`mamba create -n submission -f environment.yml`
3. Activate the environment
`mamba activate submission`
4. Run the web app
`flask --app submission run --debug`


### Testing locally

Navigate to localhost:5000 to test the webpage.

By default, three users are created for testing purposes:

------------------------------
User                | password 
--------------------|---------
admin@localhost     | changeme
submitter@localhost | changeme
reviewer@localhost  | changeme

These have different abilities witin the application.

# Team Members

Example Golang CRUD app with simple HTML/JQuery front-end included to exercise features.

'Team Members' allows a company to manage employee/contractor resources.
Name, employee or contractor, roles, contract durations and tags are all included.

Uses gorilla-mux for routing and mysql-driver for RDBMS connectivity.
Routed, API service provided by process.go written in Golang.
Release also contains unit tests in process_test.go - usage instructions below.

Docker image file provided for quick setup and use - usage instructions below. Also included is the Dockerfile and mariadb-docker script to help others package, configure and import a mysql database into a docker container.

## File Structure

The overall file structure is as follows:

```text
/
├── index.html
├── process.go
├── process_test.go
├── mysql
│   └── db_import.sql
```

## Setup

Create a mysql database and user. The code currently expects a 'team_members' database name and 'tm' database username. Default password used is 'tm'. Change as desired.

Import mysql/db_import.sql into your mysql db with mysql -u tm team_members -p < mysql/db_import.sql

Build the Golang binary with 'go build process.go'. This will produce the runnable binary named 'process'.

If you do not have the gorilla-mux or mysql driver packages installed, install them per your O/S. Debian/Ubuntu uses the follow aptitude commands for installation:

'apt-get install golang-github-gorilla-mux-dev'

I used github.com/ziutek/mymysql/godrv for our mysql driver in this example. Install it with
'go get github.com/ziutek/mymysql/godrv'

## Run

Start the application by executing './process' within your shell. Browse to http://localhost:8000 to begin using Team Members.

## Adding team members.

One can add a team member by clicking 'Add''. 
Supply their type, either employee or contractor.
If an employee, provide their role. Some examples would be Software Engineer, Project Manager, etc.
If a contractor, provide their contact duration.
Finally, tag the member with any relevant information, skills such as MS Project, Jira, or coding skills such as Golang, C#, C++, Perl, Python and so on.

## Searching of team members.

Click 'Search' to open a text-dialog. Entering search parameters will find any team member that contains a match. Once found, you can edit their information by clicking 'Edit' or delete the team member by clicking 'Delete'. 

## Testing
One can run a small suite of unit tests against the service by executing 'go test'. Have a look at 'process_test.go' to review the tests that have been implemented.

## Docker image
Download the provided docker image from here https://peakelements.com/files/teams_image.tar.gz

Load with
'sudo docker load teams_image.tar'

Then, run the image with
'sudo docker run -p 8000:8000 teams'

or similar for your Docker install. Browse to http://localhost:8000 and add some team members.

## Building/running/saving/loading docker image

sudo docker build -t teams .

sudo docker run -p 8000:8000 teams

docker save teams > ~/teams_image.tar 

docker load teams < ~/teams_image.tar 

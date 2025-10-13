# Load Paranoia
This script is a paranoia eradicator because [like The Marías say "your paranoia is ..."](https://music.youtube.com/watch?v=thB6wpwJYEk&t=1m18s)

## Run
### Pre-requisite
- To run the script you need `go 1.25.1` installed
- Then in the Repo dir run `go get`
- You need to have respective access to all the GCP projects from where we will be pulling logs and running BigQuery queries

### Starting the program
- Add in the necessary variables
- `go run .` - to run the program

## Excel Formulas
- Change EpochMicro Timestamp to Excel format  
  `=(A1/(24*60*60))/1000000 + DATE(1970,1,1)`

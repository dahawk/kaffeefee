[![Go Report Card](https://goreportcard.com/badge/github.com/dahawk/kaffeefee)](https://goreportcard.com/report/github.com/dahawk/kaffeefee)

# README #
Kaffeefee (roughly translates to coffee fairy - german pun intended) is a nice little go project that offers a browser-based coffee tracker for companies, teams, or whoever wants to track coffee consumption.

## Install ##

### docker run

In order to run Kaffeefee you need:
* Build the project
* A Postgres database
* The content of create.sql imported into the Postgres database    
* Export an environment variable DB with the postgres connection url as value

### docker-compose 

Alternatively you can use the docker-compose file. In this case a database and the required tables will be created.

**PLEASE NOTE** For security reasons it is strongly recommended to genearate a secure password for the database and not to use the default credentials in the docker-compose.yml file.

In order to display users correctly, you need a image or avatar for every user in png format. These images need to be placed in the static folder.
The filename for each image needs to be <user name>.png (e.g. user name is christof -> filename is christof.png).
**Please note that user name and filename matching is case sensitive**

Alternatively Kaffeefee also checks Gravatar. If you enter an E-Mail address for a user and no matching image is found under static, Kaffeefee will lookup the md5 hash of the E-Mail address at Gravatar and, if available display an image from there. Check out how this great service works [here](https://en.gravatar.com/site/implement/).

### Legal matters ###
For the tally marks Kaffeefee makes use of tally.js (https://github.com/jeremy-brenner/tally.js). Unfortunately this project doesn't contain any license information at the moment.
For this reason kaffeefee is published without tally.js. All you need to do to get the tally marks funtionality is download the contents of https://github.com/jeremy-brenner/tally.js and place them unter /static.
After this step you should have the directories:
* /static/coffeescripts
* /static/fonts
* /static/javascripts
* /static/stylesheets

Under /static/stylesheets/tally.css update the font file paths so that they start with /static/fonts/...

As soon as I get information about licensing, I'll either include tally.js into the project or try to alter the code so that it works without (probably without tally marks), depending on the outcome.

The default avatar image is taken from https://pixabay.com/de/profil-mann-benutzer-home-42914/. According to the website the licensing for the image is CC0.

## Developer/Maintainer ##
Christof Horschitz (horschitz@gmail.com)

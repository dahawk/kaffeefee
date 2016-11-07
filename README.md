[![Build Status](https://travis-ci.org/dahawk/kaffeefee.svg?branch=master)](https://travis-ci.org/dahawk/kaffeefee) 
[![Go Report Card](https://goreportcard.com/badge/github.com/dahawk/kaffeefee)](https://goreportcard.com/report/github.com/dahawk/kaffeefee)

# README #
Kaffeefee (roughly translates to coffee fairy - german pun intended) is a nice little go project that offers a browser-based coffee tracker for companies, teams, or whoever wants to track coffee consumption.

## Install ##
In order to run Kaffeefee you need:
* Build the project
* A Postgres database
* The content of create.sql imported into the Postgres database
* One of the following:
    * In main.go substitute the string in line 52 with the correct postgres connection url
    * Alternatively export an environment variable DB with the postgres connection url as value

In order to display users correctly, you need a image or avatar for every user in png format. These images need to be placed in the static folder.
The filename for each image needs to be <user name>.png (e.g. user name is christof -> filename is christof.png).
**Please note that user name and filename matching is case sensitive**

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

## Developer/Maintainer ##
Christof Horschitz (horschitz@gmail.com)

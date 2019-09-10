# slackBot
In this porject I have tried to create a slack bot that teaches you some [Amharic](https://en.wikipedia.org/wiki/Amharic) and help you translate English to [Amharic](https://en.wikipedia.org/wiki/Amharic) with the help of [Google Translate API](https://cloud.google.com/translate/docs/)

I used [go 1.12](https://golang.org/dl/) to develop the bot app with the help of MySQL database and [Google Tranlsate API](https://cloud.google.com/translate/docs/)

Please follow the steps below to run the bot on you machine

You can escape the step below if you dont want the translate feature.
<p>Go to [Google Cloud](https://cloud.google.com/translate/docs/quickstart-client-libraries) to set up your project on Google Cloud and follow the steps 1 and 2 to get the private json key and set the environment variable (it is free for a year and you can deactivate it anytime).<p>

Clone the repo in a folder which is not your GOPATH and run the following command

- go build

This will create a <b>slackBot.exe</b> file and run that. This will start the server and the boot app will start listening to slack events. 

Create a workspace on [slack](https://slack.com/create#email). After creating the worksapce, go to the following [link](https://api.slack.com/apps) and create a slack app.
After creating the app:
- Enable Events under Event subscription. It will ask you to provide Request URL as shown in the picture below. I used [ngrok](https://ngrok.com/) to get public url. You can sign up and download a free version and obtain a public URL. it will be something like http://5bb67309.ngrok.io and you just have to add an endpoint to it (http://5bb67309.ngrok.io/learnAmharic) fill that on the Request URL.
- Once the URL is verified, you need to Add Bot User Events by click on the button and add the events shown below
 ![Enable Event](https://github.com/hailetotaw/slackBot/blob/master/EnableEvent.JPG)
- Add a Bot User as shown in the figure below and make sure to turn Always Show My Bot as Online on 
![Bot User](https://github.com/hailetotaw/slackBot/blob/master/BotUser.JPG)
- Under OAuth & Permissions add the scopes as show in the figure below
![Scope](https://github.com/hailetotaw/slackBot/blob/master/Scopes.JPG)
- Finally Install the APP on your workspace

After creating the app please set the following environment variables

> PORT=5000 <br />
> BOT_TOKEN=[Bot User OAuth Access Token]<br />
> VERIFICATION_TOKEN=[OAuth Access Token]]<br />
> BOT_ID=[@bot_user on slack]<br />
> CHANNEL_ID=[@channe on slack]<br />
> BOTDB_USER_PASSWORD=[password]<br />
> BOT_DB_USER=[dbuser]<br />
> DB_PORT=3306<br />
> DB_ADDRESS=[localhost]<br />
> DB_NAME=goBotDB<br />
> GOOGLE_APPLICATION_CREDENTIALS=[follow Google Cloud Set up above]<br />

<p>Create a channel on slack and set the channel id to <b>CHANNEL_ID</b> environment variable value</p>

Stop the running slackBot application for now and lets spin up the database by running the follwing command from the repo
- docker-compose up
<p>I am using docker compose to spin up mysql database</p>
<p>After Spinning the database start the <b>slackBot.exe</b> app and on the channel you created type <b>'[@bot_user] hello'</b></p>
![sample](https://github.com/hailetotaw/slackBot/blob/master/SlackResponse.JPG)
 


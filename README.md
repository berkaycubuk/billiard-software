# Billiard Software

Billiard Software is a billiard venue automation software that runs on the web.

## Technologies

Frontend: Vue
Backend: Go
Database: MySQL

## Supported payment methods

Currently it only supports Vipps as a payment method and it's directly embeded to the code itself. I'm looking to convert it into a module and make it easy to add another payment methods as modules.

## Deployment

Backend and the database can be deployed as containers. Frontend does not require specific method of deployment, tested with Nginx and Netlify.

## Development

In order to develop this on your local machine, you should start up the database and backend containers. After that run the frontend using `npm run dev` command. Don't forget to configure .env files in both frontend and backend.

## Support

I do not guarantee supporting this version of Billiard Software, so make sure you can play with the code before using it on the production. If you want dedicated and premium support feel free to contact with me <berkay@berkaycubuk.com>.

## Testing

Yes, I did not practiced for this part :^)

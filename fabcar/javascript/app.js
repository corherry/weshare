const express = require('express');
const app= express();
const bodyParser = require('body-parser');

const query = require('./query.js')
const initUser = require('./initUser.js')
const completeShare = require('./completeShare.js')
const shopping = require('./shopping.js')

app.use(bodyParser.json()); // for parsing application/json
app.use(bodyParser.urlencoded({ extended: true })); // for parsing application/x-www-form-urlencoded

app.get('/query', (req, res)=>{
	query(req.query.userId).then(result => {
		res.json(result);
	});
});
app.post('/initUser', (req, res)=>{
	initUser(req.body.userId).then(result => {
		res.json(result);
	});
});
app.post('/completeShare', (req, res)=>{
	console.log(req.body.userArr);
	completeShare(req.body.userArr).then(result => {
		res.json(result);
	});
});
app.post('/shopping', (req, res)=>{
	shopping(req.body.shoppingArr).then(result => {
		res.json(result);
	});
});
app.listen(8083, ()=>{
    console.log('Server is running at http://localhost:8083');
})
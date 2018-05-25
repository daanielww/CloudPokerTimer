const http = require('http');
var finalhandler = require('finalhandler');

const hostname = '127.0.0.1';
const port = 3000;

var serveStatic = require('serve-static');
var serve = serveStatic('public', {'index': ['index.html']});

const server = http.createServer(function onRequest (req, res) {
  serve(req, res, finalhandler(req, res));
});


server.listen(port, hostname, () => {
  console.log(`Server running at http://${hostname}:${port}`);
});


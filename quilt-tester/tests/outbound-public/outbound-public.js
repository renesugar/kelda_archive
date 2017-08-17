const quilt = require('@quilt/quilt');
let infrastructure = require('../../config/infrastructure.js');

let deployment = quilt.createDeployment();
deployment.deploy(infrastructure);

let connected = new quilt.Container('outbound', 'alpine', {
  command: ['tail', '-f', '/dev/null'],
}).replicate(infrastructure.nWorker*2);
quilt.publicInternet.allowFrom(connected, 80);

let notConnected = new quilt.Container('ignoreme', 'alpine', {
  command: ['tail', '-f', '/dev/null'],
}).replicate(infrastructure.nWorker*2);

deployment.deploy(connected);
deployment.deploy(notConnected);

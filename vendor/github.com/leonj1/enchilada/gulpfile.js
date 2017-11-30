var gulp = require('gulp'),
    fs = require("fs"),
    exec = require('child_process').execSync;

gulp.task('buildRepo', gulp.parallel(getDockerTag, gulp.series(removeEnchiladaBinary, buildApp)));
gulp.task('kickStart', gulp.parallel('buildRepo', gulp.series(stopRunningContainer, removeRunningContainers)));
gulp.task('default', gulp.series('kickStart', buildDockerImage, runIntegrationSetup));

function runIntegrationSetup(cb) {
  console.log('Running integration setup');
  exec('./integration/test.sh');
  cb();
}

function stopRunningContainer(cb) {
  console.log('Stoping Enchilada container');
  exec('docker stop enchilada || true');
  cb();
}

function removeRunningContainers(cb) {
  console.log('Removing Enchilada container');
  exec('docker rm enchilada || true');
  cb();
}

function removeEnchiladaBinary(cb) {
  try {
    console.log('Removing Enchilada binary');
  	fs.accessSync(enchilada)
    exec('rm enchilada');
  	// the file exists
  }catch(e){
  	// the file doesn't exists
  }
  cb();
}

function buildApp(cb) {
  console.log('Building Go app');
  exec("/bin/sh -c 'env GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -v -o enchilada server.go'");
  console.log('Done buiding Go app');
  cb();
}

function getDockerTag(cb) {
  console.log('Generating docker tag');
  var options = {
    pipeStdout: true
  };
  exec('python get_docker_build_version.py > file.contents', options);
  cb();
}

function buildDockerImage(cb) {
  var dockerTag = fs.readFileSync("file.contents", "utf8").replace(/\n$/, '');
  console.log('Building docker image with tag: ' + dockerTag);
  var cmd = 'docker build -t www.dockerhub.us:' + dockerTag + ' .';
  console.log('Command is: ' + cmd);
  exec(cmd);
  console.log('Done building docker image');
  cb();
}


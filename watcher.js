// libs
const childProcess = require('child_process');

function getMemory(pid, reportFn, stopFn) {
  // really shouldn't spawn a process for this so frequently
  // could have it run in a shell outputting to a file
  // we can read the file at the end, but would need to kill the process

  // pmap isn't doing what I think it is.
  // Haskell showed > 1TB memory on a system with 8gb RAM and 256 GB Disk
  childProcess.exec(`pmap ${pid} | tail -n 1`, (err, stdout, stderr) => {
    if (err) {
      console.log(err);
      process.exit(1);
    }
    const outParts = stdout.split(' ');
    const mem = parseInt(outParts[outParts.length - 1]);
    reportFn(mem);
    // check 100 times per second, checking terminator
    if (!stopFn()) setTimeout(() => getMemory(pid, reportFn, stopFn), 10);
  });
}

exports.runAndWatch = function(command) {
  // vars
  let pid = 0;
  let max = 0;
  let finished = false;

  // execution
  return new Promise((res, rej) => {
    console.log(command);
    childProcess.exec(command, (err, stdout, stderr) => {
      finished = true;
      if (err) {
        console.log(err);
        rej(err);
      } else if (!(max >= 0)) {
        console.log(err);
        rej('memory not recorded');
      } else {
        // append memory to output
        res(`${stdout}memory: ${max}`);
      }
    });

    // GET PID
    // pgrep -f = match full pattern

    // this can take 100ms to complete -_-

    // this isn't working for bash and perl... do not know why
    childProcess.exec(`pgrep -f "^${command}$"`, (err, stdout, stderr) => {
      if (err) {
        console.log(err);
        // process.exit(1);
      }
      pid = parseInt(stdout);
      // wait until we get pid to check memory
      const reportFn = x => {
        if (x > max) max = x;
      }
      const stopFn = () => finished;
      getMemory(pid, reportFn, stopFn);
    });
  });
}

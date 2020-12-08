// libs
const childProcess = require('child_process');

function getMemory(pid, reportFn, stopFn) {
  // really shouldn't spawn a process for this so frequently
  // could have it run in a shell outputting to a file
  // we can read the file at the end, but would need to kill the process
  childProcess.exec(`pmap -x ${pid} | tail -n 1`, (err, stdout, stderr) => {
    if (err) {
      console.log(err);
      process.exit(1);
    }
    // "total kB ... x ... x ... x"
    const [,,vsz,rss,dirty] = stdout.replace(/\n/gi, '').split(' ').filter(x => x !== '').map(x => parseInt(x));
    reportFn({ vsz, rss, dirty });
    // check 100 times per second, checking terminator
    if (!stopFn()) setTimeout(() => getMemory(pid, reportFn, stopFn), 10);
  });
}

exports.runAndWatch = function(command) {
  // vars
  let max = { vsz: 0, rss: 0, dirty: 0 };
  let finished = false;

  // execution
  return new Promise((res, rej) => {
    console.log(command);
    childProcess.exec(command, (err, stdout, stderr) => {
      finished = true;
      if (err) {
        console.log(err);
        rej(err);
      } else {
        const mem = `vsz-mem(k): ${max.vsz}\nrss-mem(k): ${max.rss}\ndirty-mem(k): ${max.dirty}`;
        // append memory to output
        // console.log(`mem: ${max.rss}`);
        res(`${stdout}${mem}`);
      }
    });

    // GET PID
    // pgrep -f = match full pattern
    // this can take 100ms to complete -_-

    // we could try to use pidof with just the first word
    // childProcess.exec(`pgrep -af "^${command}$"`, (err, stdout, stderr) => {

    // pidof is too fast so I had to add a timeout or it would fail on ruby: thanks ruby
    const commandStr = `pidof ${command.split(/\s+/)[0]}`;
    setTimeout(() => {
      childProcess.exec(commandStr, (err, stdout, stderr) => {
        if (err) {
          console.log(err);
          rej(err);
        }
        // console.log(command);
        // console.log(commandStr);
        // console.log(stdout);
        const pid = parseInt(stdout);
        // console.log(`pid: ${pid}`);
        // wait until we get pid to check memory
        const reportFn = x => {
          ['vsz', 'rss', 'dirty'].forEach(key => {
            if (x[key] > max[key]) max[key] = x[key];
          });
        }
        const stopFn = () => finished;
        getMemory(pid, reportFn, stopFn);
      });
    }, 150);
  });
}

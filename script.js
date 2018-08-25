const fs = require('fs');
const child_process = require('child_process');

// promisify

function access(file) {
  return new Promise((resolve, reject) => {
    fs.access(file, err => {
      if (err) {
        reject(err);
      } else {
        resolve(file);
      }
    });
  });
}

function exec(command) {
  console.log(command);
  return new Promise((resolve, reject) => {
    child_process.exec(command, (err, stdOut, stdErr) => {
      if (err) {
        reject(err);
      } else if (stdErr) {
        reject(stdErr);
      } else {
        resolve(stdOut);
      }
    });
  });
}

function readFile(file) {
  return new Promise((resolve, reject) => {
    fs.readFile(file, 'utf8', (err, data) => {
      if (err) {
        reject(err);
      } else {
        resolve(data);
      }
    });
  });
}

function writeFile(file, data) {
  return new Promise((resolve, reject) => {
    fs.writeFile(file, data, err => {
      if (err) {
        reject(err);
      } else {
        resolve(data);
      }
    });
  });
}

// async

const asyncMap = (fn, arr) => arr.reduce(
  (prev, curr) => prev.then(arr => fn(curr).then(res => arr.concat(res))),
  Promise.resolve([])
);

// utils

function fromPairs(pairs) {
  return pairs.reduce(
    (acc, [k,v]) => Object.assign({}, acc, ({ [k]:v })),
    {}
  );
}

function runLanguage(lang) {
  return new Promise((res, rej) => lang.isScript ? res() : rej())
    .catch(() => access(lang.bin))
    .catch(() => exec(lang.build))
    .then(() => exec(lang.run));
}

function textToHash(text) {
  const lines = text
    .split('\n')
    .filter(Boolean);
  const pairs = lines.map(
    ele => ele
      .split(':')
      .map(x => x.trim()));
  const hash = fromPairs(pairs);
  return hash;
}

function addKey(hash, key, fn) {
  return Object.assign({}, hash, { [key]: fn(hash) });
}

function compare(a, b) {
  if (a > b) {
    return 1;
  } else if (a == b) {
    return 0;
  } else {
    return -1;
  }
}

function sort(fn, arr) {
  return arr.slice().sort(fn);
}

function sortBy(fn, arr) {
  return sort((a, b) => compare(fn(a), fn(b)), arr);
}

function average(xs) {
  return xs.reduce((acc, x) => acc + x, 0) / xs.length;
}

// constants

CONFIG = 'config.json';
README = 'README.md'

// paths

function getConfig() {
  return readFile(CONFIG).then(
    x => JSON.parse(x).languages
  );
}

function build(xs) {
  console.log('\nbuilding...');
  return Promise.all(xs.map(
      x => x.build ? exec(x.build).then(() => x) : new Promise(res => res(x))
    ));
}

function run(xs) {
  console.log('\nrunning...');
  return asyncMap(
      lang => exec(lang.run)
        .then(textToHash)
        .then(x => addKey(x, 'name', () => lang.name))
        .then(x => addKey(x, 'speed', x => parseInt(x.iterations / x.seconds))),
      xs
    ).then(xs => sortBy(x => x.speed, xs).reverse());
}

function filter(xs) {
  return xs.filter(x => !x.ignore);
}

function readme(xs) {
  console.log('\nreadme:');
  const sampleSize = average(xs.map(x => parseFloat(x['sample-size'])));
  const percent = average(xs.map(x => parseFloat(x.percent))).toFixed(2);
  const tableData = xs.map(x => `${x.name}|${x.speed.toLocaleString()}`).join('\n');
  const fileData =
`#### Birthday Paradox - Monte Carlo simulations

* sample-size: ${sampleSize}
* probability: ${percent}

language | iterations/sec
|--|--|
${tableData}`;
  return writeFile(README, fileData);
}

const paths = new function() {
  this.config = () => getConfig();
  this.build = () => this.config().then(filter).then(build);
  this.run = () => this.build().then(run);
  this.readme = () => this.run().then(readme);
  this.repl = lang => this.config().then(
    xs => xs.find(x => x.name == lang)
  ).then(x => x && x.repl || `echo repl: '${lang}' not found`);
};

const commands = `[${Object.keys(paths).join(', ')}]`;

// main
(() => {
  const [,,path, ...args] = process.argv;
  paths[path] ?
    paths[path](...args).then(console.log) :
    console.log(`command: '${path}' failed\nlist of commands: ${commands}`);
})()

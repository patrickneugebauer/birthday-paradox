const fs = require('fs');
const child_process = require('child_process');
const util = require('util');

// promisify
const exec = x => {
  console.log(x);
  return util.promisify(child_process.exec)(x)
};
const readFile = util.promisify(fs.readFile);
const writeFile = util.promisify(fs.writeFile);

// async
const asyncMap = (fn, arr) => arr.reduce(
  (prev, curr) => prev.then(arr => fn(curr).then(res => arr.concat(res))),
  Promise.resolve([])
);

// utils
const fromPairs = xs => xs.reduce(
  (acc, [k,v]) => Object.assign({}, acc, ({ [k]:v })),
  {}
);

const textToHash = text => {
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

const addKey = (key, fn) => x => Object.assign({}, x, { [key]: fn(x) });

const compare = (a, b) => {
  if (a > b) {
    return 1;
  } else if (a == b) {
    return 0;
  } else {
    return -1;
  }
}

const sort = fn => xs => xs.slice().sort(fn);
const sortBy = fn => sort((a, b) => compare(fn(a), fn(b)));
const filter = fn => xs => xs.filter(fn);
const find = fn => xs => xs.find(fn);

const average = xs => {
  return xs.reduce((acc, x) => acc + x, 0) / xs.length;
}

// constants
const CONFIG = 'config.json';
const README = 'README.md'

// commands
const getConfig = () => {
  return readFile(CONFIG).then(
    x => JSON.parse(x).languages
  );
}

const build = xs => {
  console.log('\nbuilding...');
  return Promise.all(xs.map(
      x => x.build ? exec(x.build).then(() => x) : Promise.resolve(x)
    ));
}

const run = xs => {
  console.log('\nrunning...');
  return asyncMap(
      lang => exec(lang.run)
        .then(x => x.stdout)
        .then(textToHash)
        .then(addKey('name', () => lang.name))
        .then(addKey('speed', x => parseInt(x.iterations / x.seconds))),
      xs
    ).then(sortBy(x => x.speed))
    .then(x => x.reverse());
}

const readme = xs => {
  const sampleSize = average(xs.map(x => parseFloat(x['sample-size'])));
  const percent = average(xs.map(x => parseFloat(x.percent))).toFixed(2);
  const tableData = xs.map((x, i) => `${i + 1}|${x.name}|${x.speed.toLocaleString()}`).join('\n');
  const fileData =
`#### Birthday Paradox - Monte Carlo simulations

* sample-size: ${sampleSize}
* probability: ${percent}

| | language | iterations/sec |
|--|--|--|
${tableData}

thanks [Anthony Robinson](https://github.com/anthonycrobinson) for the tip about randint and random speed in python\n`;
  return writeFile(README, fileData).then(() => fileData);
}

const paths = new function() {
  this.config = () => getConfig();
  this.build = () => this.config().then(filter(x => !x.ignore)).then(build);
  this.run = () => this.build().then(run);
  this.readme = () => this.run().then(readme);
  this.repl = lang => this.config()
    .then(find(x => x.name == lang))
    .then(x => x && x.repl || `echo repl: '${lang}' not found`);
};

const commands = `[${Object.keys(paths).join(', ')}]`;

// main
(() => {
  const [,,path, ...args] = process.argv;
  paths[path] ?
    paths[path](...args).then(x => console.log('\nresults:\n',x)) :
    console.log(`command: '${path}' failed\nlist of commands: ${commands}`);
})()

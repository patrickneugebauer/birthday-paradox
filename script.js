const fs = require('fs');
const child_process = require('child_process');
const util = require('util');

const watcher = require('./watcher');

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
const CONFIG = 'config.test.json';
const README = 'README.test.md'
const VERSIONS = 'versions.md';

// commands
const getConfig = () => {
  return readFile(CONFIG).then(
    x => JSON.parse(x).languages
  );
}

const buildSync = xs => {
  console.log('\nbuilding sync...');
  return asyncMap(x => x.build ? exec(x.build).then(() => x) : Promise.resolve(x), xs);
}

const build = xs => {
  console.log('\nbuilding...');
  return Promise.all(xs.map(
      x => x.build ? exec(x.build).then(() => x) : Promise.resolve(x)
    ));
}

const run = xs => {
  console.log('\nrunning...');
  const iterationsScale = process.argv[3] || 0.25;
  return asyncMap(
      lang => watcher.runAndWatch(`${lang.run} ${parseInt(lang.executionsPerSecond * iterationsScale)}`)
        .then(textToHash)
        .then(addKey('name', () => lang.name))
        .then(addKey('speed', x => parseInt(x.iterations / x.seconds)))
        .then(addKey('year', () => lang.year))
        .then(addKey('execution', () => lang.execution))
        .then(addKey('solution', () => lang.solution))
        .then(addKey('hasRepl', () => Boolean(lang.repl))),
      xs
    ).then(sortBy(x => parseInt(x['rss-mem(k)'])))
    .then(x => x /* .reverse() */ );
}

const readme = xs => {
  const sampleSize = average(xs.map(x => parseFloat(x['sample-size'])));
  const percent = average(xs.map(x => parseFloat(x.percent))).toFixed(2);
  const tableData = xs.map(
    (x, i) => `| ${i + 1}
      ${x.name}
      ${x.speed.toLocaleString()}
      ${parseInt(x['rss-mem(k)']).toLocaleString()}
      ${x.year}
      ${x.solution}
      ${x.hasRepl ? 'x' : ''} |`.replace(/\s*\n\s*/gi, ' | ')
  ).join('\n');
  const fileData =
`#### Birthday Paradox - Monte Carlo simulations

* sample-size: ${sampleSize}
* probability: ${percent}

| | language | iterations/sec | rss-mem(k) | year | solution type | has repl |
|--| -- | -- | -- | -- | -- | -- |
${tableData}

thanks [Anthony Robinson](https://github.com/anthonycrobinson) for the tip about randint and random speed in python\n`;
  return writeFile(README, fileData).then(() => fileData);
}

const versions = xs => {
  const version = xs.filter(x => !x.ignore).map(x =>
    exec(x.version).then(
      x => x.stdout || x.stderr
    ).then(r => {
      const name = `#### ${x.name}`;
      const command = `\`${x.version}\``;
      const version = r
        .replace(/\r/g, '\n')
        .split('\n')
        .filter(x => Boolean(x))
        .map(x => `    ${x}`)
        .join('\n');
      return `${name}\n\n${command}\n\n${version}`;
    })
  );
  return Promise.all(version).then(xs => xs.join('\n\n') + '\n').then(
    x => writeFile(VERSIONS, x).then(() => x)
  );
}

const paths = new function() {
  this.config = () => getConfig();
  this.build = () => this.config().then(filter(x => !x.ignore)).then(build);
  this.buildSync = () => this.config().then(filter(x => !x.ignore)).then(buildSync);
  this.run = () => this.config().then(run);
  this.readme = () => this.run().then(readme);
  this.repl = lang => this.config()
    .then(find(x => x.name == lang))
    .then(x => x && x.repl || `echo repl: '${lang}' not found`);
  this.doc = this.help = this.man = () => Promise.resolve(`list of commands: ${commands}`);
  this.versions = () => this.config().then(versions);
};

const commands = `[${Object.keys(paths).join(', ')}]`;

// main
(() => {
  const [,,path, ...args] = process.argv;
  paths[path] ?
    paths[path](...args).then(x => console.log('\nresults:\n',x)) :
    console.log(`command: '${path}' failed\nlist of commands: ${commands}`);
})()

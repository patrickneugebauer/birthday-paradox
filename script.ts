import * as helpers from './script-helpers';
// import * as watcher from './watcher';
import { constants } from 'fs';

// constants
const CONFIG = 'config.json'; // 'config.test.json'
const README = 'README.md' // 'readme.test.md'
const VERSIONS = 'versions.md';

type Command = string;
type Config = {
  build: Command;
  execution: string;
  executionsPerSecond: number;
  ignore: boolean;
  isScript: boolean;
  name: string;
  repl: string;
  run: Command;
  solution: string;
  version: Command;
  year: number;
};
type RawResult = {
  iterations: string;
  "sample-size": string;
  percent: string;
  seconds: string;
  'vsz-mem(k)': string;
  'rss-mem(k)': string;
  'dirty-mem(k)': string;
}
type ResultData = RawResult & {
  name: string;
  speed: number;
  year: number;
  execution: string;
  solution: string;
  hasRepl: boolean;
}

const getConfig = (): Promise<Config[]> => {
  return helpers.readFile(CONFIG).then(
    x => JSON.parse(x.toString()).languages
  );
};

const buildSync = (xs: Config[]): Promise<Config[]> => {
  console.log('\nbuilding sync...');
  return helpers.asyncMap(x => x.build ? helpers.exec(x.build).then(() => x) : Promise.resolve(x), xs);
}

type PassResult = [Config, null];
type FailResult = [null, string];
const buildSyncWithFailures = (xs: Config[]): Promise<Config[]> => {
  console.log('\nbuilding sync with failures...');
  const mapResult = helpers.asyncMap(x =>
    (x.build)
      ? helpers.exec(x.build)
        .then(() => [x, null])
        .catch(err => [null, err]) as Promise<PassResult | FailResult>
      : Promise.resolve([x, null] as PassResult)
  , xs);
  return mapResult.then(x => {
    // (x.filter(a => a[0]) as PassResult[]).map(b => b[0]);
    const positiveResults = x.filter(a => a[0]) as PassResult[];
    return positiveResults.map(b => b[0]);
  });
}

const filterForDocker = (xs: Config[]): Promise<Config[]> => {
  console.log('\nfiltering for docker...');
  return Promise.all(xs.map(
    x => helpers.access(`${x.name}/Dockerfile`, constants.F_OK)
      .then(access => Promise.resolve([access, x] as [boolean, Config]))
  )).then(
    configs => configs.filter(x => x[0]).map(x => Object.assign({}, x[1], {
      build: `docker build ${x[1].name} -t bday/${x[1].name}`,
      run: `docker run --rm bday/${x[1].name}`
    }))
  );
}

const build = (xs: Config[]): Promise<Config[]> => {
  console.log('\nbuilding...');
  return Promise.all(xs.map(
    x => x.build ? helpers.exec(`echo "\${${x.build}}"`).then(() => x) : Promise.resolve(x)
  ));
}

const run = (xs: Config[]): Promise<ResultData[]> => {
  console.log('\nrunning...');
  const iterationsScale = parseFloat(process.argv[3]) || 0.25;
  return helpers.asyncMap(
    lang => {
      const iterations = parseInt((lang.executionsPerSecond * iterationsScale).toString());
      // return watcher.runAndWatch(`${lang.run} ${iterations}`)
      return helpers.exec(`${lang.run} ${iterations}`).then(x => x.stdout)
        .then(x => helpers.textToHash(x) as RawResult)
        // .then(x => { console.log(x); return x })
        .then(x => Object.assign({}, x, {
          name: lang.name,
          speed: parseInt((iterations / parseFloat(x.seconds)).toString()),
          year: lang.year,
          execution: lang.execution,
          solution: lang.solution,
          hasRepl: Boolean(lang.repl)
        }) as ResultData)
    }, xs
  ).then(helpers.sortBy(x => parseInt(x['speed'])))
  .then(x => x.reverse());
}

const readme = (xs: ResultData[]) => {
  const sampleSize = helpers.average(xs.map(x => parseFloat(x['sample-size'])));
  const percent = helpers.average(xs.map(x => parseFloat(x.percent))).toFixed(2);
  const tableData = xs.map(
    (x, i) => `| ${i + 1}
      ${x.name}
      ${x.speed.toLocaleString()}
      ${x.year}
      ${x.solution}
      ${x.hasRepl ? 'x' : ''} |`.replace(/\s*\n\s*/gi, ' | ')
  ).join('\n');
  const fileData =
    `#### Birthday Paradox - Monte Carlo simulations

* sample-size: ${sampleSize}
* probability: ${percent}

| | language | iterations/sec | year | solution type | has repl |
|--| -- | -- | -- | -- | -- |
${tableData}

thanks [Anthony Robinson](https://github.com/anthonycrobinson) for the tip about randint and random speed in python\n`;
  return helpers.writeFile(README, fileData).then(() => fileData);
}

const versions = (xs: Config[]) => {
  const versionInfoPromises = xs.filter(x => !x.ignore).map(info =>
    helpers.exec(info.version).then(
      execResult => execResult.stdout || execResult.stderr
    ).then(r => {
      const name = `#### ${info.name}`;
      const command = `\`${info.version}\``;
      const version = r
        .replace(/\r/g, '\n')
        .split('\n')
        .filter(x => Boolean(x))
        .map(x => `    ${x}`)
        .join('\n');
      return `${name}\n\n${command}\n\n${version}`;
    })
  );
  return Promise.all(versionInfoPromises)
    .then(versionInfo => versionInfo.join('\n\n') + '\n')
    .then(x => helpers.writeFile(VERSIONS, x).then(() => x));
}

class PublicFunctions {
  config = () => getConfig();
  build = () => this.config().then(helpers.filter(x => !x.ignore)).then(build);
  buildSync = () => this.config().then(helpers.filter(x => !x.ignore)).then(buildSync);
  readmeDocker = () => this.config()
    .then(helpers.filter(x => !x.ignore))
    .then(filterForDocker)
    .then(build)
    .then(run)
    .then(readme);
  readmeDockerSync = () => this.config()
    // .then(helpers.filter(x => !x.ignore))
    .then(filterForDocker)
    // .then(buildSync)
    .then(run)
    .then(readme);
  buildDockerSyncWithFailures = () => this.config()
    .then(helpers.filter(x => !x.ignore))
    .then(filterForDocker)
    .then(buildSyncWithFailures);
  readmeDockerSyncWithFailures = () => this.config()
    .then(helpers.filter(x => !x.ignore))
    .then(filterForDocker)
    .then(buildSyncWithFailures)
    .then(run)
    .then(readme);
  run = () => this.build().then(run);
  readme = () => this.run().then(readme);
  repl = (lang?: string) => this.config()
    .then(helpers.find(x => x.name == lang))
    .then(x => x && x.repl || `echo repl: '${lang}' not found`);
  doc = () => Promise.resolve(`list of commands: ${commands}`);
  help = this.doc;
  man = this.doc;
  versions = () => this.config().then(versions);
}

const mainObj = new PublicFunctions();
type Commands = keyof typeof mainObj;
const commands = `[${Object.keys(mainObj).join(', ')}]`;

// main
(() => {
  const [, , path, ...args] = process.argv;
  mainObj[path as Commands]
    ? mainObj[path as Commands](...args).then((x) => console.log('\nresults:\n', x))
    : console.log(`command: '${path}' failed\nlist of commands: ${commands}`);
})()

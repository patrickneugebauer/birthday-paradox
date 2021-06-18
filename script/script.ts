import * as helpers from './script-helpers';
import { Printer } from './printer';
// import * as watcher from './watcher';
import { constants } from 'fs';

// constants
const CONFIG = 'config.json';
const README = 'readme.md';
const VERSIONS = 'versions.md';

const iterationsScale = parseFloat(process.argv[2]) || 0.25;

type Command = string;
type Config = {
  build: Command;
  execution: string;
  executionsPerSecond: number;
  ignore: boolean;
  isScript: boolean;
  name: string;
  displayName?: string;
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
  size: number;
  year: number;
  execution: string;
  solution: string;
  hasRepl: boolean;
}

const printer = new Printer();

const getConfig = (): Promise<Config[]> => {
  printer.startLine(1, "parallel");
  return helpers.readFile(CONFIG).then(x => {
    printer.progressTick();
    return JSON.parse(x.toString()).languages
  });
};

const filterSolutions = (xs: Config[]): Promise<Array<Config & { weigh: string }>> => {
  printer.startLine(xs.length, "parallel");
  const checkForDockerfilesPromise = Promise.all(xs.map(
    x => helpers.access(`solutions/${x.name}/Dockerfile`, constants.F_OK)
      .then(access => {
        printer.progressTick();
        return Promise.resolve([access, x] as [boolean, Config]);
      })
  ));
  return checkForDockerfilesPromise.then(configs => {
    const filteredConfigs = configs
      .filter(x => {
        const keep = x[0];
        // if (!keep) console.log(`removing: ${x[1].name}`);
        return keep;
      })
      .map(x => x[1])
      .map(x => {
        const imageName = `bday/${x.name}`;
        return Object.assign({}, x, {
          build: `docker build solutions/${x.name} -t ${imageName}`,
          run: `docker run --rm ${imageName}`,
          weigh: `docker images | grep "${imageName} " | rev | cut -d " " -f 1 | rev`
        });
      });
    return filteredConfigs;
  });
}

type PassResult = [Config, null];
type FailResult = [null, string];
const buildSyncWithFailures = (xs: Config[]): Promise<Config[]> => {
  printer.startLine(xs.length, "series");
  const mapResult = helpers.asyncMap(x =>
    (x.build)
      ? helpers.exec(x.build)
        .then(() => {
          printer.progressTick();
          return [x, null];
        })
        .catch(err => [null, err]) as Promise<PassResult | FailResult>
      : Promise.resolve([x, null] as PassResult)
  , xs);
  return mapResult.then(x => {
    const positiveResults = x.filter(a => a[0]) as PassResult[];
    const positiveResultConfigs = positiveResults.map(b => b[0]);
    return positiveResultConfigs;
  });
}

const weighImages = (xs: Array<Config & { weigh: string }>): Promise<Array<Config & { size: string }>> => {
  printer.startLine(xs.length, "parallel");
  return Promise.all(xs.map(
    x => helpers
      .exec(x.weigh)
      .then(res => {
        printer.progressTick();
        return Object.assign({}, x, { size: res.stdout });
      })
  )).then(x => {
    return x;
  });
}

const versions = (xs: Config[]) => {
  printer.startLine(xs.length, "parallel");
  const versionInfoPromises = xs.map(info => {
    const [versionExecutable, versionParams] = helpers.headTail(info.version.split(" "));
    const versionCommand = `docker run --rm --entrypoint ${versionExecutable} bday/${info.name} ${versionParams.join(" ")}`;
    return helpers.exec(versionCommand).then(
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
      const versionString = `${name}\n\n${command}\n\n${version}`;
      printer.progressTick();
      return versionString;
    })
  });
  return Promise.all(versionInfoPromises)
    .then(versionInfo => versionInfo.join('\n\n') + '\n')
    .then(x =>
      helpers.writeFile(VERSIONS, x).then(() => {
        return x;
      })
    );
}

const run = (xs: Array<Config & { size: string }>): Promise<ResultData[]> => {
  printer.startLine(xs.length, "series");
  return helpers.asyncMap(lang => {
    const iterations = parseInt((lang.executionsPerSecond * iterationsScale).toString());
    return helpers.exec(`${lang.run} ${iterations}`).then(x => x.stdout)
      .then(x => helpers.textToHash(x) as RawResult)
      .then(x => {
        const speed = parseInt((iterations / parseFloat(x.seconds)).toString());
        const size = parseInt( (parseFloat(lang.size) * (lang.size.includes('GB') ? 1024 : 1) ).toString() )
        const resultData: ResultData = Object.assign({}, x, {
          name: lang.displayName || lang.name,
          speed,
          size,
          year: lang.year,
          execution: lang.execution,
          solution: lang.solution,
          hasRepl: Boolean(lang.repl)
        });
        printer.progressTick();
        return resultData;
      })
  }, xs)
    .then(helpers.sortBy(x => x['speed']))
    .then(x => {
      const descendingResults = x.reverse()
      return descendingResults;
    });
}
const readme = (xs: ResultData[]) => {
  printer.startLine(1, "parallel");
  const sampleSize = helpers.average(xs.map(x => parseFloat(x['sample-size'])));
  const percent = helpers.average(xs.map(x => parseFloat(x.percent))).toFixed(2);
  const tableData = xs.map(
    (x, i) => `| ${i + 1}
      ${x.name}
      ${x.speed.toLocaleString()}
      ${x.size}
      ${x.year}
      ${x.solution}
      ${x.hasRepl ? 'x' : ''} |`.replace(/\s*\n\s*/gi, ' | ')
  ).join('\n');
  const fileData =
    `#### Birthday Paradox - Monte Carlo simulations

* sample-size: ${sampleSize}
* probability: ${percent}

| | language | iterations/sec | image-size(MB) | year | solution type | has repl |
| :--: | :-- | --: | --: | --: | :-- | :--: |
${tableData}

thanks [Anthony Robinson](https://github.com/anthonycrobinson) for the tip about randint and random speed in python\n`;
  return helpers
    .writeFile(README, fileData)
    .then(() => {
      printer.progressTick();
      return fileData;
    });
}

// main
(() => {
  const commands = ['read config', 'filter', 'build', 'weigh', 'versions', 'run', 'readme'];
  printer.setNames(commands);
  printer.printPlan();

  getConfig()
    .then(filterSolutions)
    .then(x => buildSyncWithFailures(x) as Promise<typeof x>)
    .then(weighImages)
    .then(x => versions(x).then(() => x))
    .then(run)
    .then(readme)
    .then(_x => printer.finish())
})()

import * as helpers from './script-helpers';
// import * as watcher from './watcher';
import { constants } from 'fs';

// constants
const CONFIG = 'config.json'; // 'config.test.json'
const README = 'README.md' // 'readme.test.md'
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

const getConfig = (): Promise<Config[]> => {
  console.log('start: read config')
  return helpers.readFile(CONFIG).then(x => {
    console.log('finish: read config\n')
    return JSON.parse(x.toString()).languages
  });
};

const filterSolutions = (xs: Config[]): Promise<Array<Config & { weigh: string }>> => {
  console.log('start: filtering');
  const checkForDockerfilesPromise = Promise.all(xs.map(
    x => helpers.access(`${x.name}/Dockerfile`, constants.F_OK)
      .then(access => Promise.resolve([access, x] as [boolean, Config]))
  ));
  return checkForDockerfilesPromise.then(configs => {
    const filteredConfigs = configs
      .filter(x => {
        const keep = x[0];
        if (!keep) console.log(`removing: ${x[1].name}`);
        return keep;
      })
      .map(x => x[1])
      .map(x => {
        const imageName = `bday/${x.name}`;
        return Object.assign({}, x, {
          build: `docker build ${x.name} -t ${imageName}`,
          run: `docker run --rm ${imageName}`,
          weigh: `docker images | grep "${imageName} " | rev | cut -d " " -f 1 | rev`
        })
      });
    console.log('finish: filtering\n');
    return filteredConfigs;
  });
}

type PassResult = [Config, null];
type FailResult = [null, string];
const buildSyncWithFailures = (xs: Config[]): Promise<Config[]> => {
  console.log('start: building');
  const mapResult = helpers.asyncMap(x =>
    (x.build)
      ? helpers.exec(x.build, true)
        .then(() => [x, null])
        .catch(err => [null, err]) as Promise<PassResult | FailResult>
      : Promise.resolve([x, null] as PassResult)
  , xs);
  return mapResult.then(x => {
    // (x.filter(a => a[0]) as PassResult[]).map(b => b[0]);
    const positiveResults = x.filter(a => a[0]) as PassResult[];
    const positiveResultConfigs = positiveResults.map(b => b[0])
    console.log('finish: building\n');
    return positiveResultConfigs;
  });

}

const weighImages = (xs: Array<Config & { weigh: string }>): Promise<Array<Config & { size: string }>> => {
  console.log('start: weigh images');
  return Promise.all(xs.map(
    x => helpers
      .exec(x.weigh)
      .then(res => Object.assign({}, x, { size: res.stdout }) )
  )).then(x => {
    console.log('finish: weigh images\n');
    return x;
  });
}

const run = (xs: Array<Config & { size: string }>): Promise<ResultData[]> => {
  console.log('start: run');
  return helpers.asyncMap(lang => {
    const iterations = parseInt((lang.executionsPerSecond * iterationsScale).toString());
    // return watcher.runAndWatch(`${lang.run} ${iterations}`)
    return helpers.exec(`${lang.run} ${iterations}`, true).then(x => x.stdout)
      .then(x => helpers.textToHash(x) as RawResult)
      // .then(x => { console.log(x); return x })
      .then(x => {
        const speed = parseInt((iterations / parseFloat(x.seconds)).toString());
        const size = parseInt( (parseFloat(lang.size) * (lang.size.includes('GB') ? 1024 : 1) ).toString() )
        return Object.assign({}, x, {
          name: lang.displayName || lang.name,
          speed,
          size,
          year: lang.year,
          execution: lang.execution,
          solution: lang.solution,
          hasRepl: Boolean(lang.repl)
        }) as ResultData
      })
  }, xs)
    .then(helpers.sortBy(x => x['speed']))
    .then(x => {
      const descendingResults = x.reverse()
      console.log('finish: run\n');
      return descendingResults;
    });
}
const readme = (xs: ResultData[]) => {
  console.log('start: readme');
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
      console.log('finish: readme\n');
      return fileData
    });
}

const versions = (xs: Config[]) => {
  console.log('start: versions');
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
      return `${name}\n\n${command}\n\n${version}`;
    })
  });
  return Promise.all(versionInfoPromises)
    .then(versionInfo => versionInfo.join('\n\n') + '\n')
    .then(x =>
      helpers.writeFile(VERSIONS, x).then(() => {
        console.log('finish: versions\n');
        return x;
      })
    );
}

// main
(() => {
  getConfig()
    .then(filterSolutions)
    .then(x => buildSyncWithFailures(x) as Promise<typeof x>)
    .then(weighImages)
    .then(x => versions(x).then(() => x))
    .then(run)
    .then(readme)
    .then(_x => console.log('complete'))
})()

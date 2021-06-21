import * as helpers from './script-helpers';
import { Printer } from './printer';
import { constants } from 'fs';

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
type PassResult = [Config, null];
type FailResult = [null, string];

export class SimulationRunner {
  private readonly printer: Printer;
  private readonly iterationsScale: number;
  private readonly configPath: string | undefined;
  private readonly readmePath: string | undefined;
  private readonly versionPath: string | undefined;

  constructor({ printer, iterationsScale = 0.25, configPath, readmePath, versionPath }: {
    printer: Printer,
    iterationsScale?: number,
    configPath?: string,
    readmePath?: string,
    versionPath?: string
  }) {
    this.printer = printer;
    this.iterationsScale = iterationsScale;
    this.configPath = configPath;
    this.readmePath = readmePath;
    this.versionPath = versionPath;
  }

  initialize = () => {
    this.printer.printPlan();
  }
  uninitialize = () => {
    this.printer.finish();
  }

  getConfig = (): Promise<Config[]> => {
    if (!this.configPath) throw "configPath not set";
    this.printer.startLine(1, "parallel");
    return helpers.readFile(this.configPath).then(x => {
      this.printer.progressTick();
      return JSON.parse(x.toString()).languages
    });
  };

  filterSolutions = (xs: Config[]): Promise<Array<Config & { weigh: string }>> => {
    this.printer.startLine(xs.length, "parallel");
    const checkForDockerfilesPromise = Promise.all(xs.map(
      x => helpers.access(`solutions/${x.name}/Dockerfile`, constants.F_OK)
        .then(access => {
          this.printer.progressTick();
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

  build = (xs: Config[]): Promise<Config[]> => {
    this.printer.startLine(xs.length, "series");
    const mapResult = helpers.asyncMap(x =>
      (x.build)
        ? helpers.exec(x.build)
          .then(() => {
            this.printer.progressTick();
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

  weighImages = (xs: Array<Config & { weigh: string }>): Promise<Array<Config & { size: string }>> => {
    this.printer.startLine(xs.length, "parallel");
    return Promise.all(xs.map(
      x => helpers
        .exec(x.weigh)
        .then(res => {
          this.printer.progressTick();
          return Object.assign({}, x, { size: res.stdout });
        })
    )).then(x => {
      return x;
    });
  }

  versions = (xs: Config[]) => {
    if (this.versionPath === undefined) throw "versionPath not set";
    this.printer.startLine(xs.length, "parallel");
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
        this.printer.progressTick();
        return versionString;
      })
    });
    return Promise.all(versionInfoPromises)
      .then(versionInfo => versionInfo.join('\n\n') + '\n')
      .then(x =>
        helpers.writeFile(this.versionPath as string, x).then(() => {
          return x;
        })
      );
  }

  run = (xs: Array<Config & { size: string }>): Promise<ResultData[]> => {
    this.printer.startLine(xs.length, "series");
    return helpers.asyncMap(lang => {
      const iterations = parseInt((lang.executionsPerSecond * this.iterationsScale).toString());
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
          this.printer.progressTick();
          return resultData;
        })
    }, xs)
      .then(helpers.sortBy(x => x['speed']))
      .then(x => {
        const descendingResults = x.reverse()
        return descendingResults;
      });
  }

  readme = (xs: ResultData[]) => {
    if (this.readmePath === undefined) throw "readmePath not set";
    this.printer.startLine(1, "parallel");
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
      .writeFile(this.readmePath, fileData)
      .then(() => {
        this.printer.progressTick();
        return fileData;
      });
  }
}

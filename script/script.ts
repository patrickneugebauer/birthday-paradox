import { Printer } from './printer';
import { SimulationRunner } from './simulation-runner';


Array.prototype.removeFirst = function(fn) {
  const index = this.findIndex(fn);
  let match;
  if (index !== -1) {
    match = this[index];
    this.splice(index, 1);
  }
  return match;
}
declare global {
  interface Array<T> {
    removeFirst(fn: (x: T) => boolean): T | undefined;
  }
}

const userArgs = process.argv.slice(2);
const silentFlag = userArgs.removeFirst(x => x === '-s');
const debugFlag = userArgs.removeFirst(x => x === '-d');
const commandFlag = userArgs.removeFirst(x => x.startsWith("c="))?.replace("c=", "");
const scaleFlag = userArgs.removeFirst(x => x.startsWith("s="))?.replace("s=", "");
if (userArgs.length > 0) {
  console.log(`ERR: leftover flags [${userArgs.map(x => `'${x}'`).join(", ")}]`);
  process.exit(1);
}
if (silentFlag !== undefined && debugFlag !== undefined) {
  console.log(`ERR: cannot have silent and debug`);
  process.exit(1);
}
if (debugFlag) console.log(`debug: ${debugFlag}, silent: ${silentFlag}, command: ${commandFlag}, scale: ${scaleFlag}`);

const printerOutputArgs = silentFlag ? { write: () => undefined, moveCursor: () => undefined } : {};
const simulationRunnerArgs = {
  iterationsScale: scaleFlag ? parseFloat(scaleFlag) : undefined,
  configPath: 'config.json',
  readmePath: 'readme.md',
  versionPath: 'versions.md'
};

if (commandFlag !== undefined && !['readme', 'version'].includes(commandFlag)) {
  console.log("ERROR: command options: readme, verion")
} else if (commandFlag === 'readme') {
  const commands = ['read config', 'filter', 'build', 'weigh', 'run', 'readme'];
  const printer = new Printer({ names: commands, ...printerOutputArgs });
  const sr = new SimulationRunner({
    printer,
    ...simulationRunnerArgs
  });
  sr.initialize();
  sr.getConfig()
    .then(sr.filterSolutions)
    .then(x => sr.build(x) as Promise<typeof x>)
    .then(sr.weighImages)
    .then(sr.run)
    .then(sr.readme)
    .then(x => (sr.uninitialize(), x))
    .then(console.log)
} else if (commandFlag === 'version') {
  const commands = ['read config', 'filter', 'versions'];
  const printer = new Printer({ names: commands, ...printerOutputArgs });
  const sr = new SimulationRunner({
    printer,
    ...simulationRunnerArgs
  });
  sr.initialize();
  sr.getConfig()
    .then(sr.filterSolutions)
    .then(x => sr.versions(x))
    .then(x => (sr.uninitialize(), x))
    .then(console.log)
} else {
  const commands = ['read config', 'filter', 'build', 'weigh', 'versions', 'run', 'readme'];
  const printer = new Printer({ names: commands, ...printerOutputArgs });
  const sr = new SimulationRunner({
    printer,
    ...simulationRunnerArgs
  });
  sr.initialize();
  sr.getConfig()
    .then(sr.filterSolutions)
    .then(x => sr.build(x) as Promise<typeof x>)
    .then(sr.weighImages)
    .then(x => sr.versions(x).then(() => x))
    .then(sr.run)
    .then(sr.readme)
    .then(_x => sr.uninitialize())
}

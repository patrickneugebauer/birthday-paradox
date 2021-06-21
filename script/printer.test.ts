import { Printer } from './printer';

(() => {
  console.log("should display")
  const commands = ['read config', 'filter', 'build', 'weigh', 'versions', 'readme'];
  const printer = new Printer({ names: commands });
  printer.printPlan();
  printer.startLine(10, "series");
  Array(9).fill(null).forEach(() => printer.progressTick());
  printer.startLine(20, "parallel");
  Array(16).fill(null).forEach(() => printer.progressTick());
  printer.startLine(5, "series");
  Array(3).fill(null).forEach(() => printer.progressTick());
  printer.startLine(19, "parallel");
  Array(19).fill(null).forEach(() => printer.progressTick());
  printer.finish();
})();

(() => {
  console.log("should not display")
  const commands = ['read config', 'filter', 'build', 'weigh', 'versions', 'readme'];
  const log: string[] = [];
  const printer = new Printer({
    names: commands,
    write: x => log.push(`write: (${x})`),
    moveCursor: (x: number, y: number) => log.push(`moveCursor: (${x},${y})`)
  });
  printer.printPlan();
  printer.startLine(5, "series");
  Array(4).fill(null).forEach(() => printer.progressTick());
  printer.startLine(2, "parallel");
  Array(1).fill(null).forEach(() => printer.progressTick());
  printer.finish();
  console.log('printer.lines:');
  console.log(printer.lines);
  console.log('log:');
  console.log(log);
})();

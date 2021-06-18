import { Printer } from './printer';

(() => {
  const printer = new Printer();
  const commands = ['read config', 'filter', 'build', 'weigh', 'versions', 'readme'];
  printer.setNames(commands);
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
})()

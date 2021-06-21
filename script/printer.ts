type PrinterInfo = {
  names: string[];
  lines: number;
  maxLength: number;
}
type Order = "series" | "parallel";
class PrinterLine {
  tick = 0;
  progress = 0;
  constructor(
    readonly itemCount: number,
    readonly order: Order
  ) {}
}
type Write = (x: string) => void;
type MoveCursor = (x: number, y: number) => void;

Object.defineProperty(Array.prototype, 'last', {
  get: function() {
    return this[this.length - 1];
  }
});
declare global {
  interface Array<T> {
    last: T | undefined;
  }
}

export class Printer {
  // constants
  readonly PROGRESS_LENGTH = 10;
  readonly write: Write;
  readonly moveCursor: MoveCursor;
  // readonly references
  readonly summaryInfo: PrinterInfo;
  readonly lines: PrinterLine[] = [];
  // variables

  constructor({
    names,
    write = (x: string) => process.stdout.write(x),
    moveCursor = (x: number, y: number) => process.stdout.moveCursor(x, y)
  }: {
    names: string[],
    write?: Write,
    moveCursor?: MoveCursor
  }) {
    this.write = write;
    this.moveCursor = moveCursor;
    this.summaryInfo = {
      names: names,
      lines: names.length,
      maxLength: names.reduce((acc, x) => (x.length > acc) ? x.length : acc, 0)
    };
  }

  printPlan = () => {
    const rows = this.summaryInfo.names.map(x => `${x}${" ".repeat(this.summaryInfo.maxLength - x.length)} [          ]`);
    this.write(rows.join("\n")+"\n");
    this.moveCursor(this.summaryInfo.maxLength + 2, -this.summaryInfo.lines);
  };

  startLine = (itemCount: number, order: Order) => {
    if (this.lines.length >= this.summaryInfo.lines) throw "all lines have been filled"
    this.moveCursor(-(this.lines.last?.progress || 0), (this.lines.length > 0) ? 1 : 0);
    this.lines.push(new PrinterLine(itemCount, order));
    if (order === "parallel") {
      this.write("-".repeat(this.PROGRESS_LENGTH));
      this.moveCursor(-this.PROGRESS_LENGTH, 0);
    } else if (order === "series") {
      this.write("-");
      this.moveCursor(-1, 0);
    }
  }

  progressTick = () => {
    const lastLine = this.lines.last;
    if (lastLine === undefined) throw "line not started"
    const ticks = lastLine.tick + 1;
    const newProgress = Math.floor(ticks / lastLine.itemCount * 10);
    const diff = newProgress - lastLine.progress;
    Object.assign(lastLine, {
      tick: ticks,
      progress: newProgress
    });
    if (diff > 0) {
      this.write("=".repeat(diff));
      if (lastLine.order === "series" && newProgress < 10) {
        this.write("-");
        this.moveCursor(-1, 0);
      }
    }
  }

  finish = () => {
    const lastLine = this.lines.last;
    if (lastLine === undefined) throw "line not started"
    this.moveCursor(-(this.summaryInfo.maxLength + 2 + lastLine.progress), this.summaryInfo.lines - this.lines.length + 1);
  }
}

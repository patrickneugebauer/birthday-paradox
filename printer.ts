type PrinterSummaryInfo = {
  names: string[];
  lines: number | null;
  maxLength: number | null;
}
type PrinterLineInfo = {
  itemCount: number | null;
  tick: number;
  progress: number;
  order: "series" | "parallel" | null;
}


export class Printer {
  // constants
  readonly PROGRESS_LENGTH = 10;
  // readonly references
  readonly summaryInfo: PrinterSummaryInfo = {
    names: [],
    lines: null,
    maxLength: null,
  }
  readonly lineInfo: PrinterLineInfo = {
    itemCount: null,
    tick: 0,
    progress: 0,
    order: null,
  }
  // variables
  linesStarted = 0;
  initialized = false;
  // debugging
  message = '';

  setNames = (names: string[]) => {
    Object.assign(this.summaryInfo, {
      names: names,
      lines: names.length,
      maxLength: names.reduce((acc, x) => (x.length > acc) ? x.length : acc, 0)
    });
    this.initialized = true;
  }

  printPlan = () => {
    if (this.initialized === false) throw "printer must be initialized with setNames"
    const rows = this.summaryInfo.names.map(x => `${x}${" ".repeat(this.summaryInfo.maxLength! - x.length)} [          ]`);
    console.log(rows.join("\n"));
    process.stdout.moveCursor(this.summaryInfo.maxLength! + 2, -this.summaryInfo.lines!);
  };

  startLine = (itemCount: number, order: "series" | "parallel") => {
    if (this.initialized === false) throw "printer must be initialized with setNames"
    if (this.linesStarted >= this.summaryInfo.lines!) throw "all lines have been filled"
    process.stdout.moveCursor(-this.lineInfo.progress, (this.linesStarted > 0) ? 1 : 0);
    Object.assign(this.lineInfo, {
      itemCount: itemCount,
      tick: 0,
      progress: 0,
      order,
    });
    if (order === "parallel") {
      process.stdout.write("-".repeat(this.PROGRESS_LENGTH));
      process.stdout.moveCursor(-this.PROGRESS_LENGTH, 0);
    } else if (order === "series") {
      process.stdout.write("-");
      process.stdout.moveCursor(-1, 0);
    }
    this.linesStarted++;
    this.message+=(`start: ${-this.lineInfo.progress}\n`);
  }

  progressTick = () => {
    if (this.initialized === false) throw "printer must be initialized with setNames"
    if (this.linesStarted === 0) throw "line must be started with startLine before calling progressTick"
    const ticks = this.lineInfo.tick + 1;
    const newProgress = Math.floor(ticks / this.lineInfo.itemCount! * 10);
    const diff = newProgress - this.lineInfo.progress;
    Object.assign(this.lineInfo, {
      tick: ticks,
      progress: newProgress
    });
    if (diff > 0) {
      process.stdout.write("=".repeat(diff));
      if (this.lineInfo.order === "series" && newProgress < 10) {
        process.stdout.write("-");
        process.stdout.moveCursor(-1, 0);
      }
    }
    this.message+=(`tick: ${this.lineInfo.tick}/${this.lineInfo.itemCount}=${this.lineInfo.progress}%${diff}\n`);
  }

  finish = () => {
    if (this.initialized === false) throw "printer must be initialized with setNames"
    process.stdout.moveCursor(-(this.summaryInfo.maxLength! + 2 + this.lineInfo.progress), this.summaryInfo.lines! - this.linesStarted! + 1);
    // console.log(this.message);
  }
}

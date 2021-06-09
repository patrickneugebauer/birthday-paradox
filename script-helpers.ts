import fs from 'fs';
import child_process from 'child_process';
import util from 'util';

// ==================================================
// exported (public) methods
// ==================================================
type Command = string;
export const exec = (x: Command) => {
  console.log(x);
  return util.promisify(child_process.exec)(x)
};
export const readFile = util.promisify(fs.readFile);
export const writeFile = util.promisify(fs.writeFile);
export const access = (x: Command, accessType: any): Promise<boolean> => {
  return new Promise((res, _rej) => {
    fs.access(x, accessType, (err) => {
      // always resolve, just pass true or false
      res(err ? false : true);
    });
  });
};

export const asyncMap = <T>(fn: (x: T) => any, arr: Array<T>) =>
arr.reduce(
  (prev: Promise<any[]>, curr) => prev.then(result => fn(curr).then((res: any) => addToArray(res, result))),
  Promise.resolve([])
);

export const textToHash = (text: string) => {
  const lines = text
    .split('\n')
    .filter(Boolean); // remove blank lines
  const pairs = lines.map(
    ele => ele
      .split(':')
      .map(x => x.trim()) as [string, string]);
  const hash = fromPairs(pairs);
  return hash;
}

export const addKey = (key: string, fn: (x: any) => any) =>
  (x: any) => Object.assign({}, x, { [key]: fn(x) });

// returns curried functions
export const sortBy = (fn: (x: any) => any) =>
  sort((a, b) => compare(fn(a), fn(b)));
export const filter = (fn: (x: any) => any) =>
  (xs: any[]) => xs.filter(fn);
export const find = (fn: (x: any) => any) =>
(xs: any[]) => xs.find(fn);

export const average = (xs: any[]) =>
  xs.reduce((acc, x) => acc + x, 0) / xs.length;

const addToArray = <T>(a: T, b: T[]) => {
  const copy = b.slice();
  copy.push(a);
  return copy;
}

// ==================================================
// un-exported (private) methods
// ==================================================
const fromPairs = (xs: Array<[string, any]>) => xs.reduce(
  (acc, [k, v]) => Object.assign({}, acc, ({ [k]: v })),
  {}
);

type CompareInt = -1 | 0 | 1;
type CompareFn = <T>(a: T, b: T) => CompareInt;

const compare = <T>(a: T, b: T): CompareInt => {
  if (a > b) {
    return 1;
  } else if (a == b) {
    return 0;
  } else {
    return -1;
  }
};

const sort = (fn: CompareFn) =>
  (xs: any[]) => xs.slice().sort(fn);

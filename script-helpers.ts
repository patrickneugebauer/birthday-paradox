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
export const access = (x: Command, accessType: number | undefined): Promise<boolean> => {
  return new Promise((res, _rej) => {
    fs.access(x, accessType, (err) => {
      // always resolve, just pass true or false
      res(err ? false : true);
    });
  });
};

export const asyncMap = <I, O>(fn: (x: I) => Promise<O>, arr: Array<I>): Promise<O[]> =>
  arr.reduce(
    (prev: Promise<O[]>, curr) => prev.then(result => fn(curr).then((res) => addToArray(res, result))),
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

// returns curried functions
export const addKey = <O, V>(key: string, fn: (x: O) => V) =>
  (x: O) =>
    Object.assign({}, x, { [key]: fn(x) });

export const sortBy = <SB>(fn: (x: SB) => Comparable): (x: SB[]) => SB[] =>
  sort<SB>((a, b) => compare(fn(a), fn(b)));

export const filter = <T>(fn: (x: T) => boolean) =>
  (xs: T[]) =>
    xs.filter(fn);

export const find = <T>(fn: (x: T) => boolean) =>
  (xs: T[]) =>
    xs.find(fn);

export const average = (xs: number[]) =>
  xs.reduce((acc, x) => acc + x, 0) / xs.length;

const addToArray = <T>(a: T, b: T[]) => {
  const copy = b.slice();
  copy.push(a);
  return copy;
}

// ==================================================
// un-exported (private) methods
// ==================================================
const fromPairs = <P>(xs: Array<[string, P]>) => xs.reduce(
  (acc, [k, v]) => Object.assign({}, acc, ({ [k]: v })),
  {} as { [k: string]: P }
);

type Comparable = number | string;
type CompareInt = -1 | 0 | 1;
// use a mapped type to help inference
type CompareFn<T> = (a: T, b: T) => CompareInt;

const compare = (a: Comparable, b: Comparable): CompareInt => {
  if (a > b) {
    return 1;
  } else if (a == b) {
    return 0;
  } else {
    return -1;
  }
};

const sort = <S>(fn: CompareFn<S>) =>
  (xs: S[]): S[] =>
    xs.slice().sort(fn);

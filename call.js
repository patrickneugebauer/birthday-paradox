const exec = require('child_process').exec;
const access = require('fs').access;

// constants

// ==================================================
// iterations is the number of process that will be run
// the code is asynchronous so processes run concurrently
// each process occupies one thread, not one core
// ex: on a 4 code 8 thread machine
//     you can give up to 8 processes their own resources
const iterations = 8;
// ==================================================
const list = [{
  lang: 'c',
  source: 'loops.c',
  executable: 'loops.out',
  run: './loops.out',
  build: 'cc loops.c -o loops.out'
},{
  lang: 'java',
  source: 'Loops.java',
  executable: 'Loops.class',
  run: 'java Loops',
  build: 'javac Loops.java'
},{
  lang: 'javascript',
  source: 'loops.js',
  executable: 'loops.js',
  run: 'node loops.js',
  build: ''
},{
  lang: 'python',
  source: 'loops.py',
  executable: 'loops.py',
  run: 'python3 loops.py',
  build: ''
},{
  lang: 'ruby',
  source: 'loops.rb',
  executable: 'loops.rb',
  run: 'ruby loops.rb',
  build: ''
},{
  lang: 'typescript',
  source: 'loops.ts',
  executable: 'loops.ts.js',
  run: 'node loops.ts.js',
  build: 'tsc loops.ts --out loops.ts.js --target esnext'
}]

// functions
const callFn = (fn) => new Promise(
  res => exec(fn, x => res(x))
);
const getTime = () => new Date().getTime();
const asyncMap = (fn, arr) => arr.reduce(
  (prev, curr) => prev.then(arr => fn(curr).then(res => arr.concat(res))),
  Promise.resolve([])
);

// check-build
console.log('checking for executable files');
const buildResults = list.map(({ lang, executable, build }) => {
  const promise = new Promise(
    res => access(executable, err => res(err))
  ).then(err => {
    if (err) {
      console.log(`building ${lang}...`);
      return callFn(build)
    }
  });
  return promise;
}, list);

// run
const runLang = ({ lang, run }) => {
  console.log(`running ${lang}...`);
  const start = getTime();
  const promises = [];
  for (let i = 0; i < iterations; i++) {
    promises.push(callFn(run));
  }
  return Promise.all(promises).then(() => {
    const time = (getTime() - start) / 1000;
    const per = time / iterations;
    return { lang, per };
  });
};

const results = Promise.all(buildResults).then(() => {
  console.log('running files')
  return asyncMap(runLang, list);
});

// output
results.then((xs) => xs.forEach(
  ({ lang, per }) => console.log(`${lang}: ${per}`)
));

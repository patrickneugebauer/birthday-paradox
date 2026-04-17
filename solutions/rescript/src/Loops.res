@val @scope("process.hrtime") external hrtimeBigInt: unit => bigint = "bigint"
let start = hrtimeBigInt()

@val external argv: array<string> = "process.argv"
let iterations = argv
  -> Array.get(2)
  -> Option.flatMap(x => Int.fromString(x))
  -> Option.getOr(0);
let sampleSize = 23;

let randomDay = () => Js.Math.random_int(0, 365);
let createSample = n => {
  let rec createSampleFn = (n, l) => {
    switch(n) {
    | 0 => l
    | _ => createSampleFn(n - 1, list{randomDay(), ...l} )
    }
  }
  createSampleFn(n, list{});
}
let createSamples = (num, size) => {
  let rec createSamplesFn = (num, l) => {
    switch(num) {
    | 0 => l
    | _ => {
      let freshSample = createSample(size)
      createSamplesFn(num - 1, list{freshSample, ...l})
    }
    }
  }
  createSamplesFn(num, list{});
}

// https://rescript-lang.org/docs/manual/latest/pattern-matching-destructuring#match-on-list
let rec hasDuplicates = l => {
  switch(l) {
  | list{} => false
  | list{x, ...xs} => List.has(xs, x, (a, b) => a == b) ? true : hasDuplicates(xs)
  }
}
let numDuplicates = list => {
  let rec numDuplicatesFn = (l, count) => {
    switch(l) {
    | list{} => count
    | list{x, ...xs} => numDuplicatesFn(xs, count + (hasDuplicates(x) ? 1 : 0))
    }
  }
  numDuplicatesFn(list, 0);
}

// calculations
let matches = iterations
  -> createSamples(sampleSize)
  -> numDuplicates
let percent = (Int.toFloat(matches) /. Int.toFloat(iterations)) *. 100.0
let fin = hrtimeBigInt()
let seconds = BigInt.toFloat(fin - start) /. 1_000_000_000.0;
// formatted
let formattedIterations = iterations -> Int.toString
let formattedSampleSize = sampleSize -> Int.toString
let formattedPercent = percent -> Float.toFixed(~digits=2)
let formattedSeconds = seconds -> Float.toFixed(~digits=6)

// print
Js.log(`iterations: ${formattedIterations}`);
Js.log(`sample-size: ${formattedSampleSize}`);
Js.log(`percent: ${formattedPercent}`);
Js.log(`seconds: ${formattedSeconds}`);

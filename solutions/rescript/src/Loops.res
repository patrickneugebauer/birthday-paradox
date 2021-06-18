let start = Js.Date.make() |> Js.Date.getTime;
let iterations = int_of_string(Node.Process.argv[2]);
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
let createSamples = (size, num) => {
  let rec createSamplesFn = (num, l) => {
    switch(num) {
    | 0 => l
    | _ => createSamplesFn(num - 1, list{createSample(size), ...l} )
    }
  }
  createSamplesFn(num, list{});
}

// https://rescript-lang.org/docs/manual/latest/pattern-matching-destructuring#match-on-list
let rec hasDuplicates = l => {
  switch(l) {
  | list{} => false
  | list{x, ...xs} => List.mem(x, xs) ? true : hasDuplicates(xs)
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

let percent = iterations
  |> createSamples(sampleSize)
  |> numDuplicates
  |> float_of_int
  |> x => x /. float_of_int(iterations)
  |> x => x *. float_of_int(100);

Js.log(j`iterations: $iterations`); // toFixed for rounding
Js.log(j`sample-size: $sampleSize`);
Js.log(`percent: ${Js.Float.toFixedWithPrecision(percent, ~digits=2)}`);

let fin = Js.Date.make() |> Js.Date.getTime;
let seconds = (fin -. start) /. float_of_int(1000);
Js.log(`seconds: ${Js.Float.toFixedWithPrecision(seconds, ~digits=3)}`);

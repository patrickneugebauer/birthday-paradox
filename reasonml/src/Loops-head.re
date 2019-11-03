let start = Js.Date.make() |> Js.Date.getTime;
let iterations = 1000;
let sampleSize = 23;
Random.self_init;

let randomDay = () => Random.int(365);
let rec createSample = n => {
  switch(n) {
  | 0 => []
  | _ => [randomDay(), ...createSample(n - 1)]
  }
}
let rec createSamples = (size, num) => {
  switch(num) {
  | 0 => []
  | _ => [createSample(size), ...createSamples(size, num - 1)]
  }
}
let rec hasDuplicates = list => {
  switch(list) {
  | [] => false
  | [x, ...xs] => List.mem(x, xs) ? true : hasDuplicates(xs)
  }
}
let rec numDuplicates = list => {
  switch(list) {
  | [] => 0
  | [x, ...xs] => (hasDuplicates(x) ? 1 : 0) + numDuplicates(xs)
  }
}

let list = createSamples(sampleSize, iterations);
let duplicates = numDuplicates(list);

Printf.printf("iterations: %d\n", iterations);
Printf.printf("sample-size: %d\n", sampleSize);
let percent = float_of_int(duplicates) /. float_of_int(iterations) *. float_of_int(100);
Printf.printf("percent: %.f\n", percent);
let fin = Js.Date.make() |> Js.Date.getTime;
let seconds = (fin -. start) /. float_of_int(1000);
Printf.printf("seconds: %.3f\n", seconds);

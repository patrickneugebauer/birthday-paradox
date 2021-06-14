/* Console.log("Running Test Program:"); */
/* let () = print_endline(Lib.Util.hello()); */

let iterations = int_of_string(Sys.argv[1]);
let () = Lib.Util.simulate(iterations);

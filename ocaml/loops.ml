(* constants *)
let start = Unix.gettimeofday ();;
let iterations = 400000;;
let sampleSize = 23;;
Random.self_init();;

(* create intset module *)
module IntSet = Set.Make(
  struct let compare = Pervasives.compare
  type t = int
end)

(* ==================================================
functions
================================================== *)
let randomDay _ = Random.int 365;;

(* unused *)
let createSamplePartialList n =
  let rec loop n acc =
    if n > 0
      then
        let item = randomDay ()
        in if List.mem item acc
          then 1
          else loop (n - 1) (item :: acc)
      else 0
  in loop n [];;

(* unused *)
let rec checkDuplicates xs =
  match xs with
  | [] -> false
  | x :: [] -> false
  | x :: xs -> if List.mem x xs then true else checkDuplicates xs;;

(* unused *)
let createSampleFullList n =
  let rec loop n acc =
    if n > 0
      then loop (n - 1) (randomDay () :: acc)
      else if checkDuplicates acc then 1 else 0
  in loop n [];;

(* unused *)
let createSamplePartialArray n =
  let rec loop n acc =
    if n > 0
      then
        let item = randomDay ()
        in if Array.get acc item == 1
          then 1 (* duplicate found *)
          else let _ = Array.set acc item 1 in loop (n - 1) acc
      else 0 (* duplicate not found *)
  in loop n (Array.make 365 0);;

(* unused *)
let createSamplePartialSet n =
  let rec loop n acc =
    if n > 0
      then
        let item = randomDay ()
        in if IntSet.mem item acc
          then 1
          else loop (n - 1) (IntSet.add item acc)
      else 0
  in loop n IntSet.empty;;

(* unused *)
let createSampleFullSet n =
  let rec loop n acc =
    if n > 0
      then loop (n - 1) (IntSet.add (randomDay ()) acc)
      else if IntSet.cardinal acc == sampleSize then 0 else 1
  in loop n IntSet.empty;;

let rec createSamples n =
  if n > 0
    then createSamplePartialList sampleSize + createSamples (n - 1)
    else 0;;

(* calcs *)
let percent = (float_of_int (createSamples iterations)) /. (float_of_int iterations) *. (float_of_int 100);;
let finish = Unix.gettimeofday ();;
let seconds = (finish -. start);;

(* output *)
Printf.printf "iterations: %d\n" iterations;;
Printf.printf "sample-size: %d\n" sampleSize;;
Printf.printf "percent: %.2f\n" percent;;
Printf.printf "seconds: %.3f\n" seconds;;

(*
iterations      400000
sampleSize      23
-------------------------
data-structure  time (s)
list partial    1.150
list full       1.200
array partial   1.375
set partial     1.935
set full        2.100
*)

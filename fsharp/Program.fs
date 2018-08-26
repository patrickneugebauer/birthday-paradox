// Learn more about F# at http://fsharp.org

open System

let random = Random()

let rec duplicates n xs =
  match n with
  | 0 -> false
  | _ -> let num = random.Next(365)
         if List.contains num xs then true else duplicates (n-1) (num::xs)

let inline println s i =
  s + string(i) |> printfn "%s"

let simulate =
  let start = DateTime.Now
  let iterations = 100000
  let sampleSize = 23
  let mutable count = 0
  for i in 1..iterations do
    if duplicates sampleSize [] then count <- count + 1
  println "iterations: " iterations
  println "sample-size: " sampleSize
  float(count) / float(iterations) * 100.0
    |> fun x -> Math.Round(x, 2)
    |> println "percent: "
  let fin = DateTime.Now
  fin - start
    |> fun x -> x.TotalSeconds
    |> fun x -> Math.Round(x, 3)
    |> println "seconds: "
  ()

[<EntryPoint>]
let main argv =
    simulate
    0 // return an integer exit code

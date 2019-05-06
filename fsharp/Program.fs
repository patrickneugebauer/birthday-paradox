// Learn more about F# at http://fsharp.org

open System

let flip f x y = f y x

let roundTo n x =
  let shifter = 10.0 ** (float n)
  x
    |> (*) shifter
    |> round
    |> flip (/) shifter

let random = Random()

let rec duplicates n xs =
  match n with
  | 0 -> false
  | _ -> let num = random.Next(365)
         if Set.contains num xs then true else duplicates (n-1) (Set.add num xs)

let inline println s i =
  s + string(i) |> printfn "%s"

let simulate =
  let start = DateTime.Now
  let iterations = 50 * 1000
  let sampleSize = 23
  let mutable count = 0
  for i in 1..iterations do
    if duplicates sampleSize Set.empty then count <- count + 1
  println "iterations: " iterations
  println "sample-size: " sampleSize
  float(count) / float(iterations) * 100.0
    |> roundTo 2
    |> println "percent: "
  let fin = DateTime.Now
  fin - start
    |> fun x -> x.TotalSeconds
    |> roundTo 3
    |> println "seconds: "
  ()

[<EntryPoint>]
let main argv =
    simulate
    0 // return an integer exit code

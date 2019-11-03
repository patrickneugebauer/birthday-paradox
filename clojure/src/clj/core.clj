(ns clj.core)

; constants
(def sample-size 23)

; functions
(defn random-day [_] (rand-int 365))
(defn create-sample [_] (map random-day (range sample-size)))
(defn vector-has-duplicates [x] (apply distinct? x))
(defn current-time-milliseconds [] (System/currentTimeMillis))

; main
(defn -main [& args]
  (def iterations (read-string (first args)))
  ; data
  (def start (current-time-milliseconds))
  (def data (map create-sample (range iterations)))
  (def duplicates (count (filter vector-has-duplicates data)))
  (def percent (float (* (/ duplicates iterations) 100)))
  (def finish (current-time-milliseconds))
  (def seconds (float (/ (- finish start) 1000)))
  ; output
  (println (str "iterations: " iterations))
  (println (str "sample-size: " sample-size))
  (printf "percent: %.2f%n" percent)
  (printf "seconds: %.3f%n" seconds))

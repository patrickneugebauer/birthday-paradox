; constants
(def iterations 100000)
(def sample-size 23)

; functions
(defn random-day [] (rand-int 365))
(defn create-sample [] (repeatedly sample-size random-day))
(defn vector-has-duplicates [x] (apply distinct? x))
(defn current-time-milliseconds [] (System/currentTimeMillis))

; main method
(def start (current-time-milliseconds))
(def data (repeatedly iterations create-sample))
(def duplicates (count (filter vector-has-duplicates data)))
(def percent (float (* (/ duplicates iterations) 100)))
(def finish (current-time-milliseconds))
(def seconds (float (/ (- finish start) 1000)))

; output
(println (str "iterations: " iterations))
(println (str "sample-size: " sample-size))
(printf "percent: %.2f%n" percent)
(printf "seconds: %.3f%n" seconds)

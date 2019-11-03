; constants
(def iterations 100000)
(def sample-size 23)
(def threads 4)

; functions
(defn random-day [_] (rand-int 365))
(defn vector-has-duplicates [x] (apply distinct? x))
(defn current-time-milliseconds [] (System/currentTimeMillis))
(defn create-sample [_] (map random-day (range sample-size)))
(defn create-samples [] (map create-sample (range iterations)))
(defn create-future [_] (future (count (filter vector-has-duplicates (create-samples)))))

; main method
(def start (current-time-milliseconds))
(def futures (map create-future (range threads)))
(def duplicates (reduce + (map deref futures)))
(def percent (float (* (/ duplicates (* iterations threads)) 100)))
(def finish (current-time-milliseconds))
(def seconds (float (/ (- finish start) 1000)))

; output
(println (str "iterations: " iterations))
(println (str "sample-size: " sample-size))
(printf "percent: %.2f%n" percent)
(printf "seconds: %.3f%n" seconds)
(printf "ips: %.0f%n" (/ (* iterations threads) seconds))
(printf "ipspw: %.0f%n" (/ iterations seconds))
(shutdown-agents)

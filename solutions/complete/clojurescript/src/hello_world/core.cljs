(ns hello-world.core)

; vars
(when (nil? (aget js/process.argv 2))
  (.error js/console "missing iterations argument")
  (.exit js/process 1))
(def iterations (js/parseInt (aget js/process.argv 2)))
(def start (js/process.hrtime.bigint))
(def sample-size 23)

; data
(def count (atom 0))
(def calendar (make-array Integer/TYPE 365))

; create sample
(defn create-sample [i]
  (loop [x 0]
    (if (< x sample-size)
      (let [day (rand-int 365)]
        (if (not= (aget calendar day) i)
          (do
            (aset calendar day i)
            (recur (inc x)))
          (swap! count inc)))
      nil)))

; create samples
(dotimes [i iterations]
  (create-sample i))

; calcs
(def end (js/process.hrtime.bigint))
(def percent (/ (* @count 100) iterations))
(def seconds (/ (js/Number (- end start)) 1e9))

; format
(def formatted-percent (/ (Math/round (* percent 100)) 100))
(def formatted-seconds (/ (Math/round (* seconds 1000000)) 1000000))

; output
(.log js/console (str "iterations: " iterations))
(.log js/console (str "sample-size: " sample-size))
(.log js/console (str "percent: " formatted-percent))
(.log js/console (str "seconds: " formatted-seconds))

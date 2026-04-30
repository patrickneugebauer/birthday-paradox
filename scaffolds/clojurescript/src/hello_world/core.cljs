(ns hello-world.core)

; vars
(when (nil? (aget js/process.argv 2))
  (.error js/console "missing iterations argument")
  (.exit js/process 1))
(def iterations (js/parseInt (aget js/process.argv 2)))
(when (js/isNaN iterations)
  (.error js/console "Error: argument must be an integer")
  (.exit js/process 1))
(def start (js/process.hrtime.bigint))
(def sample-size 23)

; data
(def count 0)

; calcs
(def end (js/process.hrtime.bigint))
(def percent (/ (* count 100) iterations))
(def seconds (/ (js/Number (- end start)) 1e9))

; format
(def formatted-percent (/ (Math/round (* percent 100)) 100))
(def formatted-seconds (/ (Math/round (* seconds 1000000)) 1000000))

; output
(.log js/console (str "iterations: " iterations))
(.log js/console (str "sample-size: " sample-size))
(.log js/console (str "percent: " formatted-percent))
(.log js/console (str "seconds: " formatted-seconds))

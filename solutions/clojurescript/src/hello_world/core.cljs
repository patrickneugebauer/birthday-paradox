(ns hello-world.core)

(def name (aget js/process.argv 2))

(println (str "Hello World " name))

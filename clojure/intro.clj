; nil
nil

; keyword
:a

; var
(def a 1)

; symbol
a

; list
(+ 1 2 3)
(list 1 2 3)

; vector
[1 2 3]
(vector 1 2 3)

; set
#{ 1 2 3 }
(hash-set 1 2 3)

; map
{:a 1, :b 2}
(hash-map :a 1, :b 2)

; function
(fn [x y] 1 2 3 (+ x y))
(defn my-fn "my precious" [x y] 1 2 3 (+ x y))
#(+ %1 %2)
#(apply + %&)

; block - conventionally used for side-effects
(do 1 2 3)

; local-scope
(let [a 1 b 2] (+ 1 2) (+ a b))

; tail recursion
(defn fact "big fact" [x acc]
  (if (> x 0)
    (recur (dec x) (* acc x))
    acc))

; recur applies to innermost function
(defn flet "wrapped fact" [x]
  (let [f (fn [x acc]
    (if (> x 0)
      (recur (dec x) (* acc x))
      acc))]
  (f x 1)))

(defn floop "loopy fact" [num]
  (loop [x num, acc 1]
    (if (> x 0)
      (recur (dec x) (* acc x))
      acc)))

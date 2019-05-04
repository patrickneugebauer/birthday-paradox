(sort '(5 4 3 2 1) <)

(map (lambda (x) (+ x 1))
  '(1 2 3))
(map +
  '(1 2 3)
  '(1 2 3))
(map (lambda (x y z) (+ x y z))
  '(1 2 3)
  '(1 2 3)
  '(1 2 3))

(define sum 0)
(for-each (lambda (x)
  (set! sum (+ sum x)))
  '(1 2 3))
sum

(keep-matching-items
  '(1 2 3 4)
  (lambda (x) (if (>= x 2) #t #f)))
(delete-matching-items
  '(1 2 3)
  (lambda (x) (if (<= x 2) #t #f)))

(reduce + 0 '(1 2 3))

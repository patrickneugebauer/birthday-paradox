
; recursive
(define factorial
  (lambda (n)
    (if (> n 0)
      (* n (factorial (- n 1)))
      1)))

(define list*2
  (lambda (ls)
    (if (null? ls)
      '()
      (cons (* 2 (car ls)) (list*2 (cdr ls))))))

(define length.2
  (lambda (xs)
    (if (null? xs)
      0
      (+ 1 (length.2 (cdr xs))))))

; tail recursive
(define fact-rec
  (lambda (n acc)
    (if (= n 1)
      acc
      (let ((n-1 (- n 1)))
        (fact-ret n-1 (* acc n-1))))))

(define factorial-tail
  (lambda (x)
    (fact-ret x x)))

; does not work
; ===============
; (define factorial-inside
;   (lambda (x)
;     (let ((fi (lambda (n acc)
;       (if (= n 1)
;         acc
;         (let ((n-1 (- n 1)))
;           (fi n-1 (* acc n-1)))))))
;     (fi x x))))

; (f 4) = (4 * (f 3))
; (f 1 2 3) =

; named let tail-recursion
(define f-rekt (lambda (x)
  (let loop ((n x) (acc x))
    (if (= n 1)
      acc
      (let ((n-1 (- n 1))) (loop n-1 (* acc n-1))) ))))

; letrec tail recursion
(define f-snake (lambda (x)
  (letrec ((iter (lambda (n acc)
    (if (= n 1)
      acc
      (let ((n-1 (- n 1))) (iter n-1 (* acc n-1))) ) )))
  (iter x x) ) ))

; do tail recursion
(define f-do (lambda (x)
  (do ((n x (- n 1)) (acc x (* acc (- n 1))))
    ((= n 1) acc)) ))

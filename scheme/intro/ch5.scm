(if (null? ()) "yes" "no")
; #f = false
; #t = true
; (not #f) = #t

; the then and else values in an if should both be s-expressions
; for side effects use a begin expression

(define absolute
  (lambda (x)
    (if (>= x 0)
      x
      (* -1 x))))

(define reciprocal
  (lambda (x)
    (if (= x 0)
      #f
      (/ 1 x))))

(define int->char
  (lambda (x)
    (if (>= 33 x 126)
      #f
      (integer->char x))))

(define prod3
  (lambda (a b c)
    (if (and (positive? a) (positive? b) (positive? c))
      (* a b c)
      #f)))

(define prodn3
  (lambda (a b c)
    (if (or (negative? a) (negative? b) (negative? c))
    (* a b c)
    #f)))

(define hash
  (lambda (x)
    (cond
      ((equal? x "a") 1)
      ((equal? x "b") 2)
      (else #f))))

(define grade
  (lambda (x)
    (cond
      ((>= x 80) "a")
      ((>= x 60) "b")
      ((>= x 40) "c")
      ((< x 40) "d")
      (else #f))))

; eq? - objects are the same reference
; eqv? - compare object types and values
; equal? - used for comparing sequences

; char #\a
; symbol 'a
; () is an empty list, and is null, it is a list, but it is not a pair

; comparison operators, <, >, = and combs can be passed multiple values, order is checked

; string-ci? - will check case insensitive match

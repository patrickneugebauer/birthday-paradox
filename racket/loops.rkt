#lang racket

; functions
(define (random-day _) (random 365))
(define (random-sample x) (build-list x random-day))
(define (round-to n x)
  (let ((shifter (expt 10 n)))
    (/
      (round (* x shifter))
      shifter)))

; data
(define start-ms (current-inexact-milliseconds))
(define iterations (string->number (vector-ref (current-command-line-arguments) 0)))
(define sample-size 23)
(define data (build-list iterations (lambda (_) (check-duplicates (random-sample sample-size)))))
(define duplicates (length (filter identity data)))
(define percent (* (/ duplicates iterations) 100))
(define end-ms (current-inexact-milliseconds))
(define seconds (/ (round (- end-ms start-ms)) 1000))

; output
(printf "iterations: ~a~n" iterations)
(printf "sample-size: ~a~n" sample-size)
(printf "percent: ~a~n" (round-to 2 (exact->inexact percent)))
(printf "seconds: ~a~n" (exact->inexact seconds))

#lang racket
; ==================================================
; functions
; ==================================================

(define (random-day) (random 365))

(define (repeat num fn)
  (let loop ((n num) (f fn) (acc '()))
    (if (> n 0)
      (loop (- n 1) f (cons (fn) acc))
      acc)))

(define (round-to n x)
  (let ((shifter (expt 10 n)))
    (/
      (round (* x shifter))
      shifter)))

; ==================================================
; main method code
; ==================================================
(define start (current-inexact-milliseconds))

(define iterations
  (string->number
    (vector-ref
      (current-command-line-arguments)
      0)))

(define sample-size 23)

(define data
  (repeat iterations
    (lambda () (repeat sample-size random-day))))

(define percent
  (*
    (/
      (length (filter check-duplicates data))
      iterations)
    100))

(define formatted-percent
  (round-to 2
    (exact->inexact percent)))

(display
  (string-append
    "iterations: "
    (number->string iterations)
    "\n"))
(display
  (string-append
    "sample-size: "
    (number->string sample-size)
    "\n"))
(display
  (string-append
    "percent: "
    (number->string formatted-percent)
    "\n"))
(define seconds
  (exact->inexact
    (/
      (round
        (-
          (current-inexact-milliseconds)
          start))
      1000)))
(display
  (string-append
    "seconds: "
    (number->string seconds)
    "\n"))

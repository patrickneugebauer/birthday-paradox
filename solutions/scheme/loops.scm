; ==================================================
; functions
; ==================================================

(define (readfile filename)
  (with-input-from-file filename
    (lambda ()
      (let loop ((lines '()) (next-line (read-line)))
        (if (eof-object? next-line)
          (reverse lines)
          (loop (cons next-line lines) (read-line)))))))

(define (random-day) (random 365))

(define (repeat num fn)
  (let loop ((n num) (f fn) (acc '()))
    (if (> n 0)
      (loop (- n 1) f (cons (fn) acc))
      acc)))

(define (elem val arr)
  (let loop ((x val) (xs arr) (acc #f))
    (if (> (length xs) 0)
      (loop x (cdr xs) (or acc (= (car xs) val)))
      acc)))

(define (has-duplicates xs)
  (if (> (length xs) 0)
    (let ((head (car xs)) (tail (cdr xs)))
      (if (elem head tail)
        #t
        (has-duplicates tail)))
    #f))

(define (round-to n x)
  (let ((shifter (expt 10 n)))
    (/
      (round (* x shifter))
      shifter)))

; data
(define start (real-time-clock))
(define iterations (string->number (car (command-line-arguments))))
(define sample-size 23)
(define data (repeat iterations (lambda () (repeat sample-size random-day))))

; calcs
(define percent (* (/ (length (filter has-duplicates data)) iterations) 100))
(define formatted-percent (round-to 2 (exact->inexact percent)))
(define fin (real-time-clock))
(define seconds (internal-time/ticks->seconds (- fin start)))
(define formatted-seconds (round-to 6 (exact->inexact seconds)))

; output
(display (string-append "iterations: " (number->string iterations) "\n"))
(display (string-append "sample-size: " (number->string sample-size) "\n"))
(display (string-append "percent: " (number->string formatted-percent) "\n"))
(display (string-append "seconds: " (number->string formatted-seconds) "\n"))

(define (main args)
  (display args))

(exit)

; (working-directory-pathname)
; (cd "absolute or relative path")
; (load "file name")
; (clear)

(define vhello "Hello world")
(define fhello (lambda () "Hello world"))

(define hello
  (lambda (name)
    (string-append "Hello " name "!")))

((lambda (x) x) "x")

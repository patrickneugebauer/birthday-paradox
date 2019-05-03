; each grouping of cons is called a cons cell
; you can make a list of beaded cons cells

; cons - construciton
; car - contents of address of the register
; cdr - contents of decrement of the register

(cons 1 (cons 2 ())) ; (1 2)
(cons (cons 1 (cons 10 ())) 100)

(define my-list (cons 1 (cons 2 (cons 3 (cons 4 ())))))
(cons "Sum of" (cons my-list (cons "is" (cons 10 ()))))

'(1) ; is equivalent to (cons 1 ())
(list 1 2 3 4 5) ; also works

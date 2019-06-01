;; constants
(defvar start (get-internal-real-time))
(defvar iterations 250000)
(defvar sample-size 23)

;; generate data
(defun random-day () (random 365))
(defun generate-sample ()
  (loop repeat sample-size collect (random-day)))
(defun generate-samples ()
  (loop repeat iterations collect (generate-sample)))
(defun has-duplicates (x)
  (= (list-length (remove-duplicates x)) (list-length x)))
(defvar duplicates
  (list-length
    (remove-if
      #'has-duplicates
      (generate-samples))))

;; calcs
(defvar percent (* (/ duplicates iterations) 100))
(defvar finish (get-internal-real-time))
(defvar milliseconds (- finish start))
(defvar seconds (/ milliseconds 1000))

;; output
(format t "iterations: ~d~%" iterations)
(format t "sample-size: ~d~%" sample-size)
(format t "percent: ~,2f~%" percent)
(format t "seconds: ~,3f~%" seconds)

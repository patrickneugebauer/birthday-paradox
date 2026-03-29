;; constants
(defvar start (get-internal-real-time))
(defvar iterations (parse-integer (nth 1 *posix-argv*)))
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
(defvar time-units (- finish start))
(defvar seconds (/ time-units internal-time-units-per-second))

;; output
(format t "iterations: ~d~%" iterations)
(format t "sample-size: ~d~%" sample-size)
(format t "percent: ~,2f~%" percent)
(format t "seconds: ~,3f~%" seconds)

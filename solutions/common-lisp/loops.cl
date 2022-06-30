;; constants
(defvar start (get-internal-real-time))
(defvar iterations (parse-integer (nth 1 *posix-argv*)))
(defvar sample-size 23)

;; generate data
(defun random-day () (random 365))
(defvar counter 0)

(defun generate-sample (data)
  (declare (type (array integer 1) data))
  (loop repeat sample-size do
    (let ((index (random-day)))
      (if (eql (aref data index) 1)
          (progn (incf counter)
                 (loop-finish))
          (setf (aref data index) 1)))))

(defun generate-samples ()
  (loop repeat iterations do
    (generate-sample (make-array 365 :element-type 'integer))))

(generate-samples)

;; calcs
(defvar percent (* (/ counter iterations) 100))
(defvar finish (get-internal-real-time))
(defvar time-units (- finish start))
(defvar seconds (/ time-units internal-time-units-per-second))

;; output
(format t "iterations: ~d~%" iterations)
(format t "sample-size: ~d~%" sample-size)
(format t "percent: ~,2f~%" percent)
(format t "seconds: ~,3f~%" seconds)

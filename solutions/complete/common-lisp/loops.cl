;; constants
(defvar start (get-internal-real-time))
(defvar iterations (parse-integer (nth 1 *posix-argv*)))
(defvar sample-size 23)

;; generate data
(defun random-day () (random 365))
(defvar duplicates 0)
(defvar calendar (make-array 365 :element-type 'integer))

(defun generate-sample (i)
  (declare (type (array integer 1) calendar))
  (loop repeat sample-size do
    (let ((index (random-day)))
      (if (eql (aref calendar index) i)
          (progn (incf duplicates)
                 (loop-finish))
          (setf (aref calendar index) i)))))

(defun generate-samples ()
  (loop repeat iterations for i from 0 do (generate-sample i)))

(generate-samples)

;; calcs
(defvar percent (* (/ duplicates iterations) 100))
(defvar finish (get-internal-real-time))
(defvar seconds (/ (- finish start) internal-time-units-per-second))

;; output
(format t "iterations: ~d~%" iterations)
(format t "sample-size: ~d~%" sample-size)
(format t "percent: ~,2f~%" percent)
(format t "seconds: ~,6f~%" seconds)

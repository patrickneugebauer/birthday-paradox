" variables "
start := Time millisecondClock.
args := Smalltalk arguments.
iterations := (args at: 1) asNumber.
sampleSize := 23.
generator := Random new.

" duplicate check "
SequenceableCollection extend [
  SequenceableCollection class >> hasDuplicates: arr [
  | tail |
  ^((arr size) = 0)
    ifTrue: [ ^false ];
    ifFalse: [
      tail := arr allButFirst.
      (tail includes: (arr first))
        ifTrue: [ ^true ];
        ifFalse: [ SequenceableCollection hasDuplicates: tail ].
    ].
  ]
  hasDuplicates [ ^SequenceableCollection hasDuplicates: self ]
]

" data "
duplicates := ((1 to: iterations) select: [ :i |
  (1 to: sampleSize collect: [ :s | generator between: 1 and: 365 ]) hasDuplicates.
]) size.

" calcs "
percent := duplicates / iterations * 100.
formattedPercent := percent asScaledDecimal: 2.
fin := Time millisecondClock.
seconds := (fin - start) / 1000.
formattedSeconds := seconds asScaledDecimal: 3.

" output "
('iterations: ',(iterations asString)) displayNl.
('sample-size: ',(sampleSize asString)) displayNl.
('percent: ',(formattedPercent asString)) displayNl.
('seconds: ',(formattedSeconds asString)) displayNl.

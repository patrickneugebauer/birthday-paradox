Imports System

Module Program
  Sub Main(args As String())
    Dim iterations as Integer = Int32.parse(args(0))
    Simulate(iterations)
  End Sub

  Sub Simulate(iterations As Integer)
    Dim start as DateTime = DateTime.Now
    Dim sampleSize as Integer = 23
    Dim count as Integer = 0
    Dim rnd as Random = new Random()

    For i as Integer = 0 To iterations
      Dim data(365) as Integer
      For l as Integer = 0 To sampleSize
        Dim num as Integer = rnd.Next(0, 365)
        If data(num) = 1
          count += 1
          Exit For
        Else
          data(num) = 1
        End If
      Next
    Next

    Console.WriteLine($"iterations: {iterations}")
    Console.WriteLine($"sample-size: {sampleSize}")
    Dim percent as Double = (count / iterations) * 100
    Console.WriteLine($"percent: {Math.Round(percent, 2)}")
    Dim endTime as DateTime = DateTime.Now
    Dim diff as Double = (endTime - start).TotalSeconds
    Console.WriteLine($"seconds: {Math.Round(diff, 3)}")
  End Sub
End Module
